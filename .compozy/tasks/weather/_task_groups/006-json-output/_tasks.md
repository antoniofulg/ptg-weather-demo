---
schema_version: "compozy.tasks/v2"
workflow: weather/TG-006
graph:
  nodes:
    - id: task_01
      file: task_01.md
    - id: task_02
      file: task_02.md
  edges:
    - from: task_01
      to: task_02
---

# TG-006 JSON output — Task List

Dependency edges are owned here: `task_01` must finish before `task_02`.

## Tasks

- `task_01` — JSON encoder for a result
- `task_02` — Document the JSON output
