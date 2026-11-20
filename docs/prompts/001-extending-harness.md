## Extending the harness

> This is a transcribed version of a brainstorming session that I dictated to an AI, asking it to summarize and better structure the content into a prompt.

I’m building the a calculator as a monorepo and I want the documentation nailed down before
any real implementation. The brief is in problem.md and that’s the lighthouse; if anything I
say here or any spec disagrees with it, the brief wins. The repo has three packages: a Go REST
API, a React + TypeScript web app, and an internal presentational UI library. I want a docs
baseline for all of it: README, AGENTS, principles, architecture, decisions, design, changelog,
a specs index plus specs, package READMEs, component conventions, and a real OpenAPI
contract. No calculator code yet. scaffolding only, and the docs should say that plainly.

The product is a stateless calculator. The React app owns input, UI state, formatting, and error
presentation, but every result comes from the Go service. No client-side arithmetic and no
eval. Seven operations: add, subtract, multiply, divide, power, square root, and percentage. The
four basic ones are mandatory from the brief; the other three are optional in the brief but I
want them included. One endpoint, POST /api/v1/calculate, discriminated by operation. Binary
operations take exactly two operands; sqrt and percent take exactly one. The response echoes the
operation and operands and returns the raw result. Errors use one envelope with code and
message, where the message is developer-facing and the client owns user copy. No database, auth,
history, or persistence. Single-user local/demo scope. The service is stateless and restarting
loses nothing.

For the stack, I’m deliberately going with pnpm workspaces for the JavaScript packages and an
independent Go module. React 19, TypeScript, Vite, Tailwind CSS v4, TanStack Query, and Zod on
the front end. Storybook for the UI library. Prioritize Go standard library on the back end: net/http,
encoding/json, log/slog, but may use fiber and zap. OpenAPI 3.1 is the wire contract. Delivery is two service Dockerfiles
orchestrated by Docker Compose, with Nginx serving the SPA and proxying /api/ and /healthz to
the Go API on one origin. GitHub Actions runs CI. Playwright is optional, not acceptance.
Vitest, React Testing Library, user-event, and MSW are the front-end testing tools. No extra
dependency or abstraction without a requirement and agreement.

The folder layout matters. At the root I want elzzes-calculator-api, elzzes-calculator-web,
elzzes-ui-lib, docs, compose.yaml, pnpm-workspace.yaml, AGENTS.md. The API has
cmd/api, internal/features/calculator/domain, internal/features/calculator/http,
internal/features/commons/errors, internal/config, docs/openapi.yaml, docs/CONVENTIONS.md, and a
Dockerfile. The web app has src/app, src/features/calculator with components, hooks, api, schemas,
and types, plus src/lib, src/styles, src/test, docs/CONVENTIONS.md, and a Dockerfile.
The UI library has src/components, src/styles.css, src/index.ts, stories, .storybook,
and docs/CONVENTIONS.md. The API should never serve the SPA. Vite serves it in development 
and proxies /api to localhost:8080; Nginx serves the built SPA in production and proxies to api:8080.
//todo validation
Architecturally, the Go service has a pure domain layer and an http adapter. The adapter does
routing, strict decoding, validation, status/code mapping, and health. I don’t want ports,
services, or mocks layers for this — there is no external dependency to abstract, and the
conventions should say that directly. Strict JSON with unknown fields rejected, content-type
enforcement, the body cap from OpenAPI, panic recovery, one structured request log line,
explicit timeouts, graceful shutdown, and JSON 404/405 for API paths. No CORS because production
is same-origin through Nginx and development uses the Vite proxy. GET /healthz returns
{"status":"ok"}. The OpenAPI file is canonical for the wire contract; the architecture doc
should explain the shape and link to it instead of duplicating every schema.

The principles should be short and standing: one source of arithmetic truth, explicit state
ownership, local-first, finite results or explicit errors, a stable OpenAPI contract, validate
at the boundary, errors are data, statelessness, accessibility and responsiveness as
requirements, a simplicity budget, docs travel with the change, and lean comments. Decision
rules: fewer moving parts, write down cross-layer choices before building, take the narrowest
reading of an ambiguous brief, prefer deleting a feature over carrying an untested path, and
treat a duplicated definition as a defect that should be moved to its canonical home and linked.

The conventions need to live in the component folders. Backend: Go naming, table-driven tests,
httptest, race detector in CI, gofmt, go vet, golangci-lint, standard library over assertion
libraries, structured logging, no third-party router. Frontend: one reducer for interaction
state, TanStack Query mutation for the request lifecycle with automatic retries off,
AbortController on AC, Zod response validation, relative /api/v1/calculate, error mapping,
Tailwind utility styling with tokens from the UI library, no global state library. UI library:
controlled components, variant styling, Tailwind custom properties, accessible primitives,
Storybook, no API or calculator state. These are implementation details, so they belong in the
component docs, not in the high-level decisions file.

Documentation organization is important to me. One canonical source per topic, and everywhere
else links to it. OpenAPI owns the wire contract. Specs own operation behavior,
edge cases, verification, coverage, and CI. DESIGN.md owns UI behavior, formatting, error copy, keyboard, accessibility,
responsive behavior, and visual tokens. ARCHITECTURE.md owns system structure, request flow,
state ownership, the numeric model, and deployment. The component CONVENTIONS.md files own
implementation detail. DECISIONS.md should stay
high-level — stack, boundaries, rejected alternatives, and a table pointing at the canonical
details; it should not become a second specification. README is the deliverable and the
documentation map: setup, run, API examples, design decisions and assumptions, out of scope,
prompts used, and links. AGENTS is the working agreement and should point at those canonical
sources. CHANGELOG records notable changes. brainstorm.md is raw notes and not authoritative;
job.md is not project documentation.

The README has to include setup prerequisites, how to run the API and web app locally and with
Compose, success and error API examples that match OpenAPI, the design decisions and
assumptions, out-of-scope items, the test/coverage/CI targets with links to the canonical
verification procedure, and the docs map. Call out the brief’s “Dockerfile to run frontend +
backend together” interpretation explicitly: two service Dockerfiles plus docker compose up. The
docs should be honest that the calculator is not implemented yet and that the scaffolding
currently has the pnpm workspace, Vite/React/Tailwind packages, and Storybook scripts while
TanStack Query, Zod, coverage scripts, Dockerfiles, Compose, CI, and the Go implementation land
later.

Finally, reconcile conflicts instead of preserving them. If two docs disagree, pick the
canonical home, move the definition there, and leave links. If a spec and architecture disagree,
stop and fix the docs before any coding. Keep the prose direct and readable, cross-link
aggressively, and don’t let the decision record turn into a duplicate of the architecture or the
component conventions.