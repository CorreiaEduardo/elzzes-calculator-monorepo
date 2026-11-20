# 004 — Delivery and verification

**Status:** implemented (2026-11-20). Depends on specs 001–003. Deliverable requirements come directly from `problem.md`; the agreed delivery stack is Docker Compose + Nginx + GitHub Actions plus a root Makefile for local development ([DECISIONS](../../DECISIONS.md)).

## What is it?

The non-functional deliverables that make the app reproducible and reviewable: a Dockerized run of frontend + backend together, a Makefile for running and verifying the app without Docker, unit tests with coverage reports on both layers, CI, and complete documentation (setup, run, API examples, design decisions, assumptions).

## Scope

**In**

- Multi-stage Dockerfiles for the web and API services, orchestrated by `compose.yaml`.
- Nginx serving the SPA and proxying `/api/` and `/healthz` to the Go API on one origin.
- A root `Makefile` for local development without Docker: install, dev servers, Storybook, lint, build, typecheck, test, coverage, and verification.
- Reproducible local commands and scripts for lint, build, typecheck, test, and coverage.
- Coverage reports for Go and both JS packages; documented target.
- GitHub Actions CI running the same checks on push and pull requests.
- README completeness against the brief; OpenAPI and docs kept current.
- A local verification checklist a reviewer can run end to end.

**Out**

- Registry publishing, cloud deployment, infrastructure as code, Kubernetes.
- Automated deployment or release versioning.
- E2E browser automation as acceptance (Playwright is an optional follow-up).
- Performance/load testing, security auditing beyond the posture in the architecture.

## Requirements

### Docker

1. `docker compose up --build` succeeds from a clean clone and serves the SPA at <http://localhost:8080> and the API at `/api/v1/calculate` on the same origin.
2. `elzzes-calculator-web/Dockerfile` is multi-stage: install with pnpm, build the SPA, and copy the static output plus the Nginx config into an unprivileged Nginx image listening on `8080`.
3. `elzzes-calculator-api/Dockerfile` is multi-stage: build `./cmd/api` with the Go toolchain, then run the static binary as a non-root user in a minimal image listening on `8080`.
4. Nginx serves the SPA with history fallback and proxies `/api/` and `/healthz` to `api:8080`. No CORS is required.
5. Final images contain no dev server, no Node or Go toolchain, no source, and no secrets. Processes run as non-root where the base image permits.
6. The stack requires no external dependencies, volumes, database, or network beyond the published port.
7. `GET /healthz` returns `200 {"status":"ok"}` through the API and through Nginx.
8. Builds are deterministic: pnpm lockfile and Go modules are used; `.dockerignore` files exclude `.git`, `node_modules`, build output, coverage, and local env files.
9. The brief's "Dockerfile to run frontend + backend together" is satisfied by the two service Dockerfiles plus Compose; `docker compose up` runs both. This reading is recorded in the README.

### Local development

1. A root `Makefile` is the local command surface; it orchestrates the existing package scripts and adds no application logic, tools, or dependencies.
2. `make install` runs `pnpm install --frozen-lockfile` and `go mod download`.
3. `make dev-api` runs the API with `go run ./cmd/api` on `:8080`; `make dev-web` runs the Vite dev server on `:5173` with `/api` proxied to `:8080`; `make dev` runs both and stops both on exit.
4. `make storybook` runs the UI library's Storybook on `:6006`.
5. `make verify` runs every non-Docker check in this spec in one command. `make compose-verify` builds the stack, checks `/healthz`, the SPA, the calculate error path, and stops it.
6. The README's **How to run** section is the canonical local-run guide and covers Docker, the Makefile, manual commands, and Storybook.

### Tests and coverage

1. Go: `go vet ./...` clean; `go test ./...` green; `go test ./... -coverprofile=coverage.out` produces a report; `go test -race ./...` green in CI.
2. Web and UI library: lint, typecheck, tests, and coverage pass, with text and HTML reports under `coverage/`; Storybook builds for the UI library.
3. Every operation and every edge case in specs 001–003 has at least one test; the edge-case matrices are the source of test names.
4. No test requires network access or sleeps; UI HTTP is mocked (MSW).
5. Coverage target: **≥ 80% statements/lines per package**, with documented edge cases fully covered regardless of the percentage. The floor is enforced, not only reported: Vitest thresholds for the web and UI packages, and `scripts/check-go-coverage.sh` for each Go `internal/...` package.
6. `make verify` runs the same checks as CI (without Docker) and fails on any of them.

### CI

1. GitHub Actions runs on push and pull request and covers the same checks as the local verification procedure below, including the Compose smoke test.
2. CI uses `pnpm install --frozen-lockfile` and the Go module cache; it fails on lint, typecheck, test, coverage, or build failures.
3. CI builds Storybook and the production frontend, and runs `go vet`, `go test -race`, and `go build ./cmd/api`.
4. CI does not deploy, publish images, or require repository secrets beyond the defaults.
5. Playwright is not part of required CI; an optional workflow may be added later.
6. CI invokes the same `Makefile` targets as local verification so the command definitions do not drift.

### Documentation

README must include, as required by `problem.md`:

- A prominent **How to run** section near the top, covering Docker Compose, local development with the Makefile, manual commands, and verification.
- Setup prerequisites and instructions.
- How to run frontend and backend locally and with Compose.
- A **Storybook** subsection explaining how to start it and where it is served.
- REST API examples with success and error responses, matching the OpenAPI contract.
- Design decisions and assumptions, including out-of-scope items.
- Test/coverage/CI targets and links to the canonical verification procedure.
- Documentation map and prompts used.

Additional docs stay consistent: `ARCHITECTURE.md` canonical for structure/numeric model/deployment; `openapi.yaml` canonical for the wire contract; `DESIGN.md` canonical for UI behavior; specs reflect implemented behavior; component conventions hold implementation detail; `CHANGELOG.md` has notable entries.

## Acceptance criteria (reviewer checklist)

- [X] Clean clone + `make install` + `make dev` runs the SPA on `:5173` and the API on `:8080`; Storybook runs with `make storybook`.
- [X] Clean clone + `pnpm install` + `go run ./cmd/api` + `pnpm --filter elzzes-calculator-web dev` works.
- [X] `docker compose up --build` works; SPA and API reachable on `:8080`.
- [X] `GET /healthz` succeeds through Nginx.
- [X] `make verify` passes (format, vet, tests, coverage floors, race, lint, typecheck, builds).
- [X] `go vet`, `go test`, and `go test -race` pass; coverage report generated and ≥ target per package.
- [X] Web and UI library lint/typecheck/test/coverage pass; Storybook builds.
- [ ] CI is green on the change.
- [X] All README-required sections present, including a prominent How to run and a Storybook subsection; API examples reproduce actual output and match OpenAPI.
- [X] Every operation accessible from the UI and via `curl`.
- [X] Error paths show documented codes/statuses and friendly UI copy.
- [X] No database, auth, cloud service, or unreviewed dependency was added.
- [X] Docs match behavior and one canonical source per topic; CHANGELOG records notable changes.

## Verification procedure

`make verify` runs sections 1 and 2; `make compose-verify` runs section 3 and stops the stack. `make storybook` serves Storybook at <http://localhost:6006>. The explicit commands remain canonical:

```bash
# 1. API verification
cd elzzes-calculator-api
gofmt -l .
go vet ./...
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
./scripts/check-go-coverage.sh coverage.out
go test -race ./...
go run ./cmd/api &
curl -s localhost:8080/healthz
curl -s localhost:8080/api/v1/calculate -H 'Content-Type: application/json' \
  -d '{"operation":"add","operands":[12.5,7.5]}'
kill %1

# 2. Web and UI verification (repo root)
pnpm install --frozen-lockfile
pnpm --filter elzzes-calculator-web lint
pnpm --filter elzzes-calculator-web typecheck
pnpm --filter elzzes-calculator-web test
pnpm --filter elzzes-calculator-web test:coverage
pnpm --filter @elzzes/ui lint
pnpm --filter @elzzes/ui typecheck
pnpm --filter @elzzes/ui test
pnpm --filter @elzzes/ui test:coverage
pnpm --filter @elzzes/ui build-storybook

# 3. Compose verification (repo root)
docker compose up --build -d
curl -s localhost:8080/healthz
curl -s localhost:8080/ | head -n 5
curl -s localhost:8080/api/v1/calculate -H 'Content-Type: application/json' \
  -d '{"operation":"divide","operands":[10,0]}'
docker compose down
```

Manual UI pass: each operation, chaining, AC, division by zero, negative sqrt, overflow, network failure (stop the API), keyboard-only pass, 360px viewport.

## Decision

Two service images (unprivileged Nginx for the SPA, minimal Go runtime for the API) orchestrated by Compose; GitHub Actions CI; a root Makefile for local development and verification without Docker; coverage as an enforced guardrail at 80% per package; Storybook as a UI-library deliverable; Playwright optional.

## Reasons

- Separate images keep static serving in Nginx and arithmetic in Go, and match the chosen stack; Nginx also provides the same-origin API proxy.
- Compose is the standard way to run both services with one command while preserving independent images.
- The Makefile keeps the development loop native and fast (Go's toolchain plus Vite HMR) and gives reviewers one command for the checks CI runs, without adding a build system or dependency.
- CI catches contract drift across OpenAPI, Go DTOs, and Zod schemas; invoking the Makefile targets keeps CI and local commands identical.
- A percentage floor catches untested files, while the edge-case matrices ensure the mandated behaviors are tested even when coverage looks high.

## Rejected alternatives

- **Single Go image embedding the SPA:** fewer units, but couples every frontend build to the API image and moves static serving into Go; rejected in favor of explicit Nginx and API services.
- **Serving the SPA from a Node server:** keeps a runtime and dependency surface the production app does not need.
- **Docker-only local workflow:** rebuilding images for every edit slows the development loop; the Makefile runs the native toolchain and leaves Docker for packaging and review.
- **Makefile-driven Docker builds:** the Dockerfiles stay self-contained and buildable without Make; the Makefile is only a developer convenience.
- **Checked-in coverage reports:** reports are generated artifacts; commands are documented instead.
- **Registry publishing / cloud deployment / Kubernetes:** beyond the brief and would add credentials and maintenance.
- **Playwright E2E as acceptance:** useful follow-up but not required by the brief; it would expand scope and dependencies.
- **Skipping CI:** the stack includes GitHub Actions; local-only verification would make regressions easier to miss.

## Open questions

None. Registry publishing, cloud deployment, and Playwright acceptance are candidates for a follow-up change, each requiring agreement.
