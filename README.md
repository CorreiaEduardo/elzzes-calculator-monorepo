# Elzzes Calculator

A full-stack calculator: a React + TypeScript single-page app backed by a stateless Go REST microservice that performs every calculation.

The frontend owns input, interaction, formatting, and error presentation; it never evaluates arithmetic. The stack is React 19 + TypeScript + Vite + Tailwind CSS v4 + TanStack Query + Zod on the front end, Go standard library on the back end, OpenAPI 3.1 as the wire contract, pnpm workspaces, Storybook for the shared UI library, and Nginx + Docker Compose for delivery. GitHub Actions runs the verification checks.

**Status:** specs 001–004 are implemented (2026-11-20): the Go API exposes all seven operations, the React app provides the calculator, typed API client, tests, and UI-library Storybook stories, and delivery includes Docker Compose, a root Makefile, and GitHub Actions CI. The manual responsive review is complete.

## How to run

Prerequisites: Docker with Compose v2 for the containerized path. For local development without Docker: Node.js 24+, pnpm 12.21.0, and Go 1.27+. The `Makefile` wraps the common commands; the equivalent package commands are shown where useful.

### Docker Compose

Builds and runs the SPA and API together on one origin:

```bash
docker compose up --build
```

- App: <http://localhost:8080>
- API health: <http://localhost:8080/healthz>

`web` builds the SPA with Vite and serves it from unprivileged Nginx on `:8080`; Nginx proxies `/api` and `/healthz` to `api:8080`. Both services are stateless; no database, secrets, or cloud account is required. Stop the stack with `docker compose down` (or `make compose-down`).

### Local development

```bash
make install   # pnpm install --frozen-lockfile + Go modules
make dev       # API on :8080 and Vite on :5173
```

Open <http://localhost:5173>. Vite proxies `/api` to the API on `:8080`, matching the production same-origin model.

To run each service in its own terminal:

```bash
make dev-api   # Go API: http://localhost:8080
make dev-web   # Vite: http://localhost:5173
```

Without `make`, run the same steps directly:

```bash
pnpm install --frozen-lockfile
(cd elzzes-calculator-api && go run ./cmd/api)
pnpm --filter elzzes-calculator-web dev
```

### Storybook

The `@elzzes/ui` component explorer is optional for running the calculator and serves on its own port:

```bash
make storybook   # http://localhost:6006
# equivalent: pnpm --filter @elzzes/ui storybook
```

### Verify a change

```bash
make verify          # format, vet, tests, coverage floors, race, lint, typecheck, builds
make compose-verify  # build the Compose stack, smoke-test /healthz and the API, then stop it
```

CI runs the same checks. The canonical procedure, coverage floor, and reviewer checklist are in [spec 004](docs/specs/004-delivery-and-verification/SPEC.md#verification-procedure).

## Operations

Seven operations, all through one endpoint:

- **Basic (mandatory):** `add`, `subtract`, `multiply`, `divide`.
- **Advanced (optional in the brief, included):** `power`, `sqrt`, `percent`.

Operation semantics, arity, and edge cases are canonical in [spec 001](docs/specs/001-basic-arithmetic-operations/SPEC.md) and [spec 002](docs/specs/002-advanced-arithmetic-operations/SPEC.md). The wire contract is [`openapi.yaml`](elzzes-calculator-api/docs/openapi.yaml); client behavior is in [Design](docs/DESIGN.md).

### API examples

```bash
curl -s localhost:8080/api/v1/calculate -H 'Content-Type: application/json' \
  -d '{"operation":"add","operands":[12.5,7.5]}'
# {"operation":"add","operands":[12.5,7.5],"result":20}

curl -s localhost:8080/api/v1/calculate -H 'Content-Type: application/json' \
  -d '{"operation":"sqrt","operands":[144]}'
# {"operation":"sqrt","operands":[144],"result":12}

curl -s localhost:8080/api/v1/calculate -H 'Content-Type: application/json' \
  -d '{"operation":"divide","operands":[10,0]}'
# {"error":{"code":"DIVISION_BY_ZERO","message":"cannot divide by zero"}}

curl -s localhost:8080/api/v1/calculate -H 'Content-Type: text/plain' -d '{}'
# {"error":{"code":"INVALID_CONTENT_TYPE","message":"Content-Type must be application/json"}}
```

## Repository layout

```text
.
├── elzzes-calculator-api/     # Go REST service — arithmetic and the OpenAPI contract
│   ├── docs/CONVENTIONS.md    # backend conventions
│   ├── docs/openapi.yaml      # canonical wire contract
│   └── Dockerfile
├── elzzes-calculator-web/     # React SPA — interaction, API client, formatting
│   ├── docs/CONVENTIONS.md    # frontend conventions
│   └── Dockerfile
├── elzzes-ui-lib/             # @elzzes/ui — presentational components + tokens
│   ├── docs/CONVENTIONS.md    # UI library conventions
│   └── stories/               # Storybook
├── docs/                      # architecture, design, decisions, principles, specs
├── .github/workflows/ci.yml   # GitHub Actions: same checks as make verify + Compose smoke
├── scripts/                   # check-go-coverage.sh (80% per-package floor)
├── Makefile                   # local commands: install, dev, storybook, verify, compose
├── compose.yaml               # runs the web and api services together
├── pnpm-workspace.yaml        # pnpm workspace for the two JS packages
├── AGENTS.md                  # working agreements for contributors and agents
├── problem.md                 # the brief; source of truth for requirements
└── brainstorm.md              # early notes, not authoritative
```

- **`elzzes-calculator-api`** — Go module. Domain arithmetic, HTTP API, OpenAPI contract. It does not serve the SPA.
- **`elzzes-calculator-web`** — React SPA. Keypad, calculator state machine, TanStack Query + Zod API client, formatting.
- **`elzzes-ui-lib`** — internal package published as `@elzzes/ui`. Controlled presentational components, design tokens, Storybook stories. No state, network, or arithmetic.
- **`docs/`** — principles, architecture, decisions, UI design, and the spec breakdown of `problem.md`.

The two JavaScript packages are pnpm workspaces; the Go module is independent. Each service has its own Dockerfile; `compose.yaml` runs them together. The `Makefile` is the local entry point (`make help` lists the targets) and CI invokes the same targets.

## Design decisions and assumptions

The high-level decision record, stack rationale, and rejected alternatives are in [docs/DECISIONS.md](docs/DECISIONS.md); system detail is in [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

Decisions that shape the app:

1. **One source of arithmetic truth.** All results come from the Go service. The browser parses input and formats output; it never computes.
2. **One calculate endpoint, discriminated by `operation`.** A fixed vocabulary with per-operation arity centralizes validation and keeps the contract small. The API is versioned under `/api/v1`.
3. **Stateless service, no database.** A calculation has no entity to persist. History and accounts were considered and deliberately deferred.
4. **Go standard library preferred.** `net/http` with method-aware routing, `encoding/json`, `log/slog`. The HTTP surface is small; a framework would add a dependency without removing complexity.
5. **IEEE-754 binary64 numbers.** Native JSON numbers, no third-party decimal dependency, and the same semantics as the browser. The trade-off is binary artifacts (`0.1 + 0.2`), hidden by display formatting ([Design § Display and formatting](docs/DESIGN.md#display-and-formatting)). Alternatives are weighed in [Architecture § Numeric model](docs/ARCHITECTURE.md#numeric-model).
6. **React reducer + TanStack Query + Zod.** The calculator has one small state machine; TanStack Query owns the request lifecycle (automatic retries disabled), and Zod validates responses at the wire boundary. No global state library.
7. **Tailwind + a presentational UI library.** `@elzzes/ui` holds controlled components, design tokens, and Storybook stories; the web app stays logic and composition.
8. **OpenAPI 3.1 as the contract.** One machine-readable source for paths, schemas, arity, statuses, and error codes; both layers pin to it.
9. **Nginx + Compose delivery.** Separate static and API services on one origin; no CORS, no Node runtime in production.
10. **Makefile for local commands.** `make dev`, `make storybook`, `make verify`, and `make compose-verify` wrap the package scripts so local development does not require Docker and CI runs the same definitions. The Makefile adds no application logic or dependency.
11. **GitHub Actions CI.** Go format/vet/test/race/build, web and UI lint/typecheck/tests/coverage, the SPA and Storybook builds, and a Compose smoke test on every push and pull request.

**Assumptions**

- Single-user, local/demo deployment: no authentication, authorization, or rate limiting.
- Percentage is unary (`x → x / 100`), not contextual (`a + b%`). The percentage key is unambiguous only this way.
- `0^0 = 1`, matching Go's `math.Pow` and common calculator behavior.
- Two-operand immediate execution (not an expression parser): pressing an operator with a complete pending expression evaluates left to right ([Design](docs/DESIGN.md)).
- Input and display limits are defined in [Design § Display and formatting](docs/DESIGN.md#display-and-formatting).
- English-only UI; light and dark follow the OS preference.
- The brief's "Dockerfile to run frontend + backend together" is read as the two service Dockerfiles plus Compose; no cloud account is required.
- Playwright E2E is an optional follow-up, not an acceptance criterion.

## Out of scope

Deliberately excluded to keep the deliverable focused ([brainstorm.md](brainstorm.md) notes were pruned against the brief):

- Persistence: database, calculation history, accounts, login ("logged environment").
- Cloud deployment: AWS/Terraform/emulators, image publishing, Kubernetes.
- Full scientific calculator: expression parsing/`eval`, memory keys, unit conversions, complex numbers, repeated-equals behavior.
- Internationalization and user-selectable themes.
- Playwright E2E as acceptance (optional follow-up).
- Extending the AI harness beyond this repo's docs and agent instructions.

## Documentation

| Document | Contents |
| --- | --- |
| [problem.md](problem.md) | The brief. Source of truth for requirements and constraints. |
| [docs/specs/INDEX.md](docs/specs/INDEX.md) | Spec breakdown, order, status, and traceability to the brief. |
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Packages, request flow, API narrative, numeric model, testing, deployment. |
| [docs/DECISIONS.md](docs/DECISIONS.md) | High-level stack and architecture decisions; index of canonical details. |
| [docs/DESIGN.md](docs/DESIGN.md) | Layout, interaction model, formatting, errors, keyboard, accessibility, responsive behavior. |
| [docs/PRINCIPLES.md](docs/PRINCIPLES.md) | Standing invariants and how new work is approached. |
| [docs/CHANGELOG.md](docs/CHANGELOG.md) | Notable changes, newest first. |
| [docs/specs/001–004](docs/specs/INDEX.md) | Operation behavior/edge cases (001–002), UI acceptance (003), delivery/verification (004). |
| [elzzes-calculator-api/docs/openapi.yaml](elzzes-calculator-api/docs/openapi.yaml) | Canonical wire contract. |
| [elzzes-calculator-api/docs/CONVENTIONS.md](elzzes-calculator-api/docs/CONVENTIONS.md) | Backend implementation conventions. |
| [elzzes-calculator-web/docs/CONVENTIONS.md](elzzes-calculator-web/docs/CONVENTIONS.md) | Frontend implementation conventions. |
| [elzzes-ui-lib/docs/CONVENTIONS.md](elzzes-ui-lib/docs/CONVENTIONS.md) | UI library implementation conventions. |
| [AGENTS.md](AGENTS.md) | Working agreements for humans and coding agents in this repo. |
| [brainstorm.md](brainstorm.md) | Raw early notes; not authoritative. |

## Prompts used

- [`problem.md`](problem.md) is the challenge brief (the tech assessment itself) and the lighthouse for every trade-off. When a spec or doc disagrees with it, the brief wins.
- [`brainstorm.md`](brainstorm.md) holds my early raw notes. Ideas beyond the brief (database, AWS emulator, auth) are recorded as out of scope rather than silently dropped.
- [`001-extending-harness.md`](docs/prompts/001-extending-harness.md) a transcripted spoken brainstorm adjusted with a clarity focused prompt with Claude Code
- [`002-creating-specs.md`](docs/prompts/002-creating-specs.md) the dictated brainstorm that turned `problem.md` into specs 001–004 before implementation
- [`003-scaffolding-go.md`](docs/prompts/003-scaffolding-go.md) the scaffold prompt for the Go module, Ports & Adapters layout, config, logging, and `/healthz`
- [`004-implementing-backend.md`](docs/prompts/004-implementing-backend.md) the prompt that implemented specs 001–002, including the transport-to-calculator boundary correction
- [`005-implementing-frontend.md`](docs/prompts/005-implementing-frontend.md) the prompt that implemented the calculator UI and UI library
- [`006-infra-and-deliverables`](docs/prompts/006-infra-and-deliverables) the prompt that closed spec 004: CI, Makefile, coverage enforcement, and README run instructions
