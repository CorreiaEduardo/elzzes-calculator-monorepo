# elzzes-ui-lib

Internal, unpublished React components, design tokens, and Storybook stories for Elzzes Calculator, published as `@elzzes/ui`. Presentational only: it renders what it is given and calls back on interaction.

**Status:** documentation baseline and scaffolding — visual and accessibility behavior is specified in [`docs/DESIGN.md`](../docs/DESIGN.md). Stack: React 19, TypeScript, Tailwind CSS v4, class-variance-authority, clsx, tailwind-merge, lucide-react, Storybook, Vitest + React Testing Library.

## Responsibilities

- The component set and variants defined in [Design § Components](../docs/DESIGN.md#components), plus theme tokens (Tailwind CSS custom properties).
- Light/dark tokens driven by `prefers-color-scheme`; responsive keypad layout.
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
| `pnpm --filter @elzzes/ui test` | Vitest |
| `pnpm --filter @elzzes/ui test:coverage` | tests + coverage report |

## Testing

Vitest + React Testing Library: rendering, variants, accessible names/roles, keyboard activation, focus visibility, and token usage. Storybook documents components; it does not replace behavior/accessibility tests. The coverage target and CI are canonical in [spec 004](../docs/specs/004-delivery-and-verification/SPEC.md).

## Conventions

The component contract, layout, accessibility obligations, Storybook rules, and boundaries are canonical in [`docs/CONVENTIONS.md`](docs/CONVENTIONS.md).

See [`docs/DESIGN.md`](../docs/DESIGN.md) for layout, interaction, tokens, and accessibility rules, and the [root README](../README.md) for the project overview.
