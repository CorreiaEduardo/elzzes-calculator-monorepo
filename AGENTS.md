# Project instructions

## What this project is

Elzzes Calculator is a monorepo with three packages:

- `elzzes-calculator-api` — a stateless Go REST microservice. It owns all arithmetic and the OpenAPI contract.
- `elzzes-calculator-web` — the React + TypeScript SPA. It owns input, UI state, formatting, and error presentation; it uses TanStack Query and Zod at the API boundary.
- `elzzes-ui-lib` — internal presentational components, design tokens, and Storybook stories, published as `@elzzes/ui`. No state, network, or arithmetic.

The agreed stack is pnpm workspaces; React 19 + TypeScript + Vite + Tailwind CSS v4 + TanStack Query + Zod; Go standard library; OpenAPI 3.1; Nginx + Docker Compose; GitHub Actions; Playwright optional.

`problem.md` is the lighthouse: every requirement and constraint there must stay satisfied. Canonical sources, one per topic:

- wire contract — [`elzzes-calculator-api/docs/openapi.yaml`](elzzes-calculator-api/docs/openapi.yaml)
- operation behavior and edge cases — [specs 001](docs/specs/001-basic-arithmetic-operations/SPEC.md) and [002](docs/specs/002-advanced-arithmetic-operations/SPEC.md)
- UI behavior — [`docs/DESIGN.md`](docs/DESIGN.md)
- system structure, request flow, numeric model, deployment — [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md)
- implementation detail — component [backend](elzzes-calculator-api/docs/CONVENTIONS.md), [frontend](elzzes-calculator-web/docs/CONVENTIONS.md), [UI](elzzes-ui-lib/docs/CONVENTIONS.md) conventions
- invariants — [`docs/PRINCIPLES.md`](docs/PRINCIPLES.md)
- verification, coverage, CI — [spec 004](docs/specs/004-delivery-and-verification/SPEC.md)
- high-level decisions and rejected alternatives — [`docs/DECISIONS.md`](docs/DECISIONS.md)

Never duplicate a definition. Update its canonical document and link from everywhere else. If docs conflict, reconcile before coding.

**Current phase:** specs 001–003 implemented (2026-11-20): the Go API, calculator UI, API client, reducer, tests, and Storybook. Next is spec 004 (delivery and verification).

## Priorities

Correctness, clarity, and maintainability over extra features. Clean, readable, idiomatic code in both layers. The brief's scope is the feature ceiling. The full order of concerns and decision rules are in [`docs/PRINCIPLES.md`](docs/PRINCIPLES.md).

## Non-negotiables

- Arithmetic has exactly one implementation, in the Go service. No client-side re-implementation, no `eval`.
- The wire contract — operation vocabulary, arity, request/response shape, error codes — changes only through a spec/OpenAPI/`ARCHITECTURE.md` update agreed with the user.
- Every documented edge case has a test; no `NaN`/`Infinity` reaches the UI.
- No dependency, abstraction, or process beyond the agreed stack without a requirement and user agreement.
- Accessibility and responsive behavior are acceptance criteria, not polish.
- No secrets in the repository.

## How we work in this repo

- Changes land one at a time: implement → tests → typecheck/vet/lint → verify.
- Every distinct change gets its own commit, and notable ones get a `docs/CHANGELOG.md` entry the same day — `git log` plus the changelog read as the story of the app, which keeps rollbacks surgical.
- Change architecture (agreed with the user FIRST — structural calls are never made unilaterally) → update the implementation → update the canonical docs (`ARCHITECTURE.md`, `openapi.yaml`, specs) → verify.
- Check a planned change against `docs/PRINCIPLES.md`; its invariants are standing rules.
- Comment budget: lean. Default is no comment — naming and structure carry the what. Comment only a constraint the code can't show (a native/timing race, an ordering requirement, a standing invariant).
- Update docs in the same change: behavior → specs/DESIGN; contract → OpenAPI/ARCHITECTURE; setup/API examples → README; implementation detail → component conventions.
- Verification is canonical in [spec 004](docs/specs/004-delivery-and-verification/SPEC.md). Run the local checks before committing; CI mirrors them.
- Definition of done: tests/typecheck/vet/lint/coverage/CI green, docs current, change revertable as one unit ([PRINCIPLES § What done means](docs/PRINCIPLES.md#what-done-means) and spec 004).

## Project documentation

The full map is in [README § Documentation](README.md#documentation). Agent-specific rules:

- `problem.md` wins over every other document.
- `brainstorm.md` is raw early notes, not authoritative; `job.md` is not project documentation.
- `docs/DECISIONS.md` is a high-level decision index; it must link canonical detail rather than restate it.
- Do not add a parallel contract, configuration table, or UI rule to a second document.
