# UI library conventions

**Status:** planned baseline, 2026-11-20. Canonical for `@elzzes/ui` (folder `elzzes-ui-lib`) implementation conventions. Visual and accessibility behavior is canonical in [Design](../../docs/DESIGN.md); acceptance is in [spec 003](../../docs/specs/003-calculator-ui/SPEC.md).

## Stack

React 19, TypeScript, Tailwind CSS v4, class-variance-authority, clsx, tailwind-merge, lucide-react, Storybook, Vitest + React Testing Library. Package manager: pnpm workspaces.

## Layout

```text
elzzes-ui-lib/
├── src/
│   ├── components/
│   │   ├── button/             # CalcButton + variants
│   │   ├── display/
│   │   ├── keypad/
│   │   └── calculator-frame/
│   ├── styles.css              # Tailwind entry and theme tokens
│   └── index.ts                # public exports
├── stories/                    # Storybook stories
├── .storybook/
├── docs/CONVENTIONS.md
├── package.json
├── tsconfig.json
└── vite.config.ts
```

Expose components through `src/index.ts`; consumers do not deep-import internals.

## Component contract

- Controlled and presentational: receive formatted strings, state flags, and callbacks; render what they are given.
- Variants via `class-variance-authority`; merge conditional classes with `clsx`/`tailwind-merge`.
- Theme tokens via Tailwind CSS custom properties in `styles.css`. No hard-coded colors, spacing, or radii in components.
- Component names, variants, and slots follow [Design § Components](../../docs/DESIGN.md#components); forward accessible names and `aria-keyshortcuts`.
- No calculator state, API calls, arithmetic, routing, or persistence.
- No knowledge of the operation vocabulary; the web layer owns behavior and copy.

## Accessibility

Follow [Design § Accessibility](../../docs/DESIGN.md#accessibility). At the component level:

- Real `<button>` elements, accessible names, visible focus, and the minimum target size from [Design § Accessibility](../../docs/DESIGN.md#accessibility).
- `aria-keyshortcuts` where a shortcut exists; no information conveyed by color alone.
- Components expose the slots and layout regions defined in [Design § Components](../../docs/DESIGN.md#components), including the alert and busy state.
- `prefers-reduced-motion` disables non-essential transitions.

## Storybook

- Storybook is part of the package deliverable: one story file per component, covering variants and states.
- Stories run independently of the web app (`pnpm --filter @elzzes/ui storybook`) and build in CI (`build-storybook`).
- Storybook documents components; it does not replace behavior or accessibility tests.

## Testing

- Vitest + React Testing Library: rendering, variants, accessible names/roles, keyboard activation, focus visibility, token usage.
- No test depends on the web app or a live API.
- Coverage target and commands: [spec 004](../../docs/specs/004-delivery-and-verification/SPEC.md).

## Scripts

| Command | Purpose |
| --- | --- |
| `pnpm --filter @elzzes/ui storybook` | Storybook dev server |
| `pnpm --filter @elzzes/ui build-storybook` | static Storybook build |
| `pnpm --filter @elzzes/ui typecheck` | TypeScript, no emit |
| `pnpm --filter @elzzes/ui lint` | ESLint |
| `pnpm --filter @elzzes/ui test` | Vitest |
| `pnpm --filter @elzzes/ui test:coverage` | tests + coverage report |

## Boundaries

- No fetch, API types, arithmetic, application state, or routing.
- No calculator operation names or user-facing copy; pass labels and values in from the web layer.
- No dependency on `elzzes-calculator-web`.
