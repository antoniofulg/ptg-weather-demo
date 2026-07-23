---
status: completed
title: Temperature formatting and tests
type: backend
complexity: low
---

# Task 1: Temperature formatting and tests

## Overview
Implement the `tempfmt` package exactly per the TechSpec "tempfmt — temperature
formatting" contract. Also record the capability in `FEATURES.md`.

<critical>
- ALWAYS READ `_prd.md`, `_techspec.md`, and `_tests.md` at the initiative root before starting
- REFERENCE the TechSpec package contract — do not invent a different signature
- Standard library only
</critical>

<requirements>
- R1: MUST implement `type Unit string` with `Celsius Unit = "c"` and `Fahrenheit Unit = "f"`.
- R2: MUST implement `func Format(tempC float64, unit Unit) string` rendering one decimal place with the correct symbol (`21.3°C`, `70.4°F`); Fahrenheit uses `f = c*9/5 + 32`.
- R3: MUST fall back to Celsius for any unknown unit.
- R4: MUST append the line `- Temperature units: Celsius and Fahrenheit output` to `FEATURES.md` at the repo root (create the file with a `# Features` heading if absent).
</requirements>

## Deliverables
- `tempfmt/format.go`
- `tempfmt/format_test.go`
- `FEATURES.md` entry
- Every test case assigned in `## Tests` implemented and passing **(REQUIRED)**

## Tests

- Unit
  - [ ] `UT-020` — `Format(21.34, Celsius) == "21.3°C"`
  - [ ] `UT-021` — `Format(21.34, Fahrenheit) == "70.4°F"`
  - [ ] `UT-022` — an unknown unit falls back to Celsius

## Success Criteria
- `go build ./...` passes and `go test ./tempfmt/...` passes with the three tests
