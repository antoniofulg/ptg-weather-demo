---
status: pending
title: Wire the weather CLI
type: backend
complexity: medium
---

# Task 1: Wire the weather CLI

## Overview
Implement `cmd/weather` per the TechSpec "cmd/weather — CLI" contract, wiring the
`geo`, `weatherapi`, and `tempfmt` packages that already exist in the module.

<critical>
- ALWAYS READ `_prd.md`, `_techspec.md`, and `_tests.md` at the initiative root before starting
- The `weatherapi`, `geo`, and `tempfmt` packages already exist — import them, do not reimplement them
- Standard library only (`flag`, `context`, `fmt`, `os`)
</critical>

<requirements>
- R1: MUST parse flags `--unit` (default `c`), `--lat`, `--lon`.
- R2: When `--lat`/`--lon` are both unset, MUST call `geo.Locate`; otherwise use the supplied coordinate and omit the location clause.
- R3: MUST call `weatherapi.Fetch`, format via `tempfmt.Format`, and print `It is <formatted> in <City>, <Country>.` (or `It is <formatted>.` for an explicit coordinate).
- R4: MUST print any error to stderr and exit non-zero; no `panic`.
</requirements>

## Deliverables
- `cmd/weather/main.go`
- Every test case assigned in `## Tests` implemented and passing **(REQUIRED)**

## Tests

- End-to-end
  - [ ] `E2E-001` — with faked `geo`/`weatherapi` base URLs, `--lat 38.72 --lon -9.14` prints a line containing the formatted temperature

## Success Criteria
- `go build ./...` passes; `go vet ./...` is clean; the CLI runs end to end
