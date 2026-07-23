---
status: completed
title: IP geolocation and tests
type: backend
complexity: medium
---

# Task 1: IP geolocation and tests

## Overview
Implement the `geo` package exactly per the TechSpec "geo — IP geolocation"
contract. Standard library only.

<critical>
- ALWAYS READ `_prd.md`, `_techspec.md`, and `_tests.md` at the initiative root before starting
- REFERENCE the TechSpec package contract — do not invent a different signature
- Standard library only (`net/http`, `encoding/json`, `context`); no third-party modules
</critical>

<requirements>
- R1: MUST implement `func Locate(ctx context.Context) (Location, error)` returning `Location{Lat, Lon float64; City, Country string}`.
- R2: MUST GET `{BaseURL}/json` and decode `{"status","lat","lon","city","country"}`; `BaseURL` defaults to `http://ip-api.com`.
- R3: MUST return a wrapped error when `status != "success"` or the request fails; honor `ctx`. No `panic`.
- R4: MUST expose package-level `var BaseURL` and `var HTTPClient = http.DefaultClient` for test injection.
</requirements>

## Deliverables
- `geo/locate.go`
- `geo/locate_test.go`
- Every test case assigned in `## Tests` implemented and passing **(REQUIRED)**

## Tests

- Unit
  - [ ] `UT-010` — a `status:"success"` response returns the decoded Location
  - [ ] `UT-011` — a `status:"fail"` response returns a non-nil error

## Success Criteria
- `go build ./...` passes and `go test ./geo/...` passes with the two tests
