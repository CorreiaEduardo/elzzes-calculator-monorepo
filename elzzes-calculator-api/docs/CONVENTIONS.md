# Backend conventions

**Status:** planned baseline, 2026-11-20. Canonical for Go implementation conventions in `elzzes-calculator-api`. The wire contract is [`openapi.yaml`](openapi.yaml); system-wide structure and numeric rules are in [Architecture](../../docs/ARCHITECTURE.md); invariants are in [Principles](../../docs/PRINCIPLES.md).

## Layout

```text
elzzes-calculator-api/
├── cmd/api/main.go              # config, wiring, graceful shutdown
├── internal/
│   ├── config/                  # environment loading (PORT)
│   └── features/
│       ├── calculator/
│       │   ├── domain/          # Operation, pure functions, domain errors
│       │   └── http/            # handler, DTOs, routing, error mapping
│       └── commons/
│           └── errors/          # API error envelope, shared mappings
├── docs/
│   ├── CONVENTIONS.md
│   └── openapi.yaml             # canonical wire contract
├── go.mod
└── Dockerfile
```

- `domain` has no HTTP or JSON types and no global state. Each operation is a pure function over finite `float64` values, returning a result or a typed domain error.
- `http` owns routing, strict decoding, operation/arity validation, status/code mapping, and the `/healthz` handler. Handlers are thin: decode → validate → call domain → encode.
- `commons/errors` holds the shared error envelope and mapping helpers. It does not depend on the calculator domain.
- No `ports`, `services`, or `mocks` layer: the calculator has no external dependency to abstract. Introduce them only when a real dependency appears ([Principles § Simplicity budget](../../docs/PRINCIPLES.md#invariants)).
- The service never serves the SPA; Nginx does that in production and Vite in development.

## HTTP and API rules

- Go standard library only: `net/http` with Go 1.22+ method/pattern routing, `encoding/json`, `log/slog`. No third-party router or framework.
- All API routes live under `/api/v1/`. JSON in, JSON out; every response under `/api/` is JSON, including `404` and `405`.
- Enforce `Content-Type: application/json` before decoding; otherwise return `415 INVALID_CONTENT_TYPE`.
- Cap request bodies at the limit in [`openapi.yaml`](openapi.yaml); larger bodies return `413 PAYLOAD_TOO_LARGE`.
- Decode with unknown fields rejected. A missing body, invalid JSON, or wrong top-level type returns `400 INVALID_JSON`; an unknown field returns `400 UNKNOWN_FIELD`.
- Validate the operation against the registry and the operand count against its arity (`400 UNKNOWN_OPERATION`, `400 INVALID_OPERAND_COUNT`). Operands must be JSON numbers inside binary64 range (`400 INVALID_OPERAND`).
- Success returns `200` with `{operation, operands, result}`. Echo the request operands; do not round the result.
- Errors use one envelope: `{"error":{"code":"...","message":"..."}}`. `message` is developer-facing; the client owns user-facing copy. Codes, statuses, and schemas are in [`openapi.yaml`](openapi.yaml); do not maintain a second list.
- Add panic-recovery middleware (log + `500 INTERNAL_ERROR`), one structured request log line (method, path, status, duration), explicit server timeouts, and graceful shutdown on `SIGINT`/`SIGTERM`.
- `GET /healthz` returns `200 {"status":"ok"}`. There are no dependencies to check.
- No CORS headers: production is same-origin through Nginx; development uses the Vite proxy with relative `/api` paths.

## Error handling

- Domain errors are typed and explicit: division by zero (including `-0`), negative square root, and non-finite result. Define them in `domain`.
- The HTTP layer maps domain errors to OpenAPI codes and statuses centrally. Handlers do not write status codes ad hoc.
- Transport errors (bad JSON, unknown field/operation, operand count/value, body size, content type, route, method) are HTTP-layer errors; they do not enter the domain.
- Unexpected errors are logged with request metadata and returned as `500 INTERNAL_ERROR`; never leak internals.
- No `NaN` or `Infinity` is ever encoded. Negative-zero results are normalized to `0`.

## Numeric rules

The full numeric model is canonical in [Architecture § Numeric model](../../docs/ARCHITECTURE.md#numeric-model). In the service:

- Inputs are finite `float64` by construction.
- Check the divisor before dividing, including `-0`.
- Reject any result that is `NaN` or ±`Infinity` as `NON_FINITE_RESULT`.
- `0^0 = 1`; `0` to a negative exponent is non-finite; a negative base with a non-integer exponent is undefined.
- `percent` is unary: `x → x / 100`.
- Underflow to `0` is accepted.

## Naming and code style

- Source files: lowercase, underscores when needed (`calculation_handler.go`); package names short and lowercase.
- Exported identifiers PascalCase; unexported camelCase; acronyms follow Go conventions (`HTTPHandler`, `JSONRequest`).
- Constructors use `New{Type}` and take dependencies explicitly; no globals or `init` wiring.
- Prefer the standard library. Introduce an interface only for real decoupling or testability.
- `gofmt` is mandatory. `go vet` and `golangci-lint` are the static gates; the race detector runs in CI.
- Comments are lean: explain constraints the code cannot show, not the obvious.

## Logging

- Use `log/slog` with structured fields.
- Log unexpected errors and request metadata; do not log business logic or success payloads.
- Keep operational logging out of `domain`.

## Testing

- `_test.go` files sit beside the implementation. Test names describe behavior (`TestCalculateDivisionByZero`).
- Use table-driven tests for operations and input validation; `httptest` for handlers.
- Cover every operation and every edge case in specs 001–002; the edge-case matrices are the source of test names.
- Use explicit float tolerance where binary64 is not exact; exact comparisons where it is.
- Prefer the standard library over assertion/mocking libraries; there are no external interfaces to mock.
- Run `go test -race ./...` in CI; `go test ./... -coverprofile=coverage.out` locally and in CI.
- No test starts a real network listener unless it is an `httptest` server.

## Configuration and delivery

- `PORT` defaults to `8080`; no other configuration, secrets, or external services.
- Multi-stage `Dockerfile`: build with the Go toolchain, run the static binary as a non-root user in a minimal image. The final image contains no source or toolchain.
- The API container is stateless; no volumes. Compose runs it with the web container, and Nginx proxies `/api/` and `/healthz` to it ([spec 004](../../docs/specs/004-delivery-and-verification/SPEC.md)).

## Boundaries

- No UI logic, persistence, auth, rate limiting, static file serving, or third-party HTTP framework.
- The wire contract changes only through a spec/OpenAPI/Architecture update agreed with the user ([AGENTS.md](../../AGENTS.md)).
