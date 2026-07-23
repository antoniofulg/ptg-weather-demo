// Package weatherapi fetches the current temperature for a coordinate from the
// free, key-less Open-Meteo API.
package weatherapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// BaseURL is the Open-Meteo API root. Tests point it at an httptest.Server.
var BaseURL = "https://api.open-meteo.com"

// HTTPClient performs the request. Tests may replace it with the server client.
var HTTPClient = http.DefaultClient

// Current is the current temperature at a coordinate.
type Current struct {
	TemperatureC float64 // current temperature in Celsius
	Unit         string  // the unit label returned by the API, e.g. "°C"
}

// errUnexpectedStatus is wrapped when Open-Meteo returns a non-2xx status.
var errUnexpectedStatus = errors.New("unexpected status")

// forecastResponse mirrors the subset of the Open-Meteo forecast payload we read.
type forecastResponse struct {
	Current struct {
		Temperature2m float64 `json:"temperature_2m"`
	} `json:"current"`
	CurrentUnits struct {
		Temperature2m string `json:"temperature_2m"`
	} `json:"current_units"`
}

// Fetch returns the current temperature for the given coordinate from Open-Meteo.
// It GETs {BaseURL}/v1/forecast?latitude=..&longitude=..&current=temperature_2m
// and decodes the JSON current temperature and its unit label. A non-2xx status
// or a decode failure returns a wrapped error; ctx cancellation is honored.
func Fetch(ctx context.Context, lat, lon float64) (Current, error) {
	q := url.Values{}
	q.Set("latitude", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Set("longitude", strconv.FormatFloat(lon, 'f', -1, 64))
	q.Set("current", "temperature_2m")
	endpoint := BaseURL + "/v1/forecast?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return Current{}, fmt.Errorf("weatherapi: build request: %w", err)
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return Current{}, fmt.Errorf("weatherapi: fetch temperature: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return Current{}, fmt.Errorf("weatherapi: fetch temperature: %w: %s", errUnexpectedStatus, resp.Status)
	}

	var body forecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return Current{}, fmt.Errorf("weatherapi: decode response: %w", err)
	}

	return Current{
		TemperatureC: body.Current.Temperature2m,
		Unit:         body.CurrentUnits.Temperature2m,
	}, nil
}
