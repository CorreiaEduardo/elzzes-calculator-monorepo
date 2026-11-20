# Principles

Standing rules that outlive individual changes. Check a planned change against this file before building, and the result after.

## Order of concerns

1. **Correctness** — right answers, explicit domains, honest errors.
2. **Clarity** — code reads like the domain; names carry meaning.
3. **Maintainability** — small modules, stable contracts, tests as documentation.
4. **Features** — only what the brief asks for.

When they conflict, the higher concern wins. Extra features never buy down correctness debt.

## How to approach a new problem

1. Read `problem.md` and the relevant spec. Confirm the requirement exists; note what it does *not* ask for.
2. Re-read the invariants below and the architecture at the boundary you will touch.
3. Write the edge-case list first; make it concrete in tests.
4. Implement the smallest correct slice, keeping one change reviewable.
5. Verify: tests, `go vet`/typecheck/lint, coverage, CI, manual run (local and Compose when packaging changes).
6. Update docs in the same change, editing the canonical source and linking from other docs: specs/DESIGN for behavior, OpenAPI/ARCHITECTURE for contract and structure, README for setup/API examples, component conventions for implementation details, CHANGELOG when notable.
7. Commit the change on its own.

## Invariants

1. **One source of arithmetic truth.** Every result the user sees is computed by the Go service. The client may parse input and format output; it may not evaluate operations.
2. **State ownership.** The service owns computed values and their validity. The web app owns transient interaction state (entry, pending operation, phase); TanStack Query owns only the request lifecycle. `@elzzes/ui` owns nothing — it is controlled presentation. No layer reaches around another.
3. **Local-first.** The whole app runs on a developer machine with no cloud account, external service, or emulator. Docker Compose is the only deployment artifact; CI runs the same checks as the documented local procedure. Nothing is required beyond the repository and the toolchain.
4. **Finite results or explicit errors.** Division by zero, negative square roots, overflow, and undefined results return documented error codes. `NaN`/`Infinity` never leak to the UI.
5. **Stable contract.** The wire contract is OpenAPI 3.1 (`elzzes-calculator-api/docs/openapi.yaml`). Operation vocabulary, arity, request/response shape, and error codes change only through a spec/OpenAPI/`ARCHITECTURE.md` update agreed with the user.
6. **Validate at the boundary, trust inside.** HTTP input is validated once at the edge; domain functions receive validated, finite values. The client validates responses with Zod and improves UX with input rules, but never replaces server validation.
7. **Errors are data.** The API carries stable codes; human copy belongs to the client.
8. **Statelessness.** No database, session, persisted cache, or history. Restarting loses nothing because nothing is stored. TanStack Query manages a request, not persisted server state.
9. **Accessibility and responsiveness are requirements.** Keyboard operation, screen-reader labels, AA contrast, touch targets, and a usable phone layout are part of "done".
10. **Simplicity budget.** Every dependency, abstraction, or process needs a requirement behind it. The agreed stack is Tailwind CSS v4, TanStack Query, Zod, Storybook, OpenAPI 3.1, Docker Compose, and GitHub Actions; add nothing beyond it without a requirement and user agreement. No database, auth, or state library.
11. **Docs travel with the change.** A behavior change that leaves docs stale is unfinished. Each definition has one canonical home; other docs link to it rather than restating it.
12. **Lean comments.** Default to none; naming and structure carry the *what*. Comment only constraints the code cannot show (timing races, ordering requirements, standing invariants).
13. **Consumer-owned ports.** A boundary interface is defined by the layer that needs the capability, in that layer's package (for example `transport/http/logger.go`), and implemented by the provider (`internal/logger`). No shared `ports` package; no interface before a real replacement boundary exists.
14. **Plain names.** Prefer one-word names for packages, files, and abstractions (`config`, `logger`, `http`, `env`). Avoid compound placeholders when a plain word is accurate.

## Decision rules

- If two options are both correct, choose the one with fewer moving parts and the clearer contract.
- If a choice crosses layer boundaries or is expensive to reverse, write it down first (spec, `ARCHITECTURE.md`, or `openapi.yaml`) and get agreement.
- If the brief is ambiguous, choose the narrowest reading, record the assumption in the spec and README, and do not add speculative generality.
- Prefer deleting a feature over carrying an untested path.
- A duplicated definition is a defect: move it to its canonical home and leave a link.

## What "done" means

Tests encode the edge cases, not just the happy path. Typecheck/vet/lint are clean, coverage is reported and meets the target, CI is green, docs match behavior, and the change can be reverted as one unit. A reviewer should be able to derive the intended behavior from the spec and the tests without reading the implementation twice.
