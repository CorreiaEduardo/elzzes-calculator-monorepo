# Backend conventions

**Status:** scaffold baseline, 2026-11-20. Canonical for Go implementation conventions in `elzzes-calculator-api`. The service runs `GET /healthz`; calculator behavior lands with specs 001–002. The wire contract is [`openapi.yaml`](openapi.yaml); system-wide structure and numeric rules are in [Architecture](../../docs/ARCHITECTURE.md); invariants are in [Principles](../../docs/PRINCIPLES.md).

## Layout

```text
elzzes-calculator-api/
├── cmd/api/main.go                  # composition root: wiring, server timeouts, graceful shutdown
├── internal/
│   ├── calculator/                  # pure core: calculator.go, errors.go (specs 001–002)
│   ├── transport/
│   │   └── http/                    # net/http adapter: routes, handler, errors, logger port
│   │       └── middleware/          # request logging, panic recovery
│   ├── config/                      # PORT loading and validation
│   └── logger/                      # slog adapter for the transport's Logger port
├── docs/
│   ├── CONVENTIONS.md
│   └── openapi.yaml                 # canonical wire contract
├── go.mod
├── Dockerfile
└── .dockerignore
```

- The backend is a lightweight Ports & Adapters (hexagonal) design with consumer-owned ports: a layer defines the small interface it needs in its own package, and an adapter implements it from the outside. There is no central `ports` package.
- `calculator` is the core: no HTTP or JSON types, no logging, config, OS, transport, or global state. Each operation is a pure function over finite `float64` values, returning a result or a typed domain error. It is a placeholder until specs 001–002.
- `transport/http` owns the net/http adapter: server-side routing, middleware, strict decoding, operation/arity validation, status/code mapping, the JSON error envelope, and `/healthz`. It declares the `Logger` port it needs in `logger.go`. Handlers are thin: decode → validate → call the calculator port → encode. Swapping to Fiber or another framework replaces this package and the wiring in `main`.
- `transport/http/middleware` depends only on the callbacks it uses (`func(msg string, args ...any)`), so it needs no logger interface of its own.
- `logger` is the only package importing `log/slog`; `config` is the only package reading the environment. Replacing either touches that package and one wiring line in `main`.
- `cmd/api` is the only composition point. It wires the adapters and builds the `http.Server`; it does not construct logging internals.
- Prefer one-word package and file names (`config`, `logger`, `http`, `env`) over compound placeholders (`envconfig`, `httpserver`).
- The service never serves the SPA; Nginx does that in production and Vite in development.

## HTTP and API rules

- Go standard library only: `net/http` with Go 1.27 method/pattern routing, `encoding/json`, `log/slog`. No third-party router or framework.
- All API routes live under `/api/v1/`. JSON in, JSON out; every response under `/api/` is JSON, including `404` and `405`.
- Routing returns the contract's JSON envelopes: an unmatched path is `404 NOT_FOUND`; a known path with the wrong method is `405 METHOD_NOT_ALLOWED` with an `Allow` header. GET patterns also serve HEAD. `GET /healthz` is registered by the HTTP adapter.
- Middleware order is request logging outermost, panic recovery next, router innermost, so a recovered panic is logged with status `500`.
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

- Domain errors are typed and explicit in `calculator`: division by zero (including `-0`), negative square root, and non-finite result.
- The HTTP adapter maps domain errors to OpenAPI codes and statuses centrally. Handlers do not write status codes ad hoc.
- The shared error envelope and `ErrorCode` constants live in `transport/http/errors.go` and mirror OpenAPI; [`openapi.yaml`](openapi.yaml) remains canonical.
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

- Prefer one-word names for packages, files, and abstractions: `config`, `logger`, `http`, `env`. Avoid compound names when a plain word is accurate.
- Source files: lowercase; a single word where possible (`handler.go`, `routes.go`), not compound placeholders.
- Exported identifiers PascalCase; unexported camelCase; acronyms follow Go conventions (`HTTPHandler`, `JSONRequest`).
- Constructors use `New{Type}` and take dependencies explicitly; no globals or `init` wiring.
- Prefer the standard library. A port is a small consumer-owned interface; introduce one only when a boundary is actually replaceable, and never as a shared package of speculative contracts.
- `gofmt` is mandatory. `go vet` and `golangci-lint` are the static gates; the race detector runs in CI.
- Comments are lean: explain constraints the code cannot show, not the obvious.

## Logging

- `transport/http/logger.go` declares the `Logger` port (`Info`, `Error`) because the HTTP transport consumes logging. `internal/logger` implements it with `log/slog` and is the only package importing slog.
- Middleware accepts the logging callbacks it uses (`logger.Info`, `logger.Error`) instead of declaring another interface.
- Log structured fields through the port. Log unexpected errors and request metadata; do not log business logic or success payloads.
- Keep operational logging out of `calculator`; `cmd/api` constructs the adapter and injects it as the transport's `Logger`.

## Testing

- `_test.go` files sit beside the implementation. Test names describe behavior (`TestCalculateDivisionByZero`).
- Use table-driven tests for operations and input validation; `httptest` for handlers.
- Test the core directly; test the HTTP adapter through `httptest` and pass fake loggers or recording callbacks to the ports and middleware. Adapters are concrete; no generated mocks.
- Cover every operation and every edge case in specs 001–002; the edge-case matrices are the source of test names.
- Use explicit float tolerance where binary64 is not exact; exact comparisons where it is.
- Prefer the standard library over assertion/mocking libraries.
- Run `go test -race ./...` in CI; `go test ./... -coverprofile=coverage.out` locally and in CI.
- No test starts a real network listener unless it is an `httptest` server.
- The ≥80% coverage floor applies to `internal/...` packages with behavior; `cmd/api` is composition wiring and has no test files.

## Configuration and delivery

- `PORT` defaults to `8080` and must be an integer in `1..65535`; no other configuration, secrets, or external services. `internal/config` is the only package that reads the environment.
- Multi-stage `Dockerfile`: build with the Go toolchain, run the static binary as a non-root user in a minimal image. The final image contains no source or toolchain; `.dockerignore` keeps the context minimal.
- The API container is stateless; no volumes. Compose runs it with the web container, and Nginx proxies `/api/` and `/healthz` to it ([spec 004](../../docs/specs/004-delivery-and-verification/SPEC.md)).

## Boundaries

- No UI logic, persistence, auth, rate limiting, static file serving, or third-party HTTP framework.
- Respect the layer boundary: `calculator` imports nothing from `internal/`; `transport/http` imports only its `middleware` and the standard library; `logger` and `config` import only the standard library; `cmd/api` is the only composition point and the only place adapters meet.
- Ports are declared by their consumer, next to its use. The provider package never imports the consumer; conformance is checked at the wiring point in `cmd/api`.
- No `/api/v1/calculate` route is registered until spec 001 implements it; do not add placeholder routes or contract codes that OpenAPI does not define.
- The wire contract changes only through a spec/OpenAPI/Architecture update agreed with the user ([AGENTS.md](../../AGENTS.md)).
