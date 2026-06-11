# CLAUDE.md — MediCabinet

## Project

MediCabinet is a drug inventory tracker for households: a web app with a single backend server. Read `README.md` for the architecture before making changes.

## Architecture Rules

- Respect the subsystem decomposition: Inventory, Medication Catalog, Interaction Checker, Reminder & Alerts, Webapp.
- The webapp talks to backend services only — never to another subsystem's internals or database.
- Safety-critical logic (contraindication checks) lives in the backend, never only in the UI.
- External providers (drug database, notifications) are accessed only through their adapter component.

## Workflow

- Work on **one user story at a time**, end to end.
- **Plan first**: present a short step-by-step plan and wait for approval before editing code.
- Make small changes — no big rewrites across many files at once.

## Definition of Done

A task is only done when ALL of these hold:

1. **Review your own code** after you think you're done: re-read the diff and check it against the acceptance criteria.
2. **Write test cases** that verify the behavior of the user story (not just the happy path).
3. **Run the full test suite** — every test must pass. Never claim something works without executing it.
4. Explain in 2–3 sentences what you changed and why.

## Never

- Never commit secrets, API keys, or real personal/health data.
- Never delete or skip failing tests to make the suite green.
- Never mark a task done with failing or unexecuted tests.
