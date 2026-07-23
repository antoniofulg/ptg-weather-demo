package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/antoniofulg/ptg-weather-demo/geo"
	"github.com/antoniofulg/ptg-weather-demo/weatherapi"
)

// weatherHandler returns an httptest handler that replies with the given status
// and body for the Open-Meteo forecast endpoint.
func weatherHandler(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}
}

// swapServers points the geo and weatherapi packages at the given test servers
// for the duration of the test. A nil server leaves that package untouched.
func swapServers(t *testing.T, weather, geoSrv *httptest.Server) {
	t.Helper()
	oldWBase, oldWClient := weatherapi.BaseURL, weatherapi.HTTPClient
	oldGBase, oldGClient := geo.BaseURL, geo.HTTPClient
	if weather != nil {
		weatherapi.BaseURL = weather.URL
		weatherapi.HTTPClient = weather.Client()
	}
	if geoSrv != nil {
		geo.BaseURL = geoSrv.URL
		geo.HTTPClient = geoSrv.Client()
	}
	t.Cleanup(func() {
		weatherapi.BaseURL, weatherapi.HTTPClient = oldWBase, oldWClient
		geo.BaseURL, geo.HTTPClient = oldGBase, oldGClient
	})
}

// TestRunExplicitCoordinate covers E2E-001: with faked geo/weatherapi servers,
// invoking the CLI with --lat/--lon prints a line containing the formatted
// temperature and omits the location clause. The geo server fails the test if it
// is ever contacted, proving explicit coordinates skip geolocation entirely.
func TestRunExplicitCoordinate(t *testing.T) {
	weatherSrv := httptest.NewServer(weatherHandler(http.StatusOK,
		`{"current":{"temperature_2m":21.3},"current_units":{"temperature_2m":"°C"}}`))
	defer weatherSrv.Close()
	geoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("geo.Locate must not be called for an explicit coordinate")
		http.Error(w, "unexpected", http.StatusInternalServerError)
	}))
	defer geoSrv.Close()
	swapServers(t, weatherSrv, geoSrv)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--lat", "38.72", "--lon", "-9.14"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run() exit code = %d, want 0 (stderr: %q)", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "21.3°C") {
		t.Errorf("output %q does not contain formatted temperature %q", out, "21.3°C")
	}
	if want := "It is 21.3°C.\n"; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
	if strings.Contains(out, " in ") {
		t.Errorf("explicit-coordinate output must omit the location clause: %q", out)
	}
}

// TestRunGeolocated covers R2/R3: with both coordinate flags unset, the CLI
// resolves the location via geo.Locate and prints the full "in City, Country"
// clause.
func TestRunGeolocated(t *testing.T) {
	weatherSrv := httptest.NewServer(weatherHandler(http.StatusOK,
		`{"current":{"temperature_2m":21.3},"current_units":{"temperature_2m":"°C"}}`))
	defer weatherSrv.Close()
	geoSrv := httptest.NewServer(weatherHandler(http.StatusOK,
		`{"status":"success","lat":38.72,"lon":-9.14,"city":"Lisbon","country":"Portugal"}`))
	defer geoSrv.Close()
	swapServers(t, weatherSrv, geoSrv)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), nil, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run() exit code = %d, want 0 (stderr: %q)", code, stderr.String())
	}
	if want := "It is 21.3°C in Lisbon, Portugal.\n"; stdout.String() != want {
		t.Errorf("output = %q, want %q", stdout.String(), want)
	}
}

// TestRunFahrenheit covers R1/US-001 AC-2: --unit f converts and labels the
// temperature in Fahrenheit.
func TestRunFahrenheit(t *testing.T) {
	weatherSrv := httptest.NewServer(weatherHandler(http.StatusOK,
		`{"current":{"temperature_2m":21.34},"current_units":{"temperature_2m":"°C"}}`))
	defer weatherSrv.Close()
	swapServers(t, weatherSrv, nil)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--lat", "38.72", "--lon", "-9.14", "--unit", "f"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run() exit code = %d, want 0 (stderr: %q)", code, stderr.String())
	}
	if want := "It is 70.4°F.\n"; stdout.String() != want {
		t.Errorf("output = %q, want %q", stdout.String(), want)
	}
}

// TestRunUnknownUnitFallsBackToCelsius covers EC-2: an unrecognized --unit value
// falls back to Celsius rather than erroring.
func TestRunUnknownUnitFallsBackToCelsius(t *testing.T) {
	weatherSrv := httptest.NewServer(weatherHandler(http.StatusOK,
		`{"current":{"temperature_2m":21.34},"current_units":{"temperature_2m":"°C"}}`))
	defer weatherSrv.Close()
	swapServers(t, weatherSrv, nil)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--lat", "38.72", "--lon", "-9.14", "--unit", "k"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run() exit code = %d, want 0 (stderr: %q)", code, stderr.String())
	}
	if want := "It is 21.3°C.\n"; stdout.String() != want {
		t.Errorf("output = %q, want %q", stdout.String(), want)
	}
}

// TestRunServiceError covers R4/EC-1: a non-2xx weather response yields a clear
// stderr message, a non-zero exit code, and no output on stdout (no panic).
func TestRunServiceError(t *testing.T) {
	weatherSrv := httptest.NewServer(weatherHandler(http.StatusInternalServerError, ""))
	defer weatherSrv.Close()
	swapServers(t, weatherSrv, nil)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--lat", "38.72", "--lon", "-9.14"}, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("run() exit code = 0, want non-zero on service error")
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want empty on error", stdout.String())
	}
	if !strings.Contains(stderr.String(), "weather:") {
		t.Errorf("stderr = %q, want a weather-prefixed error message", stderr.String())
	}
}

// TestRunPairedCoordinateFlags covers D2: exactly one of --lat/--lon is a user
// error, reported to stderr with a non-zero exit and no geolocation call.
func TestRunPairedCoordinateFlags(t *testing.T) {
	geoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("geo.Locate must not be called when a coordinate flag is set")
	}))
	defer geoSrv.Close()
	swapServers(t, nil, geoSrv)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--lat", "38.72"}, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("run() exit code = 0, want non-zero when only --lat is set")
	}
	if !strings.Contains(stderr.String(), "--lat and --lon") {
		t.Errorf("stderr = %q, want a paired-flag error message", stderr.String())
	}
}
