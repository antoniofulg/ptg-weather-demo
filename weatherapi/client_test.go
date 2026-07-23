package weatherapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		body     string
		wantErr  bool
		wantTemp float64
		wantUnit string
	}{
		{
			name:     "UT-001 happy path decodes temperature and unit",
			status:   http.StatusOK,
			body:     `{"current":{"temperature_2m":21.3},"current_units":{"temperature_2m":"°C"}}`,
			wantTemp: 21.3,
			wantUnit: "°C",
		},
		{
			name:    "UT-002 non-2xx status returns error",
			status:  http.StatusInternalServerError,
			body:    `{}`,
			wantErr: true,
		},
		{
			name:    "UT-003 malformed JSON returns error",
			status:  http.StatusOK,
			body:    `{not json`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v1/forecast" {
					t.Errorf("request path = %q, want /v1/forecast", r.URL.Path)
				}
				q := r.URL.Query()
				for _, key := range []string{"latitude", "longitude", "current"} {
					if q.Get(key) == "" {
						t.Errorf("missing query param %q in %q", key, r.URL.RawQuery)
					}
				}
				if got := q.Get("current"); got != "temperature_2m" {
					t.Errorf("current = %q, want temperature_2m", got)
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			restore := swapClient(srv.URL, srv.Client())
			defer restore()

			got, err := Fetch(context.Background(), 38.72, -9.14)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Fetch() error = nil, want non-nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Fetch() error = %v, want nil", err)
			}
			if got.TemperatureC != tt.wantTemp {
				t.Errorf("TemperatureC = %v, want %v", got.TemperatureC, tt.wantTemp)
			}
			if got.Unit != tt.wantUnit {
				t.Errorf("Unit = %q, want %q", got.Unit, tt.wantUnit)
			}
		})
	}
}

// TestFetchUnexpectedStatusIsWrapped verifies R4: a non-2xx status yields a
// wrapped, matchable error rather than a bare failure.
func TestFetchUnexpectedStatusIsWrapped(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	restore := swapClient(srv.URL, srv.Client())
	defer restore()

	_, err := Fetch(context.Background(), 0, 0)
	if !errors.Is(err, errUnexpectedStatus) {
		t.Fatalf("error = %v, want wrapped errUnexpectedStatus", err)
	}
}

// TestFetchHonorsContextCancellation verifies R4: a cancelled ctx aborts the call.
func TestFetchHonorsContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	restore := swapClient(srv.URL, srv.Client())
	defer restore()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := Fetch(ctx, 0, 0); err == nil {
		t.Fatal("Fetch() with cancelled ctx error = nil, want non-nil")
	}
}

// swapClient points the package at the test server and returns a restore func.
func swapClient(baseURL string, client *http.Client) func() {
	oldBase, oldClient := BaseURL, HTTPClient
	BaseURL = baseURL
	HTTPClient = client
	return func() {
		BaseURL = oldBase
		HTTPClient = oldClient
	}
}
