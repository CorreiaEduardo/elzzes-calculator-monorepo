# 001 — Basic arithmetic operations

**Status:** backend implemented (2026-11-20); UI acceptance remains under [spec 003](../003-calculator-ui/SPEC.md). Wire contract and examples: [`openapi.yaml`](../../../elzzes-calculator-api/docs/openapi.yaml). Operation behavior and edge cases below are canonical for these four operations. UI behavior: [Design](../../DESIGN.md). Implementation detail: [backend](../../../elzzes-calculator-api/docs/CONVENTIONS.md), [frontend](../../../elzzes-calculator-web/docs/CONVENTIONS.md).

## What is it?

Addition, subtraction, multiplication, and division — the four mandatory operations from `problem.md`, delivered end to end: Go domain functions, HTTP API, typed client, keypad, tests.

## Scope

**In**

- `add`, `subtract`, `multiply`, `divide` with exactly two operands.
- `POST /api/v1/calculate` success and error paths for these operations.
- Keypad operators `+ − × ÷`, chaining, equality, and error presentation.
- Validation shared by all operations: body, content type, strict fields, arity, numbers, finite results.

**Out**

- Advanced operations (spec 002).
- Expression parsing or precedence.
- Persistence, history, auth, rate limiting.

## Semantics

| Operation | `operation` | Arity | Result |
| --- | --- | --- | --- |
| Addition | `add` | 2 | `a + b` |
| Subtraction | `subtract` | 2 | `a − b` |
| Multiplication | `multiply` | 2 | `a × b` |
| Division | `divide` | 2 | `a ÷ b`; divisor `0`/`-0` is an error |

Operands are IEEE-754 binary64 values decoded from JSON numbers. The service returns unrounded results, echoes the request operands, and normalizes negative zero to `0`. Numeric rationale and alternatives: [Architecture § Numeric model](../../ARCHITECTURE.md#numeric-model).

## Edge cases (must be tested)

1. **Division by zero:** `10 ÷ 0`, `10 ÷ -0`, `0 ÷ 0` → `422 DIVISION_BY_ZERO`; no panic, no `Infinity`.
2. **Overflow:** e.g. `1.7976931348623157e308 × 2` → `422 NON_FINITE_RESULT`; equivalent cases for `add`, `subtract`, and `divide`.
3. **Underflow:** e.g. `5e-324 ÷ 10` returns `0`; accepted, not an error.
4. **Negative and fractional operands:** signs preserved; decimal operands accepted; `1 ÷ 3` returns the binary64 approximation (tests use tolerance or assert the formatted display, not string equality of the raw result).
5. **Negative zero:** a result of `-0` is normalized to `0` in the response.
6. **Malformed payloads:** missing body, invalid JSON, wrong top-level type, unknown field, missing/unknown `operation`, zero/one/three operands, string operand, `null` operands, out-of-range number (`1e999`), wrong content type, body over the contract limit → the documented 400/413/415 codes.
7. **Methods and routes:** `GET /api/v1/calculate` → `405 METHOD_NOT_ALLOWED`; unknown `/api/...` path → `404 NOT_FOUND`.
8. **Large and tiny operands within range:** round-trip through JSON without rejection.
9. **Duplicate object keys:** last value wins, matching `encoding/json` semantics; rejecting duplicate keys is out of scope.

Field-type errors are data errors, not syntax errors: a missing or non-string `operation` is `400 UNKNOWN_OPERATION`; a non-array `operands` value is `400 INVALID_OPERAND`; an invalid or out-of-range element is `400 INVALID_OPERAND` naming its index. `400 INVALID_JSON` is reserved for unparseable JSON, wrong top-level types, and trailing data.

Every case above maps to at least one Go test; UI-level cases for division by zero and generic errors live in spec 003.

## Acceptance criteria

- All four operations succeed for normal inputs through both `curl` and the UI, with the request visible (the UI never computes).
- Every edge case returns the documented status/code and a JSON envelope; the service never panics.
- The UI maps `DIVISION_BY_ZERO` to friendly copy and recovers per [Design § Errors and recovery](../../DESIGN.md#errors-and-recovery).
- `go vet`, `go test`, web lint/typecheck, and web tests pass; coverage target met; CI green.
- README examples match the OpenAPI contract and real behavior.

## Tests

- `domain`: table-driven tests per operation, including zero, negatives, decimals, overflow, underflow, negative zero, and division by zero. Float assertions use explicit tolerance or exact comparisons where binary64 is exact.
- `http`: `httptest` cases for success per operation and for each error row in the OpenAPI contract; strict JSON decoding; body limit; content type; method and route handling.
- Web: reducer tests for operator entry, replacement, chaining, and equality; component tests for the four operators; API-client tests with MSW fixtures derived from the OpenAPI examples; TanStack Query mutation asserted to make one request per action.

## Decision

Expose one `POST /api/v1/calculate` endpoint discriminated by `operation`, with fixed arity per operation and one error envelope.

## Reasons

- A fixed vocabulary centralizes validation, status mapping, and error codes; adding an operation is one registry entry plus one domain function.
- One endpoint keeps the client and tests small and avoids near-identical handlers.
- The versioned path (`/api/v1`) leaves room for a breaking change without silently altering behavior.
- Two operands per operation match the calculator interaction and keep validation unambiguous.

## Rejected alternatives

- **One endpoint per operation** (`/add`, `/subtract`, …): more routes, duplicated decode/validation, and a larger client; no benefit for a fixed vocabulary.
- **N-ary operands** (e.g., `operands: [1,2,3]` for add): ambiguous UI mapping and extra validation; the client is binary.
- **Client-side arithmetic** (JS operations or `eval`): violates the brief's client/service split, duplicates logic, and makes results unverifiable in one place.
- **Exact rational/decimal arithmetic:** rejected for the whole app in [Architecture § Numeric model](../../ARCHITECTURE.md#numeric-model); the UI formats binary64 results.
- **Operation-discovery endpoint:** the UI needs a bespoke keypad; discovery adds surface without a consumer.

## Open questions

None for scope. Exact decimal arithmetic and client-side history are documented future options, not commitments.
