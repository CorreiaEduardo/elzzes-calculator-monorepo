# Architecture

**Status:** implementation baseline, stack agreed 2026-11-20. The API implements all seven operations and `GET /healthz`; the calculator UI, API client, and shared UI primitives are implemented. Spec 004 delivery is implemented: Docker Compose/Nginx, the root Makefile, GitHub Actions CI, and enforced coverage floors. This document is canonical for system structure, request flow, state ownership, the numeric model, and deployment. The wire contract is canonical in [`openapi.yaml`](../elzzes-calculator-api/docs/openapi.yaml); operation behavior and edge cases in [spec 001](specs/001-basic-arithmetic-operations/SPEC.md) and [spec 002](specs/002-advanced-arithmetic-operations/SPEC.md); UI behavior in [Design](DESIGN.md); implementation details in the component conventions.

## System context

```text
┌──────────────────────────── Browser ────────────────────────────┐
│  React SPA  (elzzes-calculator-web + @elzzes/ui)                │
└───────────────────────────────┬─────────────────────────────────┘
                                │ relative /api/v1/calculate
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│  web edge: Nginx (production) / Vite dev server                 │
│    /        → static SPA                                        │
│    /api/*   → Go API                                            │
│    /healthz → Go API                                            │
└───────────────────────────────┬─────────────────────────────────┘
                                │ JSON over HTTP
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│  api service: Go (elzzes-calculator-api)                        │
│    internal/calculator                pure arithmetic core     │
│    internal/transport/http            net/http adapter         │
│    internal/logger                    slog adapter             │
│    internal/config                    environment config      │
└─────────────────────────────────────────────────────────────────┘
```

Two topologies, one contract:

- **Development:** Vite serves the SPA on `:5173` and proxies `/api` to the Go service on `:8080`. The Go service runs API-only.
- **Production:** `web` (Nginx) and `api` (Go) run as separate Compose services. Nginx serves the built SPA and proxies `/api/` and `/healthz` to the API. One origin, no CORS.

Both topologies use relative `/api/...` paths, so the frontend has no environment-specific API URL. The Go service never serves the SPA.

## Packages

| Path | Package | Stack | Owns | Must not |
| --- | --- | --- | --- | --- |
| `elzzes-calculator-api/` | Go module `elzzes-calculator-api` | Go stdlib (`net/http`, `encoding/json`, `log/slog`) | arithmetic, HTTP API, OpenAPI contract | contain UI logic, persistence, or static SPA |
| `elzzes-calculator-web/` | `elzzes-calculator-web` | React 19, TypeScript strict, Vite, Tailwind CSS v4, TanStack Query, Zod | calculator interaction, API client, formatting | compute results, define visual primitives |
| `elzzes-ui-lib/` | `@elzzes/ui` | React, TypeScript, Tailwind CSS v4, CVA/clsx, Storybook | controlled presentational components and tokens | fetch, arithmetic, application state |

The JavaScript packages are pnpm workspaces; the Go module is independent. Compose runs the web and API services together. Packages communicate only through the HTTP contract and the `@elzzes/ui` component API; there is no shared runtime state.

Implementation details live in the component conventions: [backend](../elzzes-calculator-api/docs/CONVENTIONS.md), [frontend](../elzzes-calculator-web/docs/CONVENTIONS.md), [UI library](../elzzes-ui-lib/docs/CONVENTIONS.md).

## Request flow

1. The user presses a key. The web reducer decides whether this is input editing or a calculation.
2. Calculations go through a TanStack Query mutation and the typed API client as `POST /api/v1/calculate` with `operation` and `operands`.
3. The HTTP adapter validates content type, the 4 KiB body cap, strict JSON fields, operation/arity, and operands, then calls the calculator core through its port.
4. Core success becomes `200 {"operation": ..., "operands": ..., "result": ...}`; a domain failure becomes a documented `4xx` code; an unexpected failure is logged and becomes `500 INTERNAL_ERROR`.
5. The client validates the response with Zod, maps it to a display value or user-facing error, formats it, and announces the outcome.
6. The service stores nothing. Every calculation is independent.

The client is deliberately never the authority. It may reject obviously invalid input for UX, but the service validates everything again.

## Backend

The service is a lightweight Ports & Adapters design: `internal/calculator` is the pure core; `internal/transport/http` is the net/http adapter (routing, middleware, strict decoding, validation, status/code mapping, JSON error envelope, health) and declares the `Logger` capability it needs in its own package; `internal/logger` implements that port with `log/slog`; `internal/config` loads environment configuration. Ports are consumer-owned — there is no central ports package. The transport may depend on the calculator's public API (`Operation`, `ParseOperation`, `Arity`, and domain errors), but the dependency never runs the other way: the core never imports `net/http`, `log/slog`, `os`, or transport concerns. Replacing the transport, logger, or configuration source therefore touches only that package and the composition root. `cmd/api` is the only place adapters meet. The exact layout, naming, error handling, testing, and logging rules are canonical in [backend conventions](../elzzes-calculator-api/docs/CONVENTIONS.md).

### HTTP stack

Go standard library only: `net/http` with method/pattern routing, `encoding/json`, `log/slog`. The service enforces the 4 KiB body cap and JSON content type, rejects unknown fields and out-of-range operands, validates operation arity, maps domain errors to `422`, and returns JSON `404`/`405` responses. It also provides panic recovery, request logging, explicit timeouts, and graceful shutdown. All of this lives in the `internal/transport/http` adapter; `log/slog` is confined to `internal/logger`. It serves no static files and sets no CORS headers: Nginx serves the SPA and proxies the API on one origin; the Vite dev proxy preserves that model. Exact middleware, validation, and error-mapping rules are canonical in [backend conventions](../elzzes-calculator-api/docs/CONVENTIONS.md).

## API contract

The machine-readable contract is [`openapi.yaml`](../elzzes-calculator-api/docs/openapi.yaml). It is canonical for paths, request and response schemas, operation vocabulary, arity, status codes, and error codes. This section only explains the shape; it does not restate the schema.

- `POST /api/v1/calculate` — request `{operation, operands}`; success `{operation, operands, result}` with unrounded binary64; errors `{"error":{"code":"...","message":"..."}}`.
- `GET /healthz` — `{"status":"ok"}`.
- Status classes: `400` malformed request; `413`/`415` transport limits; `422` well-formed request with an undefined result; `404`/`405` route/method; `500` unexpected failure.
- `message` is developer-facing and not a UI contract; the client owns user-facing copy keyed by `code` ([Design § Errors and recovery](DESIGN.md#errors-and-recovery)).
- No operation-discovery endpoint, per-operation endpoints, authentication, rate limiting, or persistence at this layer. Rejected alternatives are recorded in [spec 001](specs/001-basic-arithmetic-operations/SPEC.md#rejected-alternatives).

## Numeric model

**Decision:** numbers are IEEE-754 binary64 end to end. JSON numbers decode to `float64`; results are `float64`; the API does not round.

Rules enforced by the service:

- Inputs are finite by construction (JSON cannot express `NaN`/`Infinity`).
- Division checks the divisor before dividing, including `-0`.
- Any result that is `NaN` or ±`Infinity` is rejected as `NON_FINITE_RESULT`.
- Negative-zero results are normalized to `0` before encoding.
- `0^0 = 1`; `0` raised to a negative exponent is non-finite and therefore an error; a negative base with a non-integer exponent is undefined and an error.
- `percent` is unary: `x → x / 100`.
- Underflow to `0` is accepted (not an error).

**Why not exact arithmetic?** `math/big.Rat` would make the four basic operations exact, but square root and fractional powers are irrational; a rational model would still need a float path, plus string/decimal conversion at the API boundary. A decimal library adds a third-party dependency and still approximates sqrt/pow. Binary64 matches JSON, the browser, and Go's math library with one model. The cost is visible artifacts (`0.1 + 0.2 = 0.30000000000000004`), which the UI hides by formatting results for display ([Design § Display and formatting](DESIGN.md#display-and-formatting)). Exact decimal arithmetic is a documented future option if the product needs ledger-grade results; it is not this product.

The UI formats what it displays; the service returns the raw value. The displayed (formatted) value is what the user carries into the next operation — WYSIWYG — so hidden precision never surprises them ([Design § Display and formatting](DESIGN.md#display-and-formatting)).

## Frontend

The app is organized by feature: `app/` composes the shell and providers; `features/calculator/` owns the reducer, API client, schemas, and formatting; `@elzzes/ui` supplies controlled presentation. The exact layout and coding rules are canonical in [frontend conventions](../elzzes-calculator-web/docs/CONVENTIONS.md).

- **State machine.** One `useReducer` owns transient calculator state: current entry, stored result/first operand, pending operation, phase (`idle` / `calculating` / `result` / `error`). No Redux, Zustand, or Context — there is one consumer and one small state shape.
- **Request lifecycle.** TanStack Query runs the calculate mutation; automatic retries are disabled so `=` remains the explicit retry. One request is in flight at a time; `AC` aborts via `AbortController` and resets. Responses are ignored if aborted or superseded.
- **Validation.** Zod validates API responses (and relevant external input) at the boundary. The service remains authoritative for requests.
- **API client.** A single typed `calculate(operation, operands, signal)` function returns a discriminated success/error result. It applies the timeout and error mapping defined in [Design § Errors and recovery](DESIGN.md#errors-and-recovery), and never throws raw fetch errors into components.
- **Formatting.** Canonical rules live in [Design § Display and formatting](DESIGN.md#display-and-formatting); the web implements them as a pure module and does not round in the API client.
- **No arithmetic in the client.** Parsing typed input and formatting output only.

### State ownership

| State | Owner | Lifetime |
| --- | --- | --- |
| Entry string, pending operation, phase | web reducer | until AC/reload |
| Request lifecycle (in flight, success, failure) | TanStack Query mutation | per request |
| Computed results, validity | Go service | per request |
| Formatted display, error copy | web (presentation) | derived |
| Component focus/hover | `@elzzes/ui` / browser | ephemeral |

Nothing is persisted: no `localStorage`, session, URL state, or backend history. TanStack Query is not used as a history cache.

## Shared UI library

`@elzzes/ui` (folder `elzzes-ui-lib`) is an internal, unpublished package containing generic, controlled presentational components and design tokens:

- The library exposes `Button`, `Card`, `Display`, and `Keypad`; the web app maps calculator key roles to generic variants. Names, variants, and slots are defined in [Design § Components](DESIGN.md#components); theme tokens are Tailwind CSS custom properties.
- Components receive formatted strings and callbacks; they hold no calculator state and know nothing about the API or operations.
- Tailwind CSS v4 with tokenized colors, spacing, radii, and type. Light/dark via `prefers-color-scheme`.
- Accessibility contract: real buttons, accessible names, visible focus, minimum target size, `aria-keyshortcuts` hints. Documented in [Design § Accessibility](DESIGN.md#accessibility).
- Storybook documents components in isolation and is part of the package deliverable.

Component conventions: [UI library](../elzzes-ui-lib/docs/CONVENTIONS.md).

## Testing strategy

| Layer | Tools | Covers |
| --- | --- | --- |
| Go core (`internal/calculator`) | `go test`, table-driven | every operation, documented edge cases, float comparisons with tolerance |
| Go HTTP (`internal/transport/http`) | `httptest` | status/code mapping, strict decoding, limits, health, method/route handling |
| Web logic | Vitest | reducer transitions, entry rules, formatting, error mapping |
| Web UI | Vitest + React Testing Library + user-event | keypad, keyboard parity, loading/error states, accessible roles/names |
| Web API client | Vitest + MSW | success, 4xx/422, 5xx, timeout/network failure, malformed response, abort; fixtures mirror the OpenAPI contract |
| UI library | Vitest + RTL + Storybook | component behavior, accessible names, variants, visual states |

No test hits a real network. The operation edge-case matrix in specs 001–002 is the authoritative test list; coverage is a guardrail, not the goal.

Commands, the coverage target, and CI checks are canonical in [spec 004](specs/004-delivery-and-verification/SPEC.md). Test-level conventions are in the component conventions.

## Build and deployment

```text
docker compose up --build
├── web service  (elzzes-calculator-web/Dockerfile)
│   └── multi-stage: pnpm install → build SPA → Nginx unprivileged, listens :8080
└── api service  (elzzes-calculator-api/Dockerfile)
    └── multi-stage: Go build ./cmd/api → minimal non-root runtime, listens :8080
```

- Nginx serves the static SPA and proxies `/api/` and `/healthz` to `api:8080`. One origin, no CORS, no Node runtime in the final web image.
- The two service Dockerfiles plus `compose.yaml` are the agreed reading of the brief's "Dockerfile to run frontend + backend together": `docker compose up` runs both. This assumption is recorded in [README § How to run](../README.md#how-to-run).
- `PORT` defaults to `8080`. The stack is stateless and needs no volumes or secrets; `/healthz` is the probe.
- Root and package `.dockerignore` files exclude `.git`, `node_modules`, build output, coverage, and local env files.
- Images run as non-root where the base image permits; the final API image contains no Go toolchain or source.
- Local development does not use the image: `make dev` (or pnpm + Vite + `go run`) as described in [README § Local development](../README.md#local-development). The root `Makefile` (`install`, `dev`, `storybook`, `verify`, `compose-verify`) wraps the package scripts and adds no logic or dependency.
- GitHub Actions runs the same `make verify` checks plus a Compose build-and-smoke job on push and pull requests ([spec 004](specs/004-delivery-and-verification/SPEC.md#ci)).

**Rejected alternative:** a single Go image embedding the SPA via `go:embed`. It has fewer units, but couples every frontend build to the API image and moves static serving into the Go binary. The agreed topology keeps static serving in Nginx and the API focused on arithmetic.

## Persistence

None. There is no database, schema, migration, or server-side history. The domain is a pure function; a calculation leaves no record. A client-side history panel was considered and deferred (out of scope), and would not change the API.

## Security posture

- Strict JSON decoding, the documented body cap, explicit content type, no user-controlled routing.
- No CORS by design: the SPA and API share an origin through Nginx in production; the dev proxy preserves that model.
- No secrets, accounts, cookies, or persisted user data. No auth because the app is single-user/local by scope.
- Nginx/API responses set minimal hardening headers (`X-Content-Type-Options: nosniff`, frame denial, referrer policy). Rate limiting and CSP hardening are deployment concerns outside this brief.
- CI uses repository-scoped credentials only; no deployment or registry publishing.

## Trade-offs and known limitations

- **Binary64, not exact decimals** — hidden artifacts are formatted away, but the API exposes raw values (`0.30000000000000004`). Deliberate; alternative weighed above.
- **Immediate-execution model** — no operator precedence or expression grammar. Matches a physical calculator and keeps parsing out of scope; users expecting an expression parser will not find one.
- **One endpoint** — less granular than REST resources per operation, but a fixed vocabulary with one validation path is smaller and easier to evolve with `/api/v1`.
- **Two containers and an Nginx proxy** — more moving parts than one embedded-SPA binary, in exchange for clear static/API separation and a standard deployment shape.
- **Duplicated contract types** — OpenAPI schemas, Go DTOs, and Zod schemas are maintained by hand. Tests and CI pin the wire format; codegen is an option if the contract grows.
- **Playwright E2E is optional** — not part of acceptance; unit/component/API-client tests carry the contract.
- **No database or auth** — deliberate local-first scope; adding either would be a new architecture decision.

## Changing the architecture

Structural decisions are agreed with the user first, recorded in a spec, `ARCHITECTURE.md`, or `openapi.yaml`, then implemented, then reflected back into the canonical document. See [`AGENTS.md`](../AGENTS.md) and [`PRINCIPLES.md`](PRINCIPLES.md).
