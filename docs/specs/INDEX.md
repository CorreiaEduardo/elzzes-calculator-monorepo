# Specs

`problem.md` is the source of truth for the project. These specs break it into reviewable slices and record the decisions behind them. When a spec and the brief disagree, the brief wins; when a spec and [`ARCHITECTURE.md`](../ARCHITECTURE.md) or [`openapi.yaml`](../../elzzes-calculator-api/docs/openapi.yaml) disagree, stop and reconcile before building.

**Status:** implementation baseline, 2026-11-20. Specs 001–003 are implemented with automated checks green; spec 003's manual responsive/screen-reader pass remains. Spec 004 remains scoped.

## Specs

| # | Spec | Scope | Depends on | Status |
| --- | --- | --- | --- | --- |
| 001 | [Basic arithmetic operations](001-basic-arithmetic-operations/SPEC.md) | `add`, `subtract`, `multiply`, `divide`; shared request validation; basic keypad | — | Implemented |
| 002 | [Advanced arithmetic operations](002-advanced-arithmetic-operations/SPEC.md) | `power`, `sqrt`, `percent`; unary/binary domains; advanced keys | 001 | Implemented |
| 003 | [Calculator UI](003-calculator-ui/SPEC.md) | Entry model, formatting, TanStack Query + Zod API client, errors, keyboard, responsive, a11y; `@elzzes/ui` split; Storybook | 001, 002 | Implemented |
| 004 | [Delivery and verification](004-delivery-and-verification/SPEC.md) | Compose/Nginx images, tests/coverage, CI, README/docs completeness, reviewer checklist | 001–003 | Scoped |

## Traceability to the brief

| `problem.md` requirement | Addressed by |
| --- | --- |
| Addition, subtraction, multiplication, division | 001 |
| Optional: exponentiation, square root, percentage | 002 |
| React UI, intuitive input/display, validation, responsive design | 003; [Design](../DESIGN.md) |
| REST API, input validation, edge cases, JSON | 001, 002; [`openapi.yaml`](../../elzzes-calculator-api/docs/openapi.yaml); [Architecture § API contract](../ARCHITECTURE.md#api-contract) |
| Clean, idiomatic code; unit tests for both layers | All specs; [Principles](../PRINCIPLES.md); 004 |
| README setup, run, API examples, design decisions | [README](../../README.md); 004 |
| Unit tests and coverage report | 001–004; 004 verification procedure |
| Dockerfile for frontend + backend together | 004; [Architecture § Build and deployment](../ARCHITECTURE.md#build-and-deployment) |
| Backend microservice architecture | [Architecture](../ARCHITECTURE.md) |
| Out-of-scope discipline | [README § Out of scope](../../README.md#out-of-scope) |

## Order of work

`001 → 002 → 003 → 004`. The contracts are fixed by 001/002 before the UI work in 003; 004 is verified only after the app and tests exist. Implementation may overlap once contracts are agreed, but each change remains independently reviewable (see [`AGENTS.md`](../../AGENTS.md)).

## Spec format

Each spec contains: **What is it?**, **Scope**, semantics/edge cases, **Acceptance criteria**, **Tests**, **Decision**, **Reasons**, **Rejected alternatives**, and **Open questions**. Wire-format details live in [`openapi.yaml`](../../elzzes-calculator-api/docs/openapi.yaml) with narrative in [Architecture § API contract](../ARCHITECTURE.md#api-contract); UI behavior lives in [Design](../DESIGN.md); implementation detail lives in the component conventions. Specs reference them rather than duplicating them.

## Status vocabulary

- **Proposed** — idea, not agreed.
- **Scoped** — requirements and decisions documented; ready to implement.
- **In progress** — implementation under way.
- **Implemented** — code, tests, and docs complete.
- **Superseded** — replaced by a newer spec.

Update the table in the same change that completes a spec.
