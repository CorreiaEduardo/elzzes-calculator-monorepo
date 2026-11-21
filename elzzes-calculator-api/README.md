# elzzes-calculator-api

The Go REST service for Elzzes Calculator. It owns every arithmetic operation; clients never compute results. The service is stateless: each request is independent and nothing is persisted. It does not serve the SPA — Nginx does that in production.

**Status:** backend implemented — all seven operations through `POST /api/v1/calculate`, plus `GET /healthz`. Operation behavior is specified in [spec 001](../docs/specs/001-basic-arithmetic-operations/SPEC.md) and [spec 002](../docs/specs/002-advanced-arithmetic-operations/SPEC.md); the web UI remains [spec 003](../docs/specs/003-calculator-ui/SPEC.md). Stack: Go standard library (`net/http`, `encoding/json`, `log/slog`), OpenAPI 3.1, Docker.

## Responsibilities

- Expose `POST /api/v1/calculate` with the operation vocabulary, arity, and error codes from [OpenAPI](docs/openapi.yaml).
- Validate transport and payload strictly; return documented JSON errors.
- Enforce finite results and domain rules (division by zero, negative square root, overflow).
- Publish the wire contract in [`docs/openapi.yaml`](docs/openapi.yaml).
- Expose `GET /healthz`.

## Run

Requires Go 1.27 or newer.

```bash
cd elzzes-calculator-api
go run ./cmd/api              # PORT defaults to 8080
curl -s localhost:8080/healthz
curl -s localhost:8080/api/v1/calculate -H 'Content-Type: application/json' \
  -d '{"operation":"add","operands":[12.5,7.5]}'
```

All operations and error examples are in the [root README § API examples](../README.md#api-examples).

Development is API-only; Vite serves the SPA on `:5173` and proxies `/api`. In production Nginx serves the SPA and proxies `/api/` and `/healthz` to this service.

## Docker

```bash
docker build -t elzzes-calculator-api .
docker run --rm -p 8080:8080 elzzes-calculator-api
```

The multi-stage image builds the static binary and runs it as a non-root user in a minimal runtime. The root [`compose.yaml`](../compose.yaml) runs this image with the web image; Nginx proxies `/api` and `/healthz` to `api:8080` on one origin ([spec 004](../docs/specs/004-delivery-and-verification/SPEC.md)).

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `PORT` | `8080` | HTTP listen port (integer `1`–`65535`) |

No database, secrets, or external services.

## Test

```bash
gofmt -l .
go vet ./...
go build -o /tmp/elzzes-api ./cmd/api
go test ./...
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out
go test -race ./...
```

Coverage is reported for `./internal/...`; `cmd/api` is composition wiring.

The coverage target and edge-case rule, including CI, are canonical in [spec 004](../docs/specs/004-delivery-and-verification/SPEC.md).

## Conventions and layout

The layout, naming, error handling, testing, and logging rules are canonical in [`docs/CONVENTIONS.md`](docs/CONVENTIONS.md). High-level structure: `cmd/api/` (composition root), `internal/calculator/` (core), `internal/transport/http/` (net/http adapter and its consumer-owned `Logger` port), `internal/logger/` (slog adapter), and `internal/config/` (environment loading).

## Boundaries

- No UI logic, persistence, auth, static SPA serving, or third-party HTTP framework.
- The wire contract is [`docs/openapi.yaml`](docs/openapi.yaml); change it with the implementation and [Architecture § API contract](../docs/ARCHITECTURE.md#api-contract).
- Operation behavior and edge cases are canonical in specs 001–002.

See the [root README](../README.md) for setup, API examples, and design decisions.
