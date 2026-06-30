# Agent Workflow

This repository is intended to be safe for direct agent-assisted development without an external orchestrator.

## Starting Work

1. Read the task file or user request.
2. Check `git status --short`.
3. Inspect the relevant backend, frontend, migration, and docs files before editing.
4. Confirm the change is within current scope.
5. Record assumptions in the active task or final report when they affect validation or behavior.

## Task Lifecycle

Use one Markdown file per task:

- `tasks/inbox/`: proposed work not yet accepted or scoped.
- `tasks/active/`: work in progress.
- `tasks/review/`: implementation complete and awaiting review.
- `tasks/completed/`: reviewed or accepted work.

Task file names should use `YYYY-MM-DD-short-slug.md`.

Move a task file between lifecycle directories instead of keeping duplicate records.

## During Implementation

- Keep edits scoped to the task.
- Follow existing module boundaries.
- Do not add marketplace features outside the accepted task.
- Do not add new dependencies unless the task requires them for deterministic validation.
- Treat database reset commands as destructive local operations.

## Validation Evidence

Every completion report should name:

- commands run;
- whether each command passed, failed, skipped, or was not run;
- database state used for integration validation;
- known gaps and follow-up work.

`make validate` must not modify tracked source files or database data, but it may generate ignored build artifacts such as `frontend/dist`.

Prefer concise command summaries over raw terminal dumps.

## Completion Report

A completion report should include:

- what changed;
- why it changed;
- files affected;
- commands executed;
- test and validation results;
- acceptance-criteria mapping;
- assumptions;
- remaining risks;
- deferred work;
- evidence produced.

Do not claim validation passed unless the command was run for the current change.
