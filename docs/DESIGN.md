# Product and UI design

**Status:** implementation baseline, 2026-11-20. Canonical for UI behavior; acceptance criteria live in [spec 003](specs/003-calculator-ui/SPEC.md). The shared library exposes generic primitives; calculator-specific names below were superseded during implementation. Implementation conventions: [web](../elzzes-calculator-web/docs/CONVENTIONS.md), [UI library](../elzzes-ui-lib/docs/CONVENTIONS.md). Update the canonical section and the affected spec together.

## Design goals

- The calculation and its result are always obvious. One glance answers: what am I entering, what will happen, what happened.
- Errors explain the problem and leave a way forward. No dead ends, no modal interruptions.
- Thumb-friendly on a phone; keyboard-first on a desktop.
- Accessible by default: labels, focus, announcements, contrast, target size.
- The interface looks like a calculator, not a dashboard. One surface, one accent, no decoration that competes with digits.

Out of scope: calculation history, scientific modes, memory keys, unit conversion, localization, user-selectable themes.

## Layout

```text
┌─────────────────────────────────────────┐
│ 12.5 +                                  │  expression / status (live region)
│                                   0     │  display, right-aligned
├─────────────────────────────────────────┤
│   [   √   ]   [  xʸ   ]   [   %   ]     │  advanced strip
├──────┬──────┬──────┬────────────────────┤
│  AC  │  ⌫   │  ±   │   ÷                │
├──────┼──────┼──────┼────────────────────┤
│  7   │  8   │  9   │   ×                │
├──────┼──────┼──────┼────────────────────┤
│  4   │  5   │  6   │   −                │
├──────┼──────┼──────┼────────────────────┤
│  1   │  2   │  3   │   +                │
├──────┴──────┼──────┼────────────────────┤
│      0      │  .   │   =                │
└─────────────┴──────┴────────────────────┘
```

- Familiar 4-column digit grid keeps muscle memory: `7-8-9` and `4-5-6` rows, operators on the right rail, `AC/⌫/±` top-left, `=` bottom-right.
- `0` spans two columns; `=` is the only accent-filled key.
- Advanced operations sit in a separate strip above the grid so they are discoverable without displacing the basic layout. A mode toggle or second keypad page was rejected: hidden modes hurt discoverability and add state.
- Content width is capped at 420px; the card is centered, with the page background distinguishable from the calculator surface.

## Display and formatting

Two text rows:

1. **Expression / status line** — small, muted, left-aligned. Shows the pending expression without the number currently being typed (`12.5 +`), so it does not churn on every keystroke. After a completed calculation it shows the full equation and outcome (`12.5 + 7.5 = 20`). During a request it shows `Calculating…`. It is the polite live region.
2. **Display** — the current entry or formatted result, right-aligned, large, tabular numerals. It visually replaces a leading `0` when typing, shrinks in steps to a 24px floor to fit, and never wraps or scrolls horizontally.

Formatting rules (implemented as one pure module and tested independently):

- Results render with at most **12 significant digits**, trailing zeros trimmed.
- Exponential notation (lowercase `e`, explicit sign) when |x| ≥ 10¹² or 0 < |x| < 10⁻⁹: `1e+15`, `1.2e-10`. Zero is `0`.
- No thousands separators: values stay easy to read back and compare.
- Negative zero displays as `0`.
- Entries show exactly what was typed, up to **15 digits** (sign and decimal point excluded); further digits are ignored.
- **WYSIWYG chaining:** after a result, the formatted value shown is the value carried into the next operation. The user never operates on hidden precision they cannot see. Example: `0.1 + 0.2` displays `0.3`; pressing `× 3 =` computes `0.3 × 3`, not `0.30000000000000004 × 3`.

The raw API result remains unrounded; formatting is presentation ([Architecture § Numeric model](ARCHITECTURE.md#numeric-model)).

## Interaction model

**Immediate-execution, two-operand calculator.** There is no expression grammar, operator precedence, or expression parser. Operations execute when `=` is pressed, or when the user chains a new operator after a complete expression. Calculations happen left to right, as on a physical desktop calculator.

### Values and entries

The display holds either an **entry** (characters the user is typing) or a **value** (a result from the service). The distinction drives editing:

- Digits and `.` start a new entry when the display holds a value; otherwise they append.
- `.` inserts at most once per entry; an empty entry becomes `0.`.
- Leading zeros are replaced, except the `0` before a decimal point.
- `⌫` removes the last character of an entry; on a computed value it is a no-op.
- `±` toggles the sign of an entry only; on a computed value it is a no-op. Sign-toggling zero is a no-op, so `-0` cannot be typed.
- `AC` clears everything (entry, pending operation, phase) and aborts any request in flight.

### Operators and equals

- With no pending operation, an operator (`+ − × ÷ xʸ`) stores the current value as the first operand, marks the operation pending, and shows it on the expression line (`12.5 +`). The display keeps the first operand until typing starts.
- Pressing another operator while the second operand is still empty **replaces** the pending operator; nothing is computed.
- Pressing an operator after a second operand is complete **executes** the pending operation, then starts the next one from the result (`2 + 3 +` → request → `5 +`).
- `=` with a complete expression executes it and enters result phase (`12.5 + 7.5 = 20`). `=` with an incomplete expression is a no-op — no error.
- Repeated `=` does not repeat the last operation. This is deliberately out of scope.
- After a result, typing a digit starts a new calculation; pressing an operator continues from the result.

### Unary operations

- `√` and `%` apply immediately to the displayed value, never to an empty second operand. If an operation is pending and no second operand has been entered, they are a no-op.
- When applied to a second operand, the pending operation is preserved: `12.5 + 9 √` shows the square-root result and the expression `12.5 + √(9)`; `=` then computes `12.5 + 3`.
- `xʸ` is a binary operator, not unary: it enters the operator rail like `×`.
- Applying `√` or `%` to a result is allowed and uses the displayed formatted value.

### Request phase

Exactly one request is in flight at a time.

- While calculating, the expression line shows `Calculating…`, the calculator sets `aria-busy="true"`, and every math key is ignored. `AC` remains available and aborts the request.
- A late response for an aborted or superseded request is discarded.
- No input is buffered while calculating; users see that the calculator is working and act after.

### Errors and recovery

Error copy appears in its own alert region directly beneath the display; the expression line keeps the failed expression. Errors never appear as a modal or toast, never move focus, and never prevent `AC`.

| Client error | Copy | State after |
| --- | --- | --- |
| `DIVISION_BY_ZERO` | Cannot divide by zero. | pending operation cleared; last value stays displayed |
| `NEGATIVE_SQUARE_ROOT` | Square root needs a number ≥ 0. | pending cleared; last value stays |
| `NON_FINITE_RESULT` | That result is too large or undefined. | pending cleared; last value stays |
| `INVALID_*`, `UNKNOWN_*` | That calculation isn't valid. | pending cleared (should be unreachable through the UI) |
| `NETWORK_ERROR` | Can't reach the calculator. Press = to retry. | expression retained; `=` retries |
| `TIMEOUT` | The calculator took too long. Press = to retry. | expression retained; `=` retries |
| `INTERNAL_ERROR` (5xx) | Something went wrong. Press = to retry. | expression retained; `=` retries |
| `INVALID_RESPONSE` | The calculator returned an unexpected response. Press = to retry. | expression retained; `=` retries |

Server codes, statuses, and schemas are canonical in [openapi.yaml](../elzzes-calculator-api/docs/openapi.yaml); this table defines user-facing copy and recovery.

- **After a domain error** (`DIVISION_BY_ZERO`, `NEGATIVE_SQUARE_ROOT`, `NON_FINITE_RESULT`, contract errors): the pending operation is cleared and the last displayed value stays. A digit or `.` starts a fresh calculation, an operator starts a new pending operation from the displayed value, `=` is a no-op, and `AC` resets. The failed expression and message remain visible until the next input.
- **After a transient error** (`NETWORK_ERROR`, `TIMEOUT`, `INTERNAL_ERROR`, `INVALID_RESPONSE`): the expression and pending operation are retained; `=` retries the same request. A digit abandons the failed expression and starts fresh, an operator replaces the pending operation, and `AC` resets.
- Domain failures are deterministic: retrying the same expression cannot succeed, so the state is cleared. Transient failures keep the expression so retry is one keypress.

A 10-second client timeout counts as a transient failure.

## Keyboard

The calculator is fully operable by keyboard. Shortcuts are conveniences, never the only path.

| Key | Action |
| --- | --- |
| `0`–`9`, `.` | enter digits / decimal point |
| `+`, `-`, `*`, `/` | add, subtract, multiply, divide |
| `^` | power (`xʸ`) |
| `%` | percentage |
| `r` | square root |
| `s` | toggle sign |
| `Enter` or `=` | equals |
| `Backspace` | backspace (entry only) |
| `Escape` or `Delete` | all clear (AC) |

- Shortcuts ignore events with `Ctrl`/`Meta`/`Alt` so browser shortcuts keep working.
- The global handler does not hijack `Enter`/`Space` when a button has focus; native button activation wins and prevents double-fire.
- Mapped keys call `preventDefault()` (e.g., `/` quick-find, `Backspace` history navigation).
- Buttons expose `aria-keyshortcuts` where a shortcut exists, so screen-reader users can discover them.

## Accessibility

- Every key is a real `<button>` with an accessible name: digits by value; operators as words (`divide`, `multiply`, `power`, `square root`, `percent`, `clear`, `backspace`, `toggle sign`).
- Visible focus ring with at least 3:1 contrast against adjacent colors; logical tab order equals DOM order; no focus traps, no focus moves on calculation or error.
- Minimum target 44×44 CSS px (56px keys default); page supports 200% zoom and increased text spacing without loss of function.
- The expression line is `aria-live="polite"`; it changes on operators, completion (`12.5 + 7.5 = 20`), and `Calculating…`, not on each digit. The display itself is not a live region, avoiding digit-by-digit announcements.
- Errors render in a separate element with `role="alert"`; starting a new calculation or pressing AC replaces/clears it.
- The calculator root exposes `aria-busy` while a request is in flight; keys are not individually `aria-disabled` — the whole calculation region is busy and math keys are no-ops.
- No information is conveyed by color alone; error text is explicit; `prefers-reduced-motion` disables all non-essential transitions.
- Language is `en`; the UI contains no user-generated or bidirectional text.

## Responsive behavior

- **≥ 480px:** centered card, max-width 420px, page padding 24px; keys 56px tall; display 40–48px type.
- **< 480px:** full-width card, 16px page padding; keys keep ≥ 48px height; display steps down; safe-area insets respected on notched devices.
- **Short landscape (< 560px height):** reduce vertical gaps and display height first; keypad remains fully visible where possible.
- No horizontal scrolling at 320px width; buttons remain at least 44px in both dimensions.
- Layout is mobile-first CSS; the desktop view is the same card with breathing room, not a stretched interface.

## Visual system

Tokens live in `@elzzes/ui` (folder `elzzes-ui-lib`) as Tailwind CSS custom properties; the web app does not hard-code colors.

| Role | Light | Dark | Use |
| --- | --- | --- | --- |
| Page background | `#f1f5f9` | `#0b1220` | behind the card |
| Surface | `#ffffff` | `#111827` | calculator card |
| Text | `#0f172a` | `#f8fafc` | display, digits |
| Muted text | `#475569` | `#cbd5e1` | expression line |
| Accent | `#1d4ed8` | `#60a5fa` | equals key, focus |
| Accent contrast | `#ffffff` | `#0b1220` | text on accent |
| Danger | `#b91c1c` | `#f87171` | AC tint, error copy |

- System font stack; display uses `font-variant-numeric: tabular-nums`.
- Spacing on a 4px scale (4/8/12/16/24/32); radii 8px (keys) and 12px (card); one soft card shadow in light mode, border-only separation in dark.
- Key hierarchy: fill accent for `=`; raised neutral for digits; subtly tinted neutral for operators; text-only accent for advanced keys; danger text for `AC`.
- Motion is limited to a 120ms press/focus transition, disabled under `prefers-reduced-motion`.
- Contrast is verified at AA for all text and focus indicators; values above are examples, not locked — tokens may be tuned as long as contrast tests hold.

## Components

`@elzzes/ui` (folder `elzzes-ui-lib`; presentational, controlled, calculator-agnostic):

- `Card` — generic surface with padding and semantic element choice.
- `Display` — expression/status live region, main value, error alert slot, busy state. Receives formatted strings.
- `Keypad` — reusable accessory strip plus a 4-column key grid; knows nothing about operations.
- `Button` — generic variants `primary`, `secondary`, `outline`, `ghost`, `accent`, `danger`; sizes `sm`, `md`, `lg`, `xl`; forwards accessible names and `aria-keyshortcuts`.
- Theme tokens — colors, spacing, radii, typography.
- Storybook stories — variants and states, documented in isolation ([UI library conventions](../elzzes-ui-lib/docs/CONVENTIONS.md)).

The web app maps calculator key roles to the generic variants: equals → `primary`, digits → `secondary`, binary operators → `outline`, advanced keys → `accent`, AC → `danger`, backspace/sign → `ghost`. The library contains no calculator vocabulary, copy, state, API, or arithmetic.

`elzzes-calculator-web` (logic):

- `Calculator` — composition and event wiring.
- `useCalculator` — reducer state machine and request orchestration via a TanStack Query mutation.
- `formatNumber`, entry rules, expression labels — pure functions.
- `calculate` API client, Zod schemas, and error mapping.

### Resolved implementation details

- Unary, backspace, and sign after an error are no-ops; digits/`.` start fresh, operators recover per the error tables, `=` retries only transient errors, and AC resets.
- A domain error freezes the last displayed text; a transient error retains the expression and the exact retry plan.
- Expression labels: standalone `√(x) = y`, `x% = y`; power uses `^`; a unary applied to a second operand keeps the pending expression (`12.5 + √(9)`).
- Accessible names for add, subtract, equals, and the decimal point are `add`, `subtract`, `equals`, and `decimal point`; shortcuts use the mapped keys (`Enter`, `Backspace`, `Escape Delete`, and so on).
- Backspacing to an empty entry displays `0`; the next digit starts a new entry.
- Results are rounded to 12 significant digits before the exponential threshold is applied; exponents use lowercase `e`, an explicit sign, and no zero padding.
- Request ids are monotonic across AC; continuations describe how a success lands (`equals`, `chain`, unary on the current value, unary on the second operand).
- Documented 4xx responses with a valid envelope map to the domain invalid-calculation copy; malformed or unknown responses map to the transient `INVALID_RESPONSE` copy.
- The API URL is built from the relative path and `window.location.origin`; no host literal is embedded.

## UI test scenarios

- Initial render: `0`, clear expression, all keys labelled.
- Entry: digits, leading-zero replacement, single decimal point, 15-digit cap, backspace, sign toggle, AC.
- Calculation: each of the seven operations; `=` with an incomplete expression is a no-op.
- Chaining: `2 + 3 +` issues one request and shows `5 +`; the resulting value then continues with a new operator.
- Operator replacement: `2 + ×` changes the pending operation without a request.
- Unary: on an entry, on a result, and no-op when a second operand is empty.
- In-flight lock: keys ignored; AC aborts; no late result overwrites state.
- Errors: division by zero, negative sqrt, overflow; message shown, alert announced, recovery via AC/new input.
- Transient errors: MSW network/timeout; expression retained; `=` retries successfully.
- Formatting: 12 significant digits, exponential thresholds, negative zero, WYSIWYG chaining.
- Keyboard parity for every shortcut and non-interference with focused buttons.
- Accessibility: accessible names, roles, live-region updates, `aria-busy`, focus visible.
- Responsive snapshots at 360, 768, 1280px; no horizontal overflow at 320px.

## Design decisions

1. **4×5 familiar keypad plus an advanced strip.** Keeps basic arithmetic where users expect it; advanced keys remain visible. Rejected: a 5-column scientific grid (clutter on phones) and a mode switch (hidden state, discoverability cost).
2. **Immediate-execution, no expression parser.** Matches the brief's "calculator" intent, avoids grammar/eval surface and precedence questions, and keeps client logic to a small reducer. Rejected: `eval`/expression parsing (security, complexity, out of scope).
3. **Unary percentage.** `x → x/100` is unambiguous and needs no expression context. Rejected: contextual percentage (`200 + 10% = 220`) — it silently changes semantics based on the pending operator and is a common source of confusion.
4. **Formatting on the client.** The service returns raw binary64; presentation rules are UI concerns and independently testable. Rejected: rounding in the API (loses information for other clients and mixes presentation into the domain).
5. **WYSIWYG chaining.** Use the displayed formatted value for the next operation. Rejected: carrying hidden raw precision (the user sees one number while the calculator uses another).
6. **Inline, non-blocking errors.** Messages live in the layout, announce themselves, and never trap input. Rejected: modals/toasts (interrupt flow, duplicate announcements, easy to miss).
7. **No history panel.** No persistence in scope. A client-side history could be added later without an API change; not built now.
8. **Split presentation from logic.** `@elzzes/ui` stays reusable and a11y-focused; `web` owns orchestration. Rejected: one package with mixed concerns (harder to test and evolve).
9. **Tailwind CSS v4 and one token source.** Utility styling plus `@elzzes/ui` Tailwind custom properties keeps visual rules consistent and testable. Rejected: app-level CSS Modules or bespoke CSS systems (two styling models, harder token control).
10. **Storybook for component documentation.** One story file per component covers variants and states and builds in CI. Rejected: treating Storybook as a substitute for behavior/accessibility tests.

See [spec 003](specs/003-calculator-ui/SPEC.md) for acceptance criteria and [Architecture](ARCHITECTURE.md) for the system around the UI.
