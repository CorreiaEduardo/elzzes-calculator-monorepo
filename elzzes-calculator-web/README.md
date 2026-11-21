# elzzes-calculator-web

The React + TypeScript SPA for Elzzes Calculator. It collects input, orchestrates API requests, formats results, and presents errors. It does **not** perform arithmetic.

**Status:** implemented — calculator UI, immediate-execution reducer, TanStack Query + Zod API client, keyboard/accessibility behavior, and tests. UI behavior is specified in [`docs/DESIGN.md`](../docs/DESIGN.md) and [spec 003](../docs/specs/003-calculator-ui/SPEC.md). Stack: React 19, TypeScript, Vite, Tailwind CSS v4, TanStack Query, Zod.

## Responsibilities

- Keypad and display for all seven operations, keyboard support, responsive layout.
- Entry rules, loading state, error copy, and recovery/retry.
- Typed API client (`POST /api/v1/calculate`) with TanStack Query, Zod response validation, the [timeout from Design](../docs/DESIGN.md#errors-and-recovery), and error mapping.
- Result formatting per [Design § Display and formatting](../docs/DESIGN.md#display-and-formatting).
- Composition of presentational components from `@elzzes/ui`.

## Run

From the repo root:

```bash
pnpm install
pnpm --filter elzzes-calculator-web dev   # http://localhost:5173
```

The dev server proxies `/api` to `http://localhost:8080`; start the API first. In production Nginx serves the SPA and proxies `/api/` and `/healthz` to the Go service on one origin.

## Docker

The root `compose.yaml` builds this SPA and serves it through Nginx alongside the API:

```bash
docker compose up --build   # http://localhost:8080
```

The image is multi-stage: pnpm workspace install and Vite build, then unprivileged Nginx serving `dist/` and proxying `/api` and `/healthz` to `api:8080`.

## Scripts

| Command | Purpose |
| --- | --- |
| `pnpm --filter elzzes-calculator-web dev` | Vite dev server with API proxy |
| `pnpm --filter elzzes-calculator-web build` | production build to `dist/` |
| `pnpm --filter elzzes-calculator-web preview` | preview the built SPA |
| `pnpm --filter elzzes-calculator-web lint` | ESLint |
| `pnpm --filter elzzes-calculator-web typecheck` | TypeScript, no emit |
| `pnpm --filter elzzes-calculator-web test` | Vitest |
| `pnpm --filter elzzes-calculator-web test:coverage` | tests + coverage report |

## Testing

Vitest + React Testing Library + user-event, with MSW for API responses. Covered: reducer transitions, entry rules, formatting, keyboard mapping, loading/error/retry states, accessibility roles/names, API-client success and failure mapping, and StrictMode single-request behavior. Run `pnpm --filter elzzes-calculator-web test:coverage` for the coverage report. Playwright E2E is optional. The coverage target and CI are canonical in [spec 004](../docs/specs/004-delivery-and-verification/SPEC.md).

## Conventions and layout

The layout, state flow, API-client rules, styling, and testing conventions are canonical in [`docs/CONVENTIONS.md`](docs/CONVENTIONS.md). High-level structure: `src/app/`, `src/features/calculator/`, `src/styles/`, `src/test/`.

## Boundaries

- Never evaluates arithmetic or negates computed results; every result comes from the API.
- No absolute API host, CORS, or environment-specific base URL: always relative `/api/v1`.
- No global state library; the calculator has one reducer and one consumer, and TanStack Query owns only the request lifecycle.
- No visual primitives defined here — they belong to `@elzzes/ui`.

See the [root README](../README.md) for setup, API examples, and design decisions.
