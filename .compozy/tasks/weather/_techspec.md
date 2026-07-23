# TechSpec: Weather CLI

Companion to `_prd.md`. Test contract in `_tests.md`.

## Overview

Three independent library packages (`weatherapi`, `geo`, `tempfmt`) plus a thin
CLI (`cmd/weather`) that wires them together. Everything is standard-library Go
(`net/http`, `encoding/json`, `context`, `flag`) — no third-party dependencies.
Each library package exposes an injectable HTTP client and base URL so it can be
unit-tested against an `httptest.Server` without real network access.

Module path: `github.com/antoniofulg/ptg-weather-demo`.

## Package contracts

### `weatherapi` — Open-Meteo client

```go
package weatherapi

type Current struct {
    TemperatureC float64 // current temperature in Celsius
    Unit         string  // the unit label returned by the API, e.g. "°C"
}

// Fetch returns the current temperature for a coordinate from Open-Meteo.
// It GETs {BaseURL}/v1/forecast?latitude=..&longitude=..&current=temperature_2m
// (BaseURL defaults to https://api.open-meteo.com) and decodes the JSON
// {"current":{"temperature_2m":<float>},"current_units":{"temperature_2m":"<str>"}}.
// A non-2xx status or a decode error returns a wrapped error; ctx cancellation
// is honored.
func Fetch(ctx context.Context, lat, lon float64) (Current, error)
```

A package-level `var BaseURL = "https://api.open-meteo.com"` and
`var HTTPClient = http.DefaultClient` allow tests to point at an `httptest.Server`.

### `geo` — IP geolocation

```go
package geo

type Location struct {
    Lat, Lon      float64
    City, Country string
}

// Locate resolves the caller's approximate location from their public IP using
// the key-less ip-api.com service: GET {BaseURL}/json (BaseURL defaults to
// http://ip-api.com). It decodes {"status","lat","lon","city","country"} and
// returns an error when status != "success" or the request fails.
func Locate(ctx context.Context) (Location, error)
```

Same injectable `BaseURL` / `HTTPClient` package variables as `weatherapi`.

### `tempfmt` — temperature formatting

```go
package tempfmt

type Unit string

const (
    Celsius    Unit = "c"
    Fahrenheit Unit = "f"
)

// Format renders a Celsius temperature in the requested unit to one decimal
// place with the correct symbol, e.g. Format(21.34, Celsius) == "21.3°C" and
// Format(21.34, Fahrenheit) == "70.4°F". An unknown unit falls back to Celsius.
func Format(tempC float64, unit Unit) string
```

Conversion: `f = c*9/5 + 32`.

### `cmd/weather` — CLI

Wires the three: parse flags (`--unit`, `--lat`, `--lon`); when `--lat`/`--lon`
are both unset, call `geo.Locate`; call `weatherapi.Fetch`; print
`It is <formatted> in <City>, <Country>.` (omit the location clause when
coordinates were supplied explicitly). Any error is printed to stderr and the
process exits non-zero.

## Auxiliary packages

- `cache` — a tiny in-memory TTL cache `Cache[Current]` keyed by `(lat,lon)` so a
  repeated lookup within the TTL skips the network. Independent of the CLI wiring.
- `jsonout` — an optional `--json` output mode that prints the result as JSON;
  shipped with its own tests, then documented. Independent of the CLI wiring.

## Conventions

- Every exported function takes `context.Context` first where it does I/O.
- Errors are wrapped with `fmt.Errorf("...: %w", err)`; no `panic` outside `main`.
- Table-driven tests with `httptest.Server` fakes; no real network in tests.
- `FEATURES.md` at the repo root lists shipped capabilities; a group that adds a
  user-visible capability appends its one-line entry there.
