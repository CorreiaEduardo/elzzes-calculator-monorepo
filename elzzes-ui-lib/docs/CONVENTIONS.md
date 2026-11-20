# UI library conventions

**Status:** implementation baseline, 2026-11-20. Canonical for `@elzzes/ui` (folder `elzzes-ui-lib`) implementation conventions. Visual and accessibility behavior is canonical in [Design](../../docs/DESIGN.md); acceptance is in [spec 003](../../docs/specs/003-calculator-ui/SPEC.md).

## Stack

React 19, TypeScript, Tailwind CSS v4, class-variance-authority, clsx, tailwind-merge, lucide-react, Storybook, Vitest + React Testing Library. Package manager: pnpm workspaces.

## Layout

```text
elzzes-ui-lib/
├── src/
│   ├── components/
│   │   ├── button/             # generic Button + variants/sizes
│   │   ├── card/               # generic surface
│   │   ├── display/            # value, expression live region, alert, busy
│   │   └── keypad/             # accessory strip + key grid
│   ├── lib/cn.ts               # clsx + tailwind-merge
│   ├── test/setup.ts
│   ├── styles.css              # Tailwind entry and theme tokens
│   └── index.ts                # public exports
├── stories/                    # one story file per component
├── .storybook/                 # Storybook config
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
- Components are generic primitives (`Button`, `Card`, `Display`, `Keypad`), not calculator objects. The web app maps semantic roles to generic variants and passes labels, values, and callbacks.
- Forward accessible names and `aria-keyshortcuts`; keep the alert and busy slots controlled via props.
- Theme tokens use the `--color-ui-*` namespace in `styles.css` (`ui-page`, `ui-surface`, `ui-ink`, `ui-muted`, `ui-accent`, `ui-danger`, `ui-key`, `ui-operator`, `ui-border`, and hover/soft variants), plus `--radius-ui-*` and `--shadow-ui-card`; dark values override under `prefers-color-scheme`.
- No calculator state, API calls, arithmetic, routing, or persistence.
- No knowledge of the operation vocabulary, calculator key roles, or user-facing copy; the web layer owns behavior and copy.

## Accessibility

Follow [Design § Accessibility](../../docs/DESIGN.md#accessibility). At the component level:

- Real `<button>` elements, accessible names, visible focus, and the minimum target size from [Design § Accessibility](../../docs/DESIGN.md#accessibility).
- `aria-keyshortcuts` where a shortcut exists; no information conveyed by color alone.
- Components expose the slots and layout regions defined in [Design § Components](../../docs/DESIGN.md#components), including the alert and busy state.
- `prefers-reduced-motion` disables non-essential transitions.

## Storybook

- `.storybook/main.ts` uses `@storybook/react-vite`; `.storybook/preview.ts` imports `src/styles.css`.
- One story file per component in `stories/`, covering variants and states.
- Stories run independently of the web app (`pnpm --filter @elzzes/ui storybook`) and build with `build-storybook`.
- Storybook documents components; it does not replace behavior or accessibility tests.

## Testing

- Vitest + React Testing Library + user-event: rendering, variants, accessible names/roles, keyboard activation, focus visibility, token usage, live/alert/busy states, accessory slot.
- Tests live beside components; the Vitest threads pool is used because the forks pool hangs in this WSL environment.
- No test depends on the web app or a live API.
- Coverage: `pnpm --filter @elzzes/ui test:coverage` (v8 text + HTML under `coverage/`).

## Scripts

| Command | Purpose |
| --- | --- |
| `pnpm --filter @elzzes/ui storybook` | Storybook dev server |
| `pnpm --filter @elzzes/ui build-storybook` | static Storybook build |
| `pnpm --filter @elzzes/ui typecheck` | TypeScript, no emit |
| `pnpm --filter @elzzes/ui lint` | ESLint |
| `pnpm --filter @elzzes/ui test` | Vitest (run once) |
| `pnpm --filter @elzzes/ui test:watch` | Vitest watch mode |
| `pnpm --filter @elzzes/ui test:coverage` | tests + coverage report |

## Boundaries

- No fetch, API types, arithmetic, application state, or routing.
- No calculator operation names, key-role vocabulary, or user-facing copy; pass labels and values in from the web layer.
- No dependency on `elzzes-calculator-web`.
- Public exports only through `src/index.ts`; no deep imports from consumers.
