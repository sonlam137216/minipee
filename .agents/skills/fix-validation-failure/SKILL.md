---
name: fix-validation-failure
description: Use when tests, build, migration, typecheck, or make validate fails during a task. Diagnoses the failure, fixes only the root cause, updates evidence, and reruns validation.
---

# Fix Validation Failure Skill

## Purpose

Use this skill when a validation command fails.

The goal is to fix the root cause without expanding task scope.

## Required inputs

Read:

- The failing command
- Full failure output
- Active task file
- `AGENTS.md`
- Relevant docs
- Relevant source and tests
- Existing task evidence

## Diagnosis process

1. Identify the exact failing command.
2. Identify the first meaningful failure.
3. Determine whether this is:
   - Production bug
   - Test bug
   - Environment/prerequisite issue
   - Flaky or non-isolated test
   - Documentation mismatch
   - Generated artifact issue
4. Check whether the failure is inside the active task scope.
5. Propose the smallest safe fix.
6. Do not rewrite unrelated code.

## Fix rules

If the failure is a test-isolation issue:

- Prefer unique test data
- Prefer cleanup with `t.Cleanup` where practical
- Do not delete unrelated developer data
- Do not require a clean shared database unless explicitly documented

If the failure is a production bug:

- Fix only the defect required by the task
- Add regression coverage
- Preserve public API behavior unless the bug fix requires a documented change

If the failure is an environment issue:

- Improve error clarity if appropriate
- Do not hide or ignore failures
- Do not make validation destructive

## Verification

After fixing:

1. Run the focused failing command.
2. Repeat the focused command when rerunnability matters.
3. Run `make validate`.
4. Run `git diff --check`.
5. Inspect the final diff.

## Evidence

Update current task evidence with:

- Failing command
- Failure summary
- Root cause
- Fix applied
- Commands rerun
- Results
- Remaining risks

## Restrictions

Do not broaden task scope.
Do not reset or destroy the developer database unless the human explicitly approves.
Do not add dependencies unless approved.
Do not commit, push, merge, or deploy.
Do not hide failed checks.
