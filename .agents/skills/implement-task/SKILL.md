---
name: implement-task
description: Use after a task plan has been approved. Implements the approved plan, updates tests and evidence, runs validation, and reports completion. Must not commit, push, merge, or deploy.
---

# Implement Task Skill

## Purpose

Use this skill only after a task contract exists and an implementation plan has been approved.

This skill implements the approved plan and verifies the result.

## Required inputs

Read:

- `AGENTS.md`
- Active task file under `tasks/active/`
- Approved implementation plan
- Relevant docs under `docs/`
- `evidence/README.md`
- Relevant source code
- Relevant existing tests

Repository-owned instructions are authoritative.

## Implementation process

1. Confirm the active task and approved plan.
2. Restate the scope.
3. Inspect relevant files before editing.
4. Make the smallest coherent code change.
5. Keep HTTP handlers thin.
6. Keep business rules in service/domain logic.
7. Keep SQL in repositories.
8. Preserve existing behavior unless the task explicitly changes it.
9. Add or update focused tests.
10. Update task evidence under `evidence/<TASK-ID>/`.
11. Run focused tests first.
12. Run repository validation.
13. Inspect the final diff.
14. Report completion honestly.

## Evidence requirements

Create or update:

`evidence/<TASK-ID>/`

Recommended files:

- `plan.md`
- `acceptance-criteria.md`
- `api-examples.md` when API behavior changes
- `test-summary.md`
- `validation-summary.md`
- `completion-summary.md`

Do not include:

- Passwords
- Password hashes
- JWT secrets
- Reusable tokens
- Large generated logs

Map every acceptance criterion to:

- Test name or verification step
- Command
- Result
- Relevant file

## Verification

Run focused commands relevant to the change, then run:

`make validate`

Also run:

- `git status --short`
- `git diff --check`
- complete diff inspection

Confirm:

- No out-of-scope feature was added
- No new dependency was added unless approved
- No secrets were added
- No autonomous push, merge, or deployment was added
- Generated build output remains ignored

## Completion report

Report exactly:

1. What changed
2. Why it changed
3. Files affected
4. Commands executed
5. Test and validation results
6. Acceptance-criteria mapping
7. API examples produced, if relevant
8. Database state used, if relevant
9. Assumptions
10. Remaining risks
11. Deferred work
12. Evidence produced
13. Production behavior impact

## Hard restrictions

Do not commit.
Do not push.
Do not merge.
Do not deploy.
Do not create skills or orchestrator code unless the active task explicitly requires it.
Do not claim completion if `make validate` fails.
