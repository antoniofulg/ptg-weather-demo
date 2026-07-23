# PRD: Weather CLI

## Problem

A developer wants a one-command way to see the current temperature for wherever
they are, without signing up for an API key or configuring anything.

## Goal

Ship a small Go CLI, `weather`, that resolves the caller's approximate location
by IP and prints the current temperature from a free, key-less weather API.

## Users

- **Terminal user** — runs `weather` and reads the current temperature for their
  location, optionally overriding the location or the temperature unit.

## Requirements

- P1: Resolve approximate location (latitude, longitude, city, country) from the
  caller's public IP, using a key-less service.
- P2: Fetch the current temperature for a latitude/longitude from a key-less
  weather API (Open-Meteo).
- P3: Print a human-readable line, e.g. `It is 21.3°C in Lisbon, Portugal.`
- P4: Support `--unit c|f` (default `c`) and `--lat`/`--lon` to override
  auto-location.
- P5: Fail with a clear, non-panicking error message when the network is
  unavailable or a service returns an error.

## Non-goals

- No API keys, accounts, or persisted configuration.
- No forecast, history, or GUI — current temperature only.

## Success

`weather` prints the correct current temperature for the auto-detected location,
and `weather --lat 38.72 --lon -9.14 --unit f` prints it for an explicit location
in Fahrenheit. Every module is unit-tested against faked HTTP responses.
