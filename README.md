# ptg-weather-demo

A small **weather CLI** — resolve the current location by IP and print the
current temperature from the free, key-less [Open-Meteo](https://open-meteo.com)
API.

```bash
weather                       # temperature for your IP-resolved location
weather --unit f              # in Fahrenheit
weather --lat 38.72 --lon -9.14
```

## Built with Compozy Parallel Task Groups

This repository is a **dogfooding audit fixture** for Compozy's *Parallel Task
Groups* feature: dependency-independent task groups are executed **concurrently
by real Claude agents**, each isolated in its own Git worktree on its own branch,
landing as its own small pull request. The full scenario log, evidence, and
findings are in **[AUDIT.md](AUDIT.md)**; the plan lives in
[`.compozy/tasks/weather/`](.compozy/tasks/weather).

### Which groups ran in parallel

The weather CLI decomposes into six task groups. The three library packages are
mutually independent and ran **at the same time**; the CLI wiring depends on all
three and is gated until they complete.

```
        ┌── TG-001  weatherapi/  (Open-Meteo client)   ─┐   PR #1
Batch 1 ├── TG-002  geo/         (IP geolocation)       ├─►  TG-004  cmd/weather
        └── TG-003  tempfmt/     (°C/°F formatter)      ─┘   (depends on 001+002+003)

Batch 2 ┌── TG-005  cache/       (TTL cache)            PR #4
        └── TG-006  jsonout/     (JSON output)          PR #5
```

- **Independent → parallel:** `{TG-001, TG-002, TG-003}` have no dependency edges
  between them, so they ran concurrently in three isolated worktrees/branches
  (Batch 1). `{TG-005, TG-006}` likewise (Batch 2).
- **Dependent → gated:** `TG-004` (`cmd/weather`) imports all three libraries, so
  it depends on them and is **rejected from any parallel set** until they are
  complete — proving the dependency guard.
- **One PR per group:** each branch (`compozy/weather-<NNN>-<brief>-<run>`)
  contains only that group's changes, so it reviews as a small, self-contained PR
  instead of one giant branch.

### How the feature was tested

Every scenario was driven against a **real** home-scoped daemon with the **real**
Claude runtime — not the unit-test fakes. Highlights (full matrix + evidence in
[AUDIT.md](AUDIT.md)):

- **Isolation** — each group only touches its own package; your checked-out
  branch and working tree are never modified during a run.
- **Bounded concurrency** — the default limit of 2 was observed (TG-003 queued
  behind TG-001/002).
- **Settlement** — a group with commits keeps its branch and cleans its worktree;
  a no-changes group deletes its empty branch ("nothing to open").
- **Internal sequencing** — TG-006's two internal tasks ran in order (implement,
  then document) within its single worktree.
- **Dependency guard & re-launch safety** — dependent groups are rejected from a
  parallel set; re-issuing a completed selection is refused until `--new`.
- **Conflict at merge** — PR #3 and PR #5 both touch `FEATURES.md`, so the second
  merge surfaces a normal Git conflict (no pre-flight prevention, by design).

## Layout

```
weatherapi/   Open-Meteo client        (TG-001)
geo/          IP geolocation           (TG-002)
tempfmt/      temperature formatting   (TG-003)
cmd/weather/  the CLI                  (TG-004)
cache/        in-memory TTL cache      (TG-005)
jsonout/      JSON output              (TG-006)
```
