# elzzes-calculator-api

The Go REST service for Elzzes Calculator. It owns every arithmetic operation; clients never compute results. The service is stateless: each request is independent and nothing is persisted. It does not serve the SPA — Nginx does that in production.

**Status:** documentation baseline and scaffolding — implementation is specified in [spec 001](../docs/specs/001-basic-arithmetic-operations/SPEC.md) and [spec 002](../docs/specs/002-advanced-arithmetic-operations/SPEC.md). Stack: Go standard library (`net/http`, `encoding/json`, `log/slog`), OpenAPI 3.1, Docker.

## Responsibilities

- Expose `POST /api/v1/calculate` with the operation vocabulary, arity, and error codes from [OpenAPI](docs/openapi.yaml).
- Validate transport and payload strictly; return documented JSON errors.
- Enforce finite results and domain rules (division by zero, negative square root, overflow).
- Publish the wire contract in [`docs/openapi.yaml`](docs/openapi.yaml).
- Expose `GET /healthz`.

## Run

```bash
cd elzzes-calculator-api
go run ./cmd/api              # PORT defaults to 8080
```

Development is API-only; Vite serves the SPA on `:5173` and proxies `/api`. In production Nginx serves the SPA and proxies `/api/` and `/healthz` to this service.

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `PORT` | `8080` | HTTP listen port |

No database, secrets, or external services.

## Test

```bash
go vet ./...
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out
go test -race ./...
```

The coverage target and edge-case rule, including CI, are canonical in [spec 004](../docs/specs/004-delivery-and-verification/SPEC.md).

## Conventions and layout

The layout, naming, error handling, testing, and logging rules are canonical in [`docs/CONVENTIONS.md`](docs/CONVENTIONS.md). High-level structure: `cmd/api/`, `internal/features/calculator/{domain,http}/`, `internal/features/commons/errors/`, `internal/config/`.

## Boundaries

- No UI logic, persistence, auth, static SPA serving, or third-party HTTP framework.
- The wire contract is [`docs/openapi.yaml`](docs/openapi.yaml); change it with the implementation and [Architecture § API contract](../docs/ARCHITECTURE.md#api-contract).
- Operation behavior and edge cases are canonical in specs 001–002.

See the [root README](../README.md) for setup, API examples, and design decisions.
