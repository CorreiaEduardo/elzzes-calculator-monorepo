# elzzes-ui-lib

Internal, unpublished React components, design tokens, and Storybook stories for Elzzes Calculator, published as `@elzzes/ui`. Presentational only: it renders what it is given and calls back on interaction.

**Status:** implemented — generic UI primitives, variants, theme tokens, tests, and Storybook stories. Visual and accessibility behavior is specified in [`docs/DESIGN.md`](../docs/DESIGN.md). Stack: React 19, TypeScript, Tailwind CSS v4, class-variance-authority, clsx, tailwind-merge, Storybook, Vitest + React Testing Library.

## Responsibilities

- Generic primitives: `Button` (six variants, four sizes), `Card`, `Display` (expression live region, value, alert, busy), and `Keypad` (accessory strip plus key grid).
- Theme tokens as Tailwind CSS custom properties, with light/dark values driven by `prefers-color-scheme`.
- Accessible primitives: real buttons, accessible names, visible focus, minimum target size, `aria-keyshortcuts`.
- Storybook stories covering component variants and states.

## Non-responsibilities

No fetch or API knowledge, no arithmetic, no calculator state, no operation vocabulary, no routing, no persistence. Controlled by props; state stays in `elzzes-calculator-web`.

## Usage

Consumed as a pnpm workspace package by `elzzes-calculator-web`:

```bash
pnpm install                              # repo root
pnpm --filter @elzzes/ui storybook        # component docs, http://localhost:6006
pnpm --filter @elzzes/ui test
```

## Scripts

| Command | Purpose |
| --- | --- |
| `pnpm --filter @elzzes/ui storybook` | Storybook dev server |
| `pnpm --filter @elzzes/ui build-storybook` | static Storybook build |
| `pnpm --filter @elzzes/ui lint` | ESLint |
| `pnpm --filter @elzzes/ui typecheck` | TypeScript, no emit |
| `pnpm --filter @elzzes/ui test` | Vitest (run once) |
| `pnpm --filter @elzzes/ui test:watch` | Vitest watch mode |
| `pnpm --filter @elzzes/ui test:coverage` | tests + coverage report |

## Testing

Vitest + React Testing Library + user-event: rendering, variants, accessible names/roles, keyboard activation, focus visibility, token usage, live/alert/busy states, and the accessory slot. Storybook documents components; it does not replace behavior/accessibility tests. The coverage target and CI are canonical in [spec 004](../docs/specs/004-delivery-and-verification/SPEC.md).

## Conventions

The component contract, layout, accessibility obligations, Storybook rules, and boundaries are canonical in [`docs/CONVENTIONS.md`](docs/CONVENTIONS.md).

See [`docs/DESIGN.md`](../docs/DESIGN.md) for layout, interaction, tokens, and accessibility rules, and the [root README](../README.md) for the project overview.
