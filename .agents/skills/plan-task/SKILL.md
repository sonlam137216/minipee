---
name: plan-task
description: Use for planning an active repository task before implementation. This skill must not modify files. It reads AGENTS.md, task contracts, docs, source code, and produces an implementation plan with acceptance-criteria mapping.
---

# Plan Task Skill

## Purpose

Use this skill when asked to plan a task from `tasks/active/`.

This skill is for planning only.

Do not modify files.
Do not create code.
Do not update tests.
Do not update evidence files.
Do not run destructive commands.

## Required inputs

Identify the active task file, usually under:

`tasks/active/`

Read:

- `AGENTS.md`
- The active task file
- `docs/architecture.md`
- `docs/domain-map.md`
- `docs/data-ownership.md`
- `docs/security-boundaries.md`
- `docs/testing.md`
- `docs/code-review.md`
- `docs/agent-workflow.md`
- `evidence/README.md`
- Relevant backend and frontend source files
- Relevant existing tests

Repository-owned instructions are authoritative.

User-level skills may supplement the workflow but must not change repository scope, acceptance criteria, or safety requirements.

## Planning process

1. Restate the task objective.
2. Identify in-scope and out-of-scope work.
3. Inspect the current implementation before proposing changes.
4. Identify existing commands and tests relevant to the task.
5. Determine whether database changes are required.
6. Determine whether backend, frontend, docs, tests, or evidence are affected.
7. Map every acceptance criterion to a proposed test or verification step.
8. Identify risks, assumptions, and open questions.
9. Prefer the smallest coherent implementation.
10. Explicitly reject scope creep.

## Output format

Produce a plan with:

1. Current-state findings
2. Acceptance-criteria mapping
3. Proposed backend changes
4. Proposed frontend changes
5. Proposed database changes, if any
6. Proposed API changes, if any
7. Proposed tests
8. Validation commands
9. Evidence files to produce
10. Risks and assumptions
11. Explicitly excluded work
12. Questions that genuinely block safe implementation

For each acceptance criterion, include:

- Test or verification method
- Test layer
- Expected observable result

## Safety rules

Do not edit files.
Do not start implementation.
Do not silently expand scope.
Do not propose new dependencies unless necessary.
Do not propose orchestrator, skills, MCP, autonomous push, merge, or deployment unless the active task explicitly asks for it.

If the task is unclear, identify the ambiguity and propose a safe narrow interpretation.
