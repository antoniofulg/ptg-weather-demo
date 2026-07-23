---
status: pending
title: Open-Meteo client and tests
type: backend
complexity: medium
---

# Task 1: Open-Meteo client and tests

## Overview
Implement the `weatherapi` package exactly per the TechSpec "weatherapi — Open-Meteo
client" contract. Standard library only.

<critical>
- ALWAYS READ `_prd.md`, `_techspec.md`, and `_tests.md` at the initiative root before starting
- REFERENCE the TechSpec package contract — do not invent a different signature
- Standard library only (`net/http`, `encoding/json`, `context`); no third-party modules
</critical>

<requirements>
- R1: MUST implement `func Fetch(ctx context.Context, lat, lon float64) (Current, error)` returning `Current{TemperatureC float64; Unit string}`.
- R2: MUST GET `{BaseURL}/v1/forecast?latitude=..&longitude=..&current=temperature_2m` and decode `{"current":{"temperature_2m":<float>},"current_units":{"temperature_2m":"<str>"}}`.
- R3: MUST expose package-level `var BaseURL = "https://api.open-meteo.com"` and `var HTTPClient = http.DefaultClient` so tests can inject an `httptest.Server`.
- R4: MUST wrap a non-2xx status and any decode error with `fmt.Errorf("...: %w", err)` and honor `ctx` cancellation. No `panic`.
</requirements>

## Deliverables
- `weatherapi/client.go`
- `weatherapi/client_test.go`
- Every test case assigned in `## Tests` implemented and passing **(REQUIRED)**

## Tests

- Unit
  - [ ] `UT-001` — happy fetch returns the decoded temperature and unit
  - [ ] `UT-002` — a non-200 status returns a non-nil error
  - [ ] `UT-003` — malformed JSON returns a non-nil error

## Success Criteria
- `go build ./...` passes and `go test ./weatherapi/...` passes with the three tests
