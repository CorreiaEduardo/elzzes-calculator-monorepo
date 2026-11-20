# Decisions

**Status:** backend implementation baseline, 2026-11-20. This is the high-level decision index, not a second specification. Each definition has one canonical home; update that document, not this page.

## Product shape

- A stateless calculator: a React SPA for input and display, a Go REST API for every calculation.
- One versioned endpoint, `POST /api/v1/calculate`, discriminated by `operation`; the wire contract is OpenAPI 3.1.
- Immediate-execution, two-operand calculator interaction; unary keys transform the displayed value. The client parses input and formats output; it never computes.
- No persistence, accounts, history, or authentication. Single-user, local/demo scope.
- In production the SPA and API are separate services joined by Nginx on one origin.

## Architecture style

- Lightweight Ports & Adapters: the calculator core is pure; `transport/http` defines the `Logger` and `Calculator` capabilities it consumes, `logger` provides the slog adapter, and `config` owns environment loading. Ports are consumer-owned; there is no central ports package. The transport may depend on the calculator's public API; the calculator never depends on transport or HTTP.
- One-word package and file names where a plain word works (`config`, `logger`, `http`, `env`), not compound placeholders.
- Rationale and layer rules: [PRINCIPLES § Invariants](PRINCIPLES.md#invariants), [Architecture § Backend](ARCHITECTURE.md#backend), [backend conventions](../elzzes-calculator-api/docs/CONVENTIONS.md).

## Stack

| Area | Choice | Why (short) |
| --- | --- | --- |
| Repository | pnpm workspaces + independent Go module | One lockfile and fast installs for the JS packages; Go keeps its own module. |
| Frontend | React 19 + TypeScript + Vite | Required stack; Vite for dev speed and production bundling. |
| Styling | Tailwind CSS v4 | Utility styling with tokenized theme; no bespoke CSS system. |
| Request lifecycle | TanStack Query | The calculate call without a bespoke fetch layer; no cached history. |
| Validation | Zod | Runtime validation of API responses and external input at the client boundary. |
| UI library | `@elzzes/ui` + Storybook | Controlled presentational components, documented in isolation. |
| Backend | Go stdlib: `net/http`, `encoding/json`, `log/slog` | Small HTTP surface; no framework needed. |
| Contract | OpenAPI 3.1 | Machine-readable contract shared by both layers. |
| Tests | Go `testing`/`httptest`; Vitest + RTL + MSW; Playwright optional | Fast unit/component coverage; E2E optional. |
| Delivery | Docker + Compose + Nginx | Reproducible local run; static SPA and API as separate containers on one origin. |
| CI | GitHub Actions | Runs the same lint/test/build checks on push and PR. |

## Boundaries

| Module | Owns | Must not |
| --- | --- | --- |
| `elzzes-calculator-api` | arithmetic, HTTP contract, OpenAPI | UI, persistence, static SPA |
| `elzzes-calculator-web` | interaction, API client, formatting | compute results, define visual primitives |
| `elzzes-ui-lib` (`@elzzes/ui`) | controlled presentational components, tokens, stories | fetch, arithmetic, application state |

`web` depends on `@elzzes/ui` and the HTTP API; `@elzzes/ui` depends on neither; the API depends on neither frontend package. There is no shared runtime state.

## Canonical details

| Topic | Source of truth |
| --- | --- |
| Wire contract, schemas, error codes, statuses | [`openapi.yaml`](../elzzes-calculator-api/docs/openapi.yaml); narrative in [Architecture § API contract](ARCHITECTURE.md#api-contract) |
| Operation behavior and edge cases | [spec 001](specs/001-basic-arithmetic-operations/SPEC.md), [spec 002](specs/002-advanced-arithmetic-operations/SPEC.md) |
| System structure, request flow, numeric model, deployment | [ARCHITECTURE.md](ARCHITECTURE.md) |
| UI behavior, formatting, error copy, accessibility, responsive design | [DESIGN.md](DESIGN.md) |
| Backend implementation conventions | [API conventions](../elzzes-calculator-api/docs/CONVENTIONS.md) |
| Frontend implementation conventions | [Web conventions](../elzzes-calculator-web/docs/CONVENTIONS.md) |
| UI library implementation conventions | [UI library conventions](../elzzes-ui-lib/docs/CONVENTIONS.md) |
| Invariants and decision rules | [PRINCIPLES.md](PRINCIPLES.md) |
| Delivery, tests, coverage, CI verification | [spec 004](specs/004-delivery-and-verification/SPEC.md) |

## Rejected alternatives

- **One endpoint per operation** — more routes and duplicated decode/validation; one discriminated endpoint keeps one validation path. See [spec 001](specs/001-basic-arithmetic-operations/SPEC.md#rejected-alternatives).
- **Client-side arithmetic** — violates the single-source invariant and makes results unverifiable. See [PRINCIPLES.md](PRINCIPLES.md#invariants).
- **Expression parser / `eval`** — grammar, precedence, and security surface beyond the brief. See [spec 003](specs/003-calculator-ui/SPEC.md#rejected-alternatives).
- **Exact decimal arithmetic** — square root and fractional powers are irrational; binary64 plus UI formatting is the agreed model. See [Architecture § Numeric model](ARCHITECTURE.md#numeric-model).
- **Global state library (Redux/Zustand/Context)** — one reducer and one consumer; TanStack Query covers request lifecycle only. See [spec 003](specs/003-calculator-ui/SPEC.md#rejected-alternatives).
- **Axios / React Hook Form / large UI kit** — native `fetch`, a keypad, and the internal UI library cover the needs. See [web conventions](../elzzes-calculator-web/docs/CONVENTIONS.md).
- **Database, auth, history, cloud deployment** — no requirement; local-first scope. See [README § Out of scope](../README.md#out-of-scope).
- **Single Go image embedding the SPA** — fewer units, but couples frontend builds to the API image and moves static serving into Go; the agreed topology uses Nginx and API containers. See [Architecture § Build and deployment](ARCHITECTURE.md#build-and-deployment).
- **Kubernetes, gRPC, message brokers, event sourcing, ORM** — no requirement at this scale. See [PRINCIPLES.md § Simplicity budget](PRINCIPLES.md#invariants).
