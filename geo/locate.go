// Package geo resolves an approximate location from the caller's public IP
// using the key-less ip-api.com service.
package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Location is an approximate geolocation resolved from a public IP address.
type Location struct {
	Lat, Lon      float64
	City, Country string
}

// BaseURL and HTTPClient are package-level knobs so tests can point Locate at an
// httptest.Server instead of the real ip-api.com service. BaseURL defaults to
// the key-less ip-api.com endpoint.
var (
	BaseURL    = "http://ip-api.com"
	HTTPClient = http.DefaultClient
)

// locateResponse mirrors the subset of the ip-api.com /json payload we consume.
type locateResponse struct {
	Status  string  `json:"status"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	City    string  `json:"city"`
	Country string  `json:"country"`
}

// Locate resolves the caller's approximate location from their public IP using
// the key-less ip-api.com service. It GETs {BaseURL}/json and decodes
// {"status","lat","lon","city","country"}. It returns a wrapped error when the
// request fails, the service returns a non-2xx status, the body cannot be
// decoded, or the reported status is not "success". ctx cancellation is honored.
func Locate(ctx context.Context) (Location, error) {
	url := BaseURL + "/json"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Location{}, fmt.Errorf("geo: new request: %w", err)
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return Location{}, fmt.Errorf("geo: request %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Location{}, fmt.Errorf("geo: unexpected status %d from %s", resp.StatusCode, url)
	}

	var body locateResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return Location{}, fmt.Errorf("geo: decode response: %w", err)
	}

	if body.Status != "success" {
		return Location{}, fmt.Errorf("geo: lookup failed with status %q", body.Status)
	}

	return Location{
		Lat:     body.Lat,
		Lon:     body.Lon,
		City:    body.City,
		Country: body.Country,
	}, nil
}
