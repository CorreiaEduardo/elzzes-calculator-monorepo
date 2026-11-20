# Frontend conventions

**Status:** planned baseline, 2026-11-20. Canonical for `elzzes-calculator-web` implementation conventions. UI behavior is canonical in [Design](../../docs/DESIGN.md); acceptance is in [spec 003](../../docs/specs/003-calculator-ui/SPEC.md); the wire contract is in [openapi.yaml](../../elzzes-calculator-api/docs/openapi.yaml); system structure is in [Architecture](../../docs/ARCHITECTURE.md).

## Stack

React 19, TypeScript strict, Vite, Tailwind CSS v4, TanStack Query, Zod, Vitest + React Testing Library + user-event + MSW, ESLint + Prettier. Package manager: pnpm workspaces.

## Layout

```text
elzzes-calculator-web/
├── src/
│   ├── app/                    # App shell, providers (QueryClientProvider), global styles
│   ├── features/calculator/
│   │   ├── components/         # Calculator composition, display, keypad wiring
│   │   ├── hooks/              # useCalculator reducer and request orchestration
│   │   ├── api/                # fetch client and error mapping
│   │   ├── schemas/            # Zod request/response schemas
│   │   └── types/
│   ├── lib/                    # shared client utilities
│   ├── styles/                 # Tailwind entry
│   └── test/                   # MSW handlers and fixtures
├── docs/CONVENTIONS.md
├── Dockerfile
├── package.json
├── tsconfig.json
└── vite.config.ts
```

- Feature code stays inside `features/calculator/`; `app/` composes and provides.
- No cross-feature imports. Shared utilities go in `lib/` only when used by more than one feature.

## State and data flow

- One `useReducer` owns transient calculator state: entry string, stored first operand/result, pending operation, phase (`idle` / `calculating` / `result` / `error`).
- TanStack Query owns the request lifecycle for the calculate call. Use a mutation, not a cached query: calculations are actions, not server state. Disable automatic retries (`retry: false`); the user retries with `=` per [Design](../../docs/DESIGN.md#errors-and-recovery).
- Exactly one request in flight. The reducer locks math keys while `calculating`; `AC` aborts via `AbortController` and resets. Ignore aborted or superseded responses.
- No Redux, Zustand, or Context for calculator state. TanStack Query is not a global client-state store.
- Nothing is persisted: no `localStorage`, session, URL state, or history.

## API client and validation

- Always call the relative path `/api/v1/calculate`; never an absolute host, CORS URL, or environment-specific base URL. Vite proxies in development; Nginx proxies in production.
- A single typed `calculate(operation, operands, signal)` function returns a discriminated success/error result. Components never call `fetch` directly.
- Validate the response with Zod before it reaches the reducer; a malformed response becomes a client error, not a crash.
- Apply the timeout defined in [Design § Errors and recovery](../../docs/DESIGN.md#errors-and-recovery). Map transport failures to `NETWORK_ERROR`, timeouts to `TIMEOUT`, 5xx to `INTERNAL_ERROR`, malformed responses to `INVALID_RESPONSE`, and documented 4xx/422 codes per that table.
- Keep server codes and statuses sourced from [openapi.yaml](../../elzzes-calculator-api/docs/openapi.yaml); do not duplicate the contract in TypeScript comments.
- The client never computes, negates computed results, or uses `eval`.

## Formatting

Formatting is a pure presentation module. The rules are canonical in [Design § Display and formatting](../../docs/DESIGN.md#display-and-formatting); implement them there, test them independently, and do not round in the API client.

## Rejected libraries

- **Axios** — native `fetch` is sufficient and keeps the dependency surface small.
- **React Hook Form** — the calculator is a keypad, not a form; the reducer owns entry state.
- **Large UI kit** — Tailwind plus `@elzzes/ui` cover the visual needs.
- **Redux/Zustand/Context for calculator state** — one reducer and one consumer; TanStack Query owns request lifecycle only.

## Styling

- Tailwind CSS v4 utility classes. Theme tokens come from `@elzzes/ui` (CSS custom properties/Tailwind theme); do not hard-code colors.
- No CSS Modules or ad-hoc CSS files for component styling; the Tailwind entry lives in `src/styles/`.
- Do not define visual primitives here; they belong to `@elzzes/ui`.

## Components and behavior

- `Calculator` composes `CalculatorFrame`, `Display`, and `Keypad` from `@elzzes/ui`; it wires events to `useCalculator`.
- Keep display strings formatted in the web layer; `@elzzes/ui` receives formatted strings and callbacks only.
- Keyboard handling and accessibility behavior follow [Design § Keyboard](../../docs/DESIGN.md#keyboard) and [Design § Accessibility](../../docs/DESIGN.md#accessibility); do not fork the rules.

## Testing

- Vitest + React Testing Library + user-event for components; MSW for API responses. No test hits a real network.
- Reducer tests cover every transition in Design; pure tests cover entry rules, formatting, and error mapping; API-client tests cover success, 4xx/422, 5xx, network failure, timeout, malformed response, and abort.
- Assert the client call happened; do not assert a result the client computed.
- Coverage target and commands: [spec 004](../../docs/specs/004-delivery-and-verification/SPEC.md). Playwright E2E is optional and not an acceptance criterion.

## Scripts

| Command | Purpose |
| --- | --- |
| `pnpm --filter elzzes-calculator-web dev` | Vite dev server with API proxy |
| `pnpm --filter elzzes-calculator-web build` | production build to `dist/` |
| `pnpm --filter elzzes-calculator-web preview` | preview the built SPA |
| `pnpm --filter elzzes-calculator-web typecheck` | TypeScript, no emit |
| `pnpm --filter elzzes-calculator-web lint` | ESLint |
| `pnpm --filter elzzes-calculator-web test` | Vitest |
| `pnpm --filter elzzes-calculator-web test:coverage` | tests + coverage report |

## Boundaries

- No arithmetic, API contract definitions, or visual primitives.
- No absolute API host or CORS configuration.
- No global state library.
- No direct `fetch` in components or hooks outside the API client.
