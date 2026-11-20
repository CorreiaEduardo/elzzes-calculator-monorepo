# 003 — Calculator UI

**Status:** scoped (documentation baseline, 2026-11-20). Depends on the contracts in specs [001](../001-basic-arithmetic-operations/SPEC.md) and [002](../002-advanced-arithmetic-operations/SPEC.md). UI behavior is canonical in [Design](../../DESIGN.md); this spec defines scope and acceptance. Implementation conventions: [frontend](../../../elzzes-calculator-web/docs/CONVENTIONS.md), [UI library](../../../elzzes-ui-lib/docs/CONVENTIONS.md). Wire contract: [`openapi.yaml`](../../../elzzes-calculator-api/docs/openapi.yaml).

## What is it?

The React + TypeScript front end: keypad, display, entry rules, request orchestration, formatting, error handling, keyboard support, responsive layout, and accessibility — consuming the API for every calculation.

## Scope

**In**

- Calculator shell and keypad for all seven operations.
- Immediate-execution state machine (`useReducer`), one request in flight, abort on AC.
- Client-side entry rules and formatting (no arithmetic).
- Typed API client with TanStack Query, Zod response validation, the timeout from Design, and error mapping.
- Loading, success, domain-error, and transient-error states.
- Keyboard shortcuts, responsive layout, accessibility requirements.
- Presentation components in `@elzzes/ui`; logic in `elzzes-calculator-web`.
- Unit/component tests with Vitest + RTL + MSW, Storybook stories for the UI library, and a coverage report.
- Optional Playwright E2E as a follow-up (not acceptance).

**Out**

- Calculation history, persistence, accounts, `localStorage`.
- Expression parsing/precedence, repeated-equals, memory keys.
- Internationalization, user-selectable themes (OS light/dark only).
- Playwright E2E as an acceptance criterion.

## Requirements

Behavior is canonical in [Design](../../DESIGN.md). The implementation must satisfy it:

| Area | Canonical rule |
| --- | --- |
| Layout | [Design § Layout](../../DESIGN.md#layout) |
| Entry model | [Design § Interaction model](../../DESIGN.md#interaction-model) |
| Operators, equals, chaining | [Design § Interaction model](../../DESIGN.md#interaction-model) |
| Unary operations | [Design § Interaction model](../../DESIGN.md#interaction-model) |
| Request phase and in-flight lock | [Design § Request phase](../../DESIGN.md#request-phase) |
| Error copy and recovery | [Design § Errors and recovery](../../DESIGN.md#errors-and-recovery) |
| Formatting | [Design § Display and formatting](../../DESIGN.md#display-and-formatting) |
| Keyboard | [Design § Keyboard](../../DESIGN.md#keyboard) |
| Accessibility | [Design § Accessibility](../../DESIGN.md#accessibility) |
| Responsive behavior | [Design § Responsive behavior](../../DESIGN.md#responsive-behavior) |
| Visual tokens | [Design § Visual system](../../DESIGN.md#visual-system) |

Implementation constraints:

1. **API integration:** relative `/api/v1/calculate`; JSON; no CORS or absolute host in code. Use a TanStack Query mutation with automatic retry disabled; pass an `AbortController` signal. Validate responses with Zod; map timeout to `TIMEOUT`, transport failure to `NETWORK_ERROR`, malformed responses to `INVALID_RESPONSE`, and 4xx/422/5xx per [Design § Errors and recovery](../../DESIGN.md#errors-and-recovery).
2. **One request in flight:** the reducer locks math keys while calculating; `AC` aborts and resets; late responses cannot mutate state.
3. **No client arithmetic:** the web app may parse typed operands and format results; it may not compute, negate computed results, or use `eval`.
4. **Package boundaries:** `@elzzes/ui` is controlled presentation only; no fetch, arithmetic, operation vocabulary, or global state. Web composes it.
5. **Storybook:** one story file per `@elzzes/ui` component, covering variants and states; `build-storybook` is a CI gate, not a replacement for tests.

## Acceptance criteria

- All seven operations can be performed with mouse/touch and keyboard, with the result always produced by an API request (tests assert the client call rather than trusting output).
- Every entry rule and formatting rule in Design has a unit test; every state transition has a reducer test.
- Domain, validation, network, timeout, malformed-response, and server-error paths show the specified copy and allow recovery/retry without a reload.
- In-flight requests cannot be double-fired; AC aborts and late responses cannot mutate state.
- Accessibility assertions pass: roles/names for all keys, live-region updates on completion, alert on error, `aria-busy` during requests, focus visible and unmoved.
- Responsive checks at the widths and target sizes in [Design § Responsive behavior](../../DESIGN.md#responsive-behavior); no horizontal scroll.
- Web and UI library lint, typecheck, tests, and coverage pass; Storybook builds; CI green.
- No arithmetic implementation, API call, or fetch exists in `@elzzes/ui`.

## Tests

- **Pure units:** entry reducer/rules, formatting rules from [Design § Display and formatting](../../DESIGN.md#display-and-formatting), error mapping, WYSIWYG chaining.
- **Reducer:** all transitions in Design, operator replacement, chaining, unary application/no-op, in-flight lock, AC abort, domain vs transient error recovery.
- **Components:** render, names, keypad interactions, keyboard shortcuts, busy state, alert, focus.
- **API client:** MSW success/400/422/5xx/network/timeout/malformed response; one request per action; abort.
- **UI library:** Storybook stories for variants/states; RTL tests for accessible names, roles, keyboard activation, focus visibility.
- **Responsive/a11y assertions:** DOM order, roles, `aria-keyshortcuts`, no reliance on color; visual checks manual at the listed widths.
- **Optional Playwright:** end-to-end happy path and one error path against the Compose stack, if adopted.

## Decision

A single `useReducer` state machine with immediate execution, TanStack Query for the request lifecycle, Zod at the API boundary, and presentation split into `@elzzes/ui` with Storybook.

## Reasons

- The calculator has one small state shape and one consumer; a reducer is testable without rendering and needs no state library.
- TanStack Query handles the request lifecycle without a bespoke cache; disabling retries preserves the explicit `=` retry contract.
- Zod validates the wire boundary at runtime, catching contract drift that TypeScript types cannot.
- Immediate execution matches the brief and Design; no expression parser means no precedence ambiguity, no `eval`, and less code than a formatter.
- Splitting controlled presentation lets accessibility and visual rules be tested independently of API behavior; Storybook documents the variants.
- One in-flight request with abort gives deterministic tests and prevents out-of-order responses without queues.

## Rejected alternatives

- **Expression parser/evaluator:** out of scope; it needs a grammar, precedence, error recovery, and an `eval`-free evaluator — larger than the whole app.
- **Global state library (Redux/Zustand/Context):** no cross-tree state, no persistence, one consumer. TanStack Query is request lifecycle, not client state.
- **Local arithmetic for responsiveness:** violates the single-source invariant and makes server/client results diverge.
- **Debounced typing requests / live previews:** the brief asks for calculations via the API on action; previews would add requests and state without being requested.
- **Storybook as acceptance:** useful for the UI package; RTL covers behavior and accessibility, Storybook documents components.
- **Playwright E2E as acceptance:** optional follow-up; unit/component/API-client tests carry the contract without browser infrastructure.

## Open questions

None for this scope. History, memory keys, repeated-equals, and Playwright acceptance are explicitly deferred.
