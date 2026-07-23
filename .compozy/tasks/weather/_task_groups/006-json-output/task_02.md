---
status: pending
title: Document the JSON output
type: docs
complexity: low
---

# Task 2: Document the JSON output

## Overview
Document the JSON output capability now that the `jsonout` encoder exists.

<critical>
- Depends on `task_01` (the `jsonout` package must already exist)
- Documentation only — do not change package behavior
</critical>

<requirements>
- R1: MUST append the line `- JSON output: --json prints the result as a JSON object` to `FEATURES.md` at the repo root (create the file with a `# Features` heading if absent).
- R2: MUST add a short "JSON output" section to `README.md` describing the `{"temperature_c",...}` shape produced by `jsonout.Encode`.
</requirements>

## Deliverables
- `FEATURES.md` entry
- `README.md` "JSON output" section

## Tests

- None (documentation task; no assigned test IDs)

## Success Criteria
- `FEATURES.md` and `README.md` describe the JSON output; `go build ./...` still passes
