# Changelog

Notable changes, newest first. Implementation changes are added as they land; this file starts with the documentation baseline. Format loosely follows Keep a Changelog; dates are ISO 8601.

## [Unreleased] — 2026-11-20

### Added

- Documentation baseline derived from `problem.md`: root README (setup, operations, API examples and contract links, design decisions, assumptions, out of scope, prompts used).
- `AGENTS.md` working agreements: priorities, invariants, verification references, definition of done.
- `docs/PRINCIPLES.md`: standing invariants (arithmetic ownership, state ownership, local-first, finite results, stable contract) and decision rules.
- `docs/ARCHITECTURE.md`: system context, package boundaries, backend/frontend design, wire contract, numeric model, testing, Docker topology, trade-offs.
- `docs/DESIGN.md`: layout, display formatting, immediate-execution interaction, error recovery, keyboard map, accessibility, responsive behavior, design tokens.
- Specs 001–004 with index and traceability to the brief; spec 004 includes the reviewer checklist.
- Package READMEs for the API, web app, and shared UI library.
- Component conventions: [backend](../elzzes-calculator-api/docs/CONVENTIONS.md), [frontend](../elzzes-calculator-web/docs/CONVENTIONS.md), [UI library](../elzzes-ui-lib/docs/CONVENTIONS.md).
- Complete OpenAPI 3.1 contract in [`elzzes-calculator-api/docs/openapi.yaml`](../elzzes-calculator-api/docs/openapi.yaml).

### Changed

- Aligned the documentation with the agreed stack: pnpm workspaces, Tailwind CSS v4, TanStack Query, Zod, Storybook, OpenAPI 3.1, Nginx + Docker Compose, and GitHub Actions.
- Rewrote [`docs/DECISIONS.md`](DECISIONS.md) as a high-level decision index; detailed implementation rules moved to the component conventions.
- Made one document canonical per topic and replaced duplicated definitions with links: OpenAPI for the wire contract, specs 001–002 for operation behavior, DESIGN for UI behavior, ARCHITECTURE for structure/numeric model/deployment, spec 004 for verification.
- Updated README, AGENTS, package READMEs, INDEX, and specs to the canonical links and `cmd/api`/pnpm commands.

### Notes

- No calculator implementation yet; the repository has scaffolding only. `problem.md` remains authoritative; `brainstorm.md` is raw early notes, not a commitment.
