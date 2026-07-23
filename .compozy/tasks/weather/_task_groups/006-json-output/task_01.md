---
status: pending
title: JSON encoder for a result
type: backend
complexity: low
---

# Task 1: JSON encoder for a result

## Overview
Implement the `jsonout` package: encode a weather result as JSON so the CLI can
offer a `--json` mode later.

<critical>
- ALWAYS READ `_prd.md`, `_techspec.md`, and `_tests.md` at the initiative root before starting
- Standard library only (`encoding/json`)
- This package is independent of the CLI wiring — do not import `cmd/weather`
</critical>

<requirements>
- R1: MUST implement `type Result struct { TemperatureC float64; Unit, City, Country string }` with JSON tags `temperature_c`, `unit`, `city`, `country`.
- R2: MUST implement `func Encode(r Result) ([]byte, error)` returning the compact JSON object.
</requirements>

## Deliverables
- `jsonout/jsonout.go`
- `jsonout/jsonout_test.go`
- Every test case assigned in `## Tests` implemented and passing **(REQUIRED)**

## Tests

- Unit
  - [ ] `UT-040` — `Encode` emits `{"temperature_c":21.3,"unit":"°C","city":"Lisbon","country":"Portugal"}`

## Success Criteria
- `go build ./...` passes and `go test ./jsonout/...` passes
