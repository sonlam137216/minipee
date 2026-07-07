---
name: review-task
description: Use for independent review of a completed task implementation. Must not modify files. Reviews diff against main, task contract, docs, tests, evidence, validation, security, scope, and readiness.
---

# Review Task Skill

## Purpose

Use this skill to independently review a task implementation before human commit or merge.

This skill must not modify files.

## Required inputs

Read:

- `AGENTS.md`
- Active task file
- Relevant docs under `docs/`
- `docs/code-review.md`
- `evidence/README.md`
- Task evidence under `evidence/<TASK-ID>/`
- Complete diff against `main`
- Relevant source and tests

## Review process

1. List active instruction sources.
2. Confirm repository-owned instructions are authoritative.
3. Review task scope.
4. Review acceptance-criteria coverage.
5. Review code correctness.
6. Review backend architecture.
7. Review frontend behavior, if applicable.
8. Review database and migrations, if applicable.
9. Review authentication and authorization, if applicable.
10. Review security and sensitive data exposure.
11. Review test quality.
12. Review evidence quality.
13. Run non-destructive focused tests.
14. Run `make validate`.
15. Inspect `git status --short`.
16. Inspect `git diff --check`.
17. Report findings by severity.

## Severity levels

Use:

- Blocker: unsafe or unusable; must fix before commit or merge
- High: likely correctness, security, validation, or scope failure
- Medium: meaningful reliability or maintainability issue
- Low: optional improvement or polish

For each finding include:

- File and location
- Why it matters
- How it can fail
- Minimal recommended fix

## Review checks

Confirm:

- The implementation satisfies the task contract
- Tests prove the behavior and do not pass for the wrong reason
- Evidence maps acceptance criteria to results
- `make validate` passes
- No out-of-scope feature was added
- No secrets or reusable tokens were added
- No autonomous push, merge, or deployment was introduced
- Generated output remains ignored
- Remaining risks are stated honestly

## Output format

Report:

1. Active instruction sources
2. Files reviewed
3. Commands executed
4. Result of each command
5. Acceptance-criteria coverage
6. Correctness review
7. Security review
8. Test quality
9. Evidence quality
10. Scope compliance
11. Blocker findings
12. High findings
13. Medium and Low findings
14. Readiness for human commit and merge

## Restrictions

Do not modify files.
Do not fix issues in this review pass.
Do not commit.
Do not push.
Do not merge.
Do not deploy.
