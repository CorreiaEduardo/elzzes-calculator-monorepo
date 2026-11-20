# Frontend conventions

**Status:** implementation baseline, 2026-11-20. Canonical for `elzzes-calculator-web` implementation conventions. UI behavior is canonical in [Design](../../docs/DESIGN.md); acceptance is in [spec 003](../../docs/specs/003-calculator-ui/SPEC.md); the wire contract is in [openapi.yaml](../../elzzes-calculator-api/docs/openapi.yaml); system structure is in [Architecture](../../docs/ARCHITECTURE.md).

## Stack

React 19, TypeScript strict, Vite, Tailwind CSS v4, TanStack Query, Zod, Vitest + React Testing Library + user-event + MSW, ESLint + Prettier. Package manager: pnpm workspaces.

## Layout

```text
elzzes-calculator-web/
├── src/
│   ├── app/                    # App shell and QueryClientProvider
│   ├── features/calculator/
│   │   ├── api/                # calculate client, error mapping, Zod schemas
│   │   ├── components/         # Calculator composition and key buttons
│   │   ├── hooks/              # useCalculator, useKeyboard
│   │   ├── copy.ts             # error copy, busy text
│   │   ├── entry.ts            # entry rules
│   │   ├── expression.ts       # display/expression projections and symbols
│   │   ├── format.ts           # result formatting
│   │   ├── keys.ts             # key definitions, names, variants, shortcuts
│   │   ├── keyboard.ts         # keydown mapping
│   │   ├── reducer.ts          # calculator state machine
│   │   └── types.ts
│   ├── styles/                 # Tailwind entry, base page styles
│   └── test/                   # setup and MSW server
├── docs/CONVENTIONS.md
├── Dockerfile
├── package.json
├── tsconfig.json
└── vite.config.ts
```

- Feature code stays inside `features/calculator/`; `app/` composes and provides.
- No cross-feature imports. Shared utilities go in a new `lib/` only when used by more than one feature.

## State and data flow

- One `useReducer` owns transient calculator state: entry string, stored first operand/result, pending operation, phase (`idle` / `calculating` / `result` / `error`).
- TanStack Query owns the request lifecycle for the calculate call. Use a mutation, not a cached query: calculations are actions, not server state. Disable automatic retries (`retry: false`); the user retries with `=` per [Design](../../docs/DESIGN.md#errors-and-recovery).
- Exactly one request in flight. The reducer locks math keys while `calculating`; `AC` aborts via `AbortController` and resets. Ignore aborted or superseded responses.
- Requests carry a monotonic id; the reducer stores a request plan (`operation`, `operands`, continuation, expression, display) and ignores stale ids. Continuations are `equals`, `chain`, `unary-current`, or `unary-second`, so retries reapply the original outcome.
- `useCalculator` starts the mutation in an effect keyed by the request id, guarded by a started-id ref so React StrictMode does not double-fire. `AC` aborts explicitly; effect cleanup does not abort.
- No Redux, Zustand, or Context for calculator state. TanStack Query is not a global client-state store.
- Nothing is persisted: no `localStorage`, session, URL state, or history.

## API client and validation

- Always use the relative path `/api/v1/calculate`, resolved with `window.location.origin`; never an absolute host, CORS URL, or environment-specific base URL. Vite proxies in development; Nginx proxies in production.
- A single typed `calculate(operation, operands, signal)` function returns a discriminated `success | error | aborted` result. Components never call `fetch` directly.
- Validate the response with strict Zod schemas before it reaches the reducer; a malformed response becomes a client error, not a crash.
- Apply the 10-second timeout from [Design § Errors and recovery](../../docs/DESIGN.md#errors-and-recovery). Map transport failures to `NETWORK_ERROR`, timeouts to `TIMEOUT`, 5xx to `INTERNAL_ERROR`, malformed/unknown responses to `INVALID_RESPONSE`, 422 domain codes to their copy, and other valid-envelope 4xx responses to the domain invalid-calculation copy.
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

- `Calculator` composes `Card`, `Display`, and `Keypad` from `@elzzes/ui`; `KeyButton` maps key roles to generic `Button` variants. It wires events to `useCalculator`.
- Keep display strings formatted in the web layer; `@elzzes/ui` receives formatted strings, flags, labels, and callbacks only. Never pass operation vocabulary or user copy into the library.
- Keyboard handling and accessibility behavior follow [Design § Keyboard](../../docs/DESIGN.md#keyboard) and [Design § Accessibility](../../docs/DESIGN.md#accessibility); do not fork the rules.

## Testing

- Vitest + React Testing Library + user-event for components; MSW for API responses. No test hits a real network or sleeps; in-flight tests use deferred gate promises, and only the API-client timeout test uses fake timers.
- Reducer tests cover every transition in Design; pure tests cover entry rules, formatting, and keyboard mapping; API-client tests cover success, 400/422/5xx, network failure, timeout, malformed response, schema drift, and abort.
- Component tests assert the API request body, not a client-computed result; include WYSIWYG and StrictMode single-request checks.
- Vitest uses the threads pool because forks hang in this WSL environment.
- Coverage: `pnpm --filter elzzes-calculator-web test:coverage` (v8 text + HTML under `coverage/`). Playwright E2E is optional and not an acceptance criterion.

## Scripts

| Command | Purpose |
| --- | --- |
| `pnpm --filter elzzes-calculator-web dev` | Vite dev server with API proxy |
| `pnpm --filter elzzes-calculator-web build` | production build to `dist/` |
| `pnpm --filter elzzes-calculator-web preview` | preview the built SPA |
| `pnpm --filter elzzes-calculator-web typecheck` | TypeScript, no emit |
| `pnpm --filter elzzes-calculator-web lint` | ESLint |
| `pnpm --filter elzzes-calculator-web test` | Vitest (run once) |
| `pnpm --filter elzzes-calculator-web test:watch` | Vitest watch mode |
| `pnpm --filter elzzes-calculator-web test:coverage` | tests + coverage report |

## Boundaries

- No arithmetic, API contract definitions, or visual primitives.
- No absolute API host or CORS configuration.
- No global state library.
- No direct `fetch` in components or hooks outside the API client.
