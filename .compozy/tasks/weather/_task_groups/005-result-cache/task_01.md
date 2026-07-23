---
status: pending
title: In-memory TTL cache
type: backend
complexity: medium
---

# Task 1: In-memory TTL cache

## Overview
Implement the `cache` package per the TechSpec "cache" note: a tiny in-memory TTL
cache keyed by coordinate so a repeated lookup within the TTL skips the network.

<critical>
- ALWAYS READ `_prd.md`, `_techspec.md`, and `_tests.md` at the initiative root before starting
- Standard library only (`time`, `sync`); no third-party modules
- This package is independent of the CLI wiring — do not import `cmd/weather`
</critical>

<requirements>
- R1: MUST implement a generic `type Cache[T any]` with `New[T](ttl time.Duration) *Cache[T]`.
- R2: MUST implement `Get(key string) (T, bool)` and `Set(key string, value T)`; `Get` returns `false` after the TTL has elapsed for that key.
- R3: MUST be safe for concurrent use (guard with a mutex). No `panic`.
- R4: Tests MUST control time deterministically (inject a clock or a small elapsed sleep with a short TTL) — no flaky wall-clock waits over ~50ms.
</requirements>

## Deliverables
- `cache/cache.go`
- `cache/cache_test.go`
- Every test case assigned in `## Tests` implemented and passing **(REQUIRED)**

## Tests

- Unit
  - [ ] `UT-030` — a value stored is returned by `Get` within the TTL
  - [ ] `UT-031` — `Get` reports a miss after the TTL has elapsed

## Success Criteria
- `go build ./...` passes and `go test -race ./cache/...` passes with the two tests
