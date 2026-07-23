# Test Specification: Weather CLI

Table-driven Go tests with `httptest.Server` fakes; no real network access.

## weatherapi (TG-001)

- **UT-001** (happy): `Fetch` against a fake server returning
  `{"current":{"temperature_2m":21.3},"current_units":{"temperature_2m":"°C"}}`
  returns `Current{TemperatureC: 21.3, Unit: "°C"}`.
- **UT-002** (error): a fake server returning HTTP 500 → `Fetch` returns a
  non-nil wrapped error.
- **UT-003** (error): a fake server returning malformed JSON → `Fetch` returns a
  non-nil error.

## geo (TG-002)

- **UT-010** (happy): a fake server returning
  `{"status":"success","lat":38.72,"lon":-9.14,"city":"Lisbon","country":"Portugal"}`
  → `Locate` returns that `Location`.
- **UT-011** (error): a fake server returning `{"status":"fail"}` → `Locate`
  returns a non-nil error.

## tempfmt (TG-003)

- **UT-020** (happy): `Format(21.34, Celsius) == "21.3°C"`.
- **UT-021** (happy): `Format(21.34, Fahrenheit) == "70.4°F"` (conversion + round).
- **UT-022** (boundary): `Format(21.34, Unit("k")) == "21.3°C"` (unknown → Celsius).

## cache (TG-005)

- **UT-030** (happy): a value stored in the cache is returned on a lookup within
  the TTL.
- **UT-031** (state): a lookup after the TTL has elapsed reports a miss.

## jsonout (TG-006)

- **UT-040** (happy): the JSON encoder emits
  `{"temperature_c":21.3,"unit":"°C","city":"Lisbon","country":"Portugal"}` for a
  known result.

## cmd/weather (TG-004)

- **E2E-001** (happy): with faked `geo` and `weatherapi` servers, running the CLI
  with `--lat 38.72 --lon -9.14` prints a line containing the formatted
  temperature.
