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
- Go API scaffold: `go.mod` (Go 1.27, stdlib only), `cmd/api` composition root, lightweight Ports & Adapters layout (`internal/calculator` core placeholder, `internal/transport/http` net/http adapter with a consumer-owned `Logger` port and middleware, `internal/logger` slog adapter, `internal/config` environment loader), structured request logging, panic recovery, JSON `404`/`405` routing with `Allow`, config validation, graceful shutdown, and a working `GET /healthz`.
- Frontend scaffold: root pnpm workspace (`pnpm-workspace.yaml` + lockfile); `elzzes-calculator-web` with React 19, TypeScript, Vite, Tailwind CSS v4, ESLint/Prettier/EditorConfig configs, and the Vite starter screen rendering `@elzzes/ui`; `elzzes-ui-lib` (`@elzzes/ui`) with a Vite library config, Tailwind `@theme` color tokens, and a controlled `Button` component (variants and sizes) exported with its types.
- Web delivery: multi-stage `elzzes-calculator-web/Dockerfile` (pnpm workspace install, Vite build, unprivileged Nginx on `:8080`), `nginx.conf` with SPA history fallback and `/api`/`/healthz` proxies, root `.dockerignore`, and `compose.yaml` running `web` and `api` on one origin.
- Calculator backend for specs 001–002: pure core with the seven operations, one operation/arity registry, and typed domain errors; `POST /api/v1/calculate` with a 4 KiB body cap, JSON content-type and strict-field validation, binary64 operand checks, finite-result enforcement, and centralized domain-error mapping.
- Backend tests: domain matrix for basic and advanced operations (division by zero, overflow, underflow, negative zero, power/sqrt/percent domains) and HTTP matrix for success, 400/413/415/422, 405/404, body limits, and unexpected errors.
- Multi-stage non-root `Dockerfile` and `.dockerignore` for the API. Docker was unavailable locally, so the image was written but not built.
- Table-driven tests for config loading, the logging adapter, routing `404`/`405`, middleware, panic recovery, the error envelope, and health.

### Changed

- Aligned the documentation with the agreed stack: pnpm workspaces, Tailwind CSS v4, TanStack Query, Zod, Storybook, OpenAPI 3.1, Nginx + Docker Compose, and GitHub Actions.
- Rewrote [`docs/DECISIONS.md`](DECISIONS.md) as a high-level decision index; detailed implementation rules moved to the component conventions.
- Made one document canonical per topic and replaced duplicated definitions with links: OpenAPI for the wire contract, specs 001–002 for operation behavior, DESIGN for UI behavior, ARCHITECTURE for structure/numeric model/deployment, spec 004 for verification.
- Updated README, AGENTS, package READMEs, INDEX, and specs to the canonical links and `cmd/api`/pnpm commands.
- Rewrote the backend conventions and the architecture's backend section around the lightweight Ports & Adapters design; added consumer-owned-ports and plain-names invariants to PRINCIPLES and the decision index. The API README now documents the runnable scaffold, health smoke, and Docker build.
- Corrected the backend layer rule: `transport/http` may depend on the calculator's public API, while `calculator` must not depend on transport or HTTP; updated conventions, architecture, and decisions.
- Updated API/root READMEs and spec statuses for the implemented backend, and corrected spec 001's underflow example to `5e-324 ÷ 10 → 0` (`1e-308 ÷ 1e10` is `1e-318`).
- Refined calculate validation: field-type errors now use `UNKNOWN_OPERATION`/`INVALID_OPERAND` instead of the misleading `INVALID_JSON`; operand errors name their index; duplicate keys are documented as last-value-wins; non-limit body-read failures are logged for diagnostics.

### Notes

- The Go backend implements all seven operations; the web app and UI library remain scaffolds pending spec 003. `problem.md` remains authoritative; `brainstorm.md` is raw early notes, not a commitment.
- Web and UI packages are scaffolds: the web app still renders the Vite starter screen, and the calculator UI, API client, and Storybook stories are not implemented.
- The frozen pnpm install and SPA production build were verified locally; the web image and Compose stack were not built because Docker is unavailable in this environment.
- Backend verification: `gofmt`, `go vet`, `go test`, `go test -race`, and the full curl smoke matrix pass. Docker/Compose remain unverified.
