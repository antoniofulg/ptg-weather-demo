# Audit: dogfooding Compozy Parallel Task Groups

This repository is a live audit fixture for Compozy's **Parallel Task Groups**
feature. A real weather CLI was decomposed into dependency-independent task
groups and executed **concurrently by real Claude agents**, each group isolated
in its own Git worktree on its own branch, landing as its own small pull request.

- Compozy build under test: `v0.2.15-97-g3320524`
- Runtime: Claude (ACP), home-scoped daemon
- Weather API: Open-Meteo (key-less); geolocation: ip-api.com (key-less)

## The plan: what runs in parallel

The initiative plan lives in [`.compozy/tasks/weather/`](.compozy/tasks/weather)
(PRD, TechSpec, user stories, test contract, ADR, and the `_task_groups.md`
dependency graph). Six task groups:

| Group | Module | Depends on | Runs in parallel with | PR |
|-------|--------|-----------|-----------------------|----|
| TG-001 | `weatherapi/` — Open-Meteo client | — | TG-002, TG-003 | [#1](https://github.com/antoniofulg/ptg-weather-demo/pull/1) |
| TG-002 | `geo/` — IP geolocation | — | TG-001, TG-003 | [#2](https://github.com/antoniofulg/ptg-weather-demo/pull/2) |
| TG-003 | `tempfmt/` — °C/°F formatter | — | TG-001, TG-002 | [#3](https://github.com/antoniofulg/ptg-weather-demo/pull/3) |
| TG-004 | `cmd/weather` — CLI wiring | **TG-001, TG-002, TG-003** | — (dependent) | pending |
| TG-005 | `cache/` — TTL cache | — | TG-006 | [#4](https://github.com/antoniofulg/ptg-weather-demo/pull/4) |
| TG-006 | `jsonout/` — JSON output | — | TG-005 | [#5](https://github.com/antoniofulg/ptg-weather-demo/pull/5) |

**Independent groups run concurrently; the dependent group (TG-004) is gated
until its prerequisites are complete.** The three library packages have no edges
between them, so `{TG-001, TG-002, TG-003}` is a valid parallel set. `TG-004`
imports all three, so it depends on them and cannot join that set. `TG-005` and
`TG-006` are independent of everything.

### Execution batches actually run

- **Batch 1** — `tasks run --multiple weather/TG-001,weather/TG-002,weather/TG-003 --parallel-task-groups`
  → three isolated worktrees/branches, three commits, three PRs. With the default
  concurrency limit of 2, TG-003 waited in the queue until a slot freed.
- **Batch 2** — `tasks run --multiple weather/TG-005,weather/TG-006 --parallel-task-groups`
  → two isolated worktrees/branches. TG-006's two internal tasks ran
  **sequentially** within its worktree (implement, then document).

Each branch is named `compozy/weather-<NNN>-<brief>-<run>` and contains only its
own group's changes — verified against the merge base:

```
TG-001  3dc242f  weatherapi/client.go + client_test.go        (only weatherapi/**)
TG-002  507af29  geo/locate.go + locate_test.go               (only geo/**)
TG-003  e962f38  tempfmt/format.go + format_test.go + FEATURES.md
TG-005  25351d4  cache/cache.go + cache_test.go               (only cache/**)
TG-006  aba380f  jsonout encoder  →  4b2ed27  jsonout docs     (two sequential commits)
```

## How the feature was tested

Every scenario below was driven against a **real** daemon and **real** Claude
agents — not the unit-test fakes. Compozy stops at local branches (it never
pushes, opens PRs, or merges), so the branches were pushed and PRs opened
manually, exactly as the intended workflow prescribes.

| # | Scenario | What it proves | Result |
|---|----------|----------------|--------|
| 1 | Parallel launch of `{TG-001,002,003}` | One isolated worktree+branch per group; per-group diffs; user checkout untouched | ✅ PRs #1–3, `main` never modified |
| 2 | Parallel launch of `{TG-005,006}` | Concurrency again; internal task sequencing (TG-006) | ✅ PRs #4–5, two ordered commits on TG-006 |
| 3 | Default concurrency bound | At most 2 groups run at once | ✅ TG-003 queued behind TG-001/002 |
| 4 | No-changes group | Empty branch deleted, "nothing to open" | ✅ observed on a dry-run |
| 5 | Branch retention | Group with commits keeps its branch, worktree cleaned | ✅ "retained by branch …" |
| 6 | Dependency guard | A dependent group cannot join a parallel set | ✅ TG-004 → `task_group_dependencies_unmet` |
| 7 | Re-launch after completion | Re-issuing a completed selection is refused | ✅ "already completed … use --new" |
| 8 | Fresh namespace (`--new`) | A fresh run/branch namespace, prior branches intact | ✅ used to start Batch 1 cleanly |
| 9 | Conflict at merge | Two groups editing the same file conflict at PR merge, no pre-flight prevention | ⏳ set up: PR #3 and PR #5 both touch `FEATURES.md` |
| 10 | Multi-workspace | Daemon runs other workspaces' runs concurrently | ✅ ran alongside two unrelated review rounds |
| — | Dependency **convergence** | Completing TG-001/002/003 hydrates `_task_groups.md`, unblocking TG-004 | ⏳ in progress |

### To reproduce

```bash
git clone https://github.com/antoniofulg/ptg-weather-demo && cd ptg-weather-demo
compozy sync --name weather
compozy tasks validate --name weather
compozy tasks run --multiple weather/TG-001,weather/TG-002,weather/TG-003 --parallel-task-groups --ide claude
git branch --list 'compozy/*'   # one branch per group
```

## Findings

- **BUG-CANDIDATE — `runs purge` aborts globally on one orphaned worktree.** The
  home-scoped `compozy runs purge` hard-errors (`has committed changes not
  retained by a branch`) on the first problematic terminal run and aborts the
  entire purge, rather than skipping that run and continuing. Refusing to delete
  committed-but-unretained work is correct; aborting *all* purging is the open
  question. (Surfaced on a pre-existing, unrelated run — not caused by this
  feature.) Severity: low/medium (robustness).
- **UX note — `--dry-run` consumes the selection fingerprint.** A `--dry-run` of
  a selection creates a completed run record bearing that selection's fingerprint,
  so a subsequent *real* run of the same set is refused until `--new`. Minor
  gotcha worth documenting.

_This document is regenerated as scenarios complete._
