# User Stories: Weather CLI

## US-001: See the current temperature for my location

**As a** terminal user, **I want** to run `weather` and see the current
temperature for wherever I am, **so that** I get a fast answer without setup.

- AC-1: `weather` resolves my location by IP and prints the current temperature
  with the city and country.
- AC-2: `weather --unit f` prints the temperature in Fahrenheit.
- AC-3: `weather --lat <x> --lon <y>` prints the temperature for that coordinate
  and omits the location clause.

Edge cases:
- EC-1: The weather or geolocation service returns a non-2xx status → a clear
  error is printed to stderr and the process exits non-zero (no panic).
- EC-2: An unknown `--unit` value falls back to Celsius.

## US-002: Reuse recent lookups

**As a** terminal user, **I want** a repeated lookup for the same coordinate to
avoid a second network call within a short window, **so that** back-to-back runs
are fast.

- AC-1: A cache returns a stored result for the same coordinate within its TTL and
  fetches again after the TTL expires.

## US-003: Machine-readable output

**As a** script author, **I want** `weather --json` to print the result as JSON,
**so that** I can pipe it into other tools.

- AC-1: `--json` prints `{"temperature_c":..,"unit":..,"city":..,"country":..}`.
