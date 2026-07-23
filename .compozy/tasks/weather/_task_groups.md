---
schema_version: compozy.task-groups/v1
initiative: weather
graph:
  nodes:
    - id: TG-001
      directory: _task_groups/001-weatherapi-client
    - id: TG-002
      directory: _task_groups/002-ip-geolocation
    - id: TG-003
      directory: _task_groups/003-temperature-format
    - id: TG-004
      directory: _task_groups/004-cli-delivery
    - id: TG-005
      directory: _task_groups/005-result-cache
    - id: TG-006
      directory: _task_groups/006-json-output
  edges:
    - from: TG-001
      to: TG-004
      rationale: The CLI wires the Open-Meteo client, so the client must land first
    - from: TG-002
      to: TG-004
      rationale: The CLI resolves location before fetching, so geolocation must land first
    - from: TG-003
      to: TG-004
      rationale: The CLI formats the temperature for output, so the formatter must land first
---

# weather Task Groups

## [x] TG-001 — Open-Meteo client

- Reference: `weather/TG-001`
- Outcome: The `weatherapi` package fetches the current temperature for a coordinate
- Owns:
  - The `weatherapi` package
  - Open-Meteo request and JSON decoding
- Dependencies: None

## [x] TG-002 — IP geolocation

- Reference: `weather/TG-002`
- Outcome: The `geo` package resolves an approximate location from the caller's IP
- Owns:
  - The `geo` package
  - ip-api.com request and JSON decoding
- Dependencies: None

## [x] TG-003 — Temperature formatting

- Reference: `weather/TG-003`
- Outcome: The `tempfmt` package renders Celsius or Fahrenheit output
- Owns:
  - The `tempfmt` package
  - Celsius/Fahrenheit conversion and rounding
- Dependencies: None

## [ ] TG-004 — CLI delivery

- Reference: `weather/TG-004`
- Outcome: The `weather` CLI wires geolocation, fetch, and formatting into one command
- Owns:
  - The `cmd/weather` command
  - Flag parsing and output assembly
- Dependencies:
  - `TG-001` — The CLI wires the Open-Meteo client, so the client must land first
  - `TG-002` — The CLI resolves location before fetching, so geolocation must land first
  - `TG-003` — The CLI formats the temperature for output, so the formatter must land first

## [ ] TG-005 — Result cache

- Reference: `weather/TG-005`
- Outcome: The `cache` package memoizes a lookup for a coordinate within a TTL
- Owns:
  - The `cache` package
  - TTL expiry behavior
- Dependencies: None

## [ ] TG-006 — JSON output

- Reference: `weather/TG-006`
- Outcome: The `jsonout` package renders a result as JSON and documents the flag
- Owns:
  - The `jsonout` package
  - The FEATURES.md entry for JSON output
- Dependencies: None
