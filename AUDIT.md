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
| TG-001 | `weatherapi/` — Open-Meteo client | — | TG-002, TG-003 | [#1](https://github.com/antoniofulg/ptg-weather-demo/pull/1) (merged) |
| TG-002 | `geo/` — IP geolocation | — | TG-001, TG-003 | [#2](https://github.com/antoniofulg/ptg-weather-demo/pull/2) (merged) |
| TG-003 | `tempfmt/` — °C/°F formatter | — | TG-001, TG-002 | [#3](https://github.com/antoniofulg/ptg-weather-demo/pull/3) (merged) |
| TG-004 | `cmd/weather` — CLI wiring | **TG-001, TG-002, TG-003** | — (dependent) | [#6](https://github.com/antoniofulg/ptg-weather-demo/pull/6) |
| TG-005 | `cache/` — TTL cache | — | TG-006 | [#4](https://github.com/antoniofulg/ptg-weather-demo/pull/4) |
| TG-006 | `jsonout/` — JSON output | — | TG-005 | [#5](https://github.com/antoniofulg/ptg-weather-demo/pull/5) |

**Independent groups ran concurrently; the dependent group (TG-004) was gated
until its prerequisites were complete.** The three library packages have no edges
between them, so `{TG-001, TG-002, TG-003}` is a valid parallel set. `TG-004`
imports all three, so it depends on them and cannot join that set. `TG-005` and
`TG-006` are independent of everything.

### Execution batches actually run

- **Batch 1** — `--multiple weather/TG-001,weather/TG-002,weather/TG-003 --parallel-task-groups`
  → three isolated worktrees/branches, three commits, three PRs. With the default
  concurrency limit of 2, TG-003 waited in the queue until a slot freed.
- **Batch 2** — `--multiple weather/TG-005,weather/TG-006 --parallel-task-groups`
  → two isolated worktrees/branches. TG-006's two internal tasks ran
  **sequentially** within its worktree (implement, then document).
- **Batch 3** — `--multiple weather/TG-004 --parallel-task-groups`, run only
  **after** TG-001/002/003 were completed (hydrating `_task_groups.md` to `[x]`)
  and their PRs merged, so the dependent group built against its dependencies.

Each branch is `compozy/weather-<NNN>-<brief>-<run>` and contains only its own
group's changes (verified via the three-dot merge-base diff). After merging the
three dependency PRs, `go test ./...` passes across all four packages, and TG-004
produced a working end-to-end CLI.

## How the feature was tested

Every scenario was driven against a **real** daemon and **real** Claude agents —
not the unit-test fakes. Compozy stops at local branches, so the branches were
pushed and PRs opened manually, exactly as the intended workflow prescribes.

| # | Scenario | What it proves | Result |
|---|----------|----------------|--------|
| 1 | Parallel launch `{TG-001,002,003}` | Isolated worktree+branch per group; per-group diffs; checkout untouched | ✅ PRs #1–3, `main` never modified |
| 2 | Parallel launch `{TG-005,006}` | Concurrency again; internal task sequencing | ✅ PRs #4–5, TG-006 two ordered commits |
| 3 | Default concurrency bound | At most 2 groups at once | ✅ TG-003 queued behind TG-001/002 |
| 4 | `--parallel-limit 1` | Limit accepted + applied on the group-parallel kind | ✅ |
| 5 | Enqueued vs parallel | Same targets, no flag → `task_multi_enqueued` (serial, no worktrees) | ✅ vs `task_multi_group_parallel` |
| 6 | No-changes group | Empty branch deleted, "nothing to open" | ✅ (dry-run) |
| 7 | Branch retention | Group with commits keeps its branch, worktree cleaned | ✅ "retained by branch …" |
| 8 | Dependency guard | Dependent group rejected from a parallel set | ✅ TG-004 → `task_group_dependencies_unmet` |
| 9 | **Convergence / hydration (ADR-009)** | Completing 001/002/003 hydrates `[x]` → unblocks TG-004 | ✅ TG-004 then ran, working CLI |
| 10 | Re-launch after completion + `--new` | Completed selection refused; `--new` = fresh namespace | ✅ "already completed … use --new" |
| 11 | Conflict at merge | Two groups editing one file conflict at merge, no pre-flight prevention | ✅ `git merge-tree` exit 1 (PR #5 vs merged #3 on `FEATURES.md`) |
| 12 | Multi-workspace | Daemon runs other workspaces' runs concurrently | ✅ alongside `fitnesshub-web` / `perdura-tg003` reviews |
| 13 | **Fault isolation** | One group's fault never stops its sibling | ✅ TG-006 worktree deleted mid-run → **TG-005 completed unaffected** |
| 14 | Review rounds (all 6 groups) | Per-group review produces artifacts + real findings | ✅ `reviews-001/` per group (see below) |
| — | Plan-drift rejection | Checksum drift between preflight and start → rejected | ⚙️ unit-tested (UT-080); not hand-inducible live |

### Review-round findings (one per group)

| Group | Finding | Severity |
|-------|---------|----------|
| TG-001 | A 2xx response missing `current` yields a silent `0°C` (no field-presence check) | medium |
| TG-002 | Non-2xx error isn't a typed sentinel (inconsistent with `weatherapi`) | low |
| TG-003 | Unit matching is case-sensitive — `--unit F` silently falls back to Celsius | low |
| TG-004 | No timeout on network calls — a hung API blocks the CLI forever | medium |
| TG-005 | Expired cache entries are never evicted (unbounded map growth) | low |
| TG-006 | `Encode` returns an unwrapped, effectively-unreachable error | low |

These are exactly the class of issue a green test suite cannot surface — the code
is correct and tested, but a human/agent review finds real polish and robustness
gaps. That is the point of the per-group review-round step (US-007).

## Findings & observations

- **BUG-CANDIDATE 1 — `runs purge` aborts globally on one orphaned worktree.** The
  home-scoped `compozy runs purge` hard-errors (`has committed changes not retained
  by a branch`) on the first problematic terminal run and aborts the *entire*
  purge, rather than skipping it and continuing. Refusing to delete committed work
  is correct; aborting all purging is the open question. (Surfaced on a
  pre-existing, unrelated run.) Severity: low/medium.
- **BUG-CANDIDATE 2 — a deleted worktree is recovered, not failed.** US-005.EC-2 /
  IT-012 contract: a group whose worktree is deleted mid-run "fails cleanly."
  Observed under a real `rm -rf` against the live daemon: the recovery machinery
  **re-allocated the worktree and the group completed with valid code**. Positive
  robustness, but a genuine behavior-vs-spec divergence the faked-git-runner test
  (IT-012) cannot model. Team triage: update the spec/test to reflect recovery, or
  scope "fails cleanly" to non-recoverable timing. Severity: low (positive), but a
  real mismatch.
- **UX note — `--dry-run` consumes the selection fingerprint.** A `--dry-run`
  creates a completed run record bearing the selection's fingerprint, so a
  subsequent real run of the same set is refused until `--new`.
- **Observation — a stray untracked file preserves a completed worktree.** When an
  agent leaves any untracked file behind, a *completed* group's worktree is
  preserved (not cleaned) — conservative and safe, but the "completed → removed"
  happy path is skipped.

## To reproduce

```bash
git clone https://github.com/antoniofulg/ptg-weather-demo && cd ptg-weather-demo
compozy sync --name weather
compozy tasks validate --name weather
compozy tasks run --multiple weather/TG-001,weather/TG-002,weather/TG-003 --parallel-task-groups --ide claude
git branch --list 'compozy/*'   # one branch per group
```

_All matrix scenarios executed; findings above are open for triage._
