# Evidence

Evidence records what was verified for a task. Prefer concise summaries in task files or completion reports instead of committing raw logs.

## Evidence Standard

Task evidence must:

- map every acceptance criterion to observable evidence;
- record the exact commands executed;
- record the result of each command as pass, fail, skipped, or not run;
- include failed and skipped checks honestly, with short failure excerpts when helpful;
- explain the database state used for integration testing, including whether PostgreSQL was running, migrated, or reset;
- prohibit passwords, password hashes, JWT secrets, access tokens, and other secrets;
- avoid committing unnecessarily large generated logs;
- include API request and response examples when API behavior changes;
- include screenshots only when visual behavior is part of the task;
- include a final review report;
- include a completion summary.

Do not report a command as passing unless it was run against the current change.

## Recommended Layout

Use a task-specific directory when evidence needs more than a short completion report:

```text
evidence/<TASK-ID>/
  plan.md
  acceptance-criteria.md
  validation-summary.md
  test-summary.md
  api-examples.md
  review-report.md
  completion-summary.md
```

Do not create irrelevant evidence files. Small documentation-only tasks may keep evidence in the task file or completion report.

## Raw Logs

Avoid committing raw terminal output. If a raw log is necessary for review, keep it short or place it under `evidence/raw/`; raw logs in that directory are ignored by git.
