# 002 — Advanced arithmetic operations

**Status:** scoped (documentation baseline, 2026-11-20). Depends on the contract and validation from [spec 001](../001-basic-arithmetic-operations/SPEC.md). Wire contract and examples: [`openapi.yaml`](../../../elzzes-calculator-api/docs/openapi.yaml). Operation behavior and edge cases below are canonical for these three operations. UI behavior: [Design](../../DESIGN.md).

## What is it?

The three optional operations from `problem.md` — exponentiation, square root, and percentage — delivered through the same endpoint and keypad, with explicit domains.

## Scope

**In**

- `power` (binary), `sqrt` (unary), `percent` (unary).
- Domain validation, finite-result enforcement, error codes.
- Keypad `xʸ`, `√`, `%`, including unary application to entries and results.
- Tests for every documented case.

**Out**

- Complex numbers, reciprocal/other unary functions, roots other than square root.
- Expression parsing, precedence, and scientific-notation entry beyond digit entry.
- Contextual percentage semantics.

## Semantics

| Operation | `operation` | Arity | Definition | Domain |
| --- | --- | --- | --- | --- |
| Exponentiation | `power` | 2 | `aᵇ` via `math.Pow` | result must be finite |
| Square root | `sqrt` | 1 | principal square root via `math.Sqrt` | operand ≥ 0 |
| Percentage | `percent` | 1 | `x / 100` | always defined for finite `x` |

All other contract rules (JSON shape, strict fields, body limit, statuses, response envelope) are inherited from spec 001 and [`openapi.yaml`](../../../elzzes-calculator-api/docs/openapi.yaml).

## Edge cases (must be tested)

1. **`0^0 = 1`** — documented convention, matching Go's `math.Pow`; not an error.
2. **Zero to a negative exponent** (`0^-1`) → non-finite → `422 NON_FINITE_RESULT`.
3. **Negative base with a fractional exponent** (`(-8)^0.5`) → `NaN` → `422 NON_FINITE_RESULT`.
4. **Negative base with an integer exponent** (`(-2)^3 = -8`) → valid.
5. **Fractional exponent with a positive base** (`9^0.5 = 3`) → valid.
6. **Overflow** (`10^400`) → `422 NON_FINITE_RESULT`; large-but-finite results pass.
7. **Negative exponent** (`2^-3 = 0.125`) → valid.
8. **`sqrt` of a negative operand** (`sqrt(-1)`) → `422 NEGATIVE_SQUARE_ROOT`.
9. **`sqrt(0) = 0`**, `sqrt(-0)` is accepted and normalized to `0`, `sqrt(2)` returns the binary64 approximation.
10. **Percentage:** `12.5 → 0.125`, `-12.5 → -0.125`, `0 → 0`; huge-but-finite values remain finite or produce `NON_FINITE_RESULT` if they overflow.
11. **Wrong arity:** `sqrt` with two operands, `power` with one, etc. → `400 INVALID_OPERAND_COUNT`.
12. **Unknown operation** → `400 UNKNOWN_OPERATION`.
13. **Negative-zero result** normalized to `0` (e.g., `sqrt(-0)`, `percent(-0)`).

## Acceptance criteria

- All three operations work through `curl` and the UI and use the same error envelope as 001.
- Every domain case above is enforced by the service and covered by tests; no `NaN`/`Infinity` is ever returned.
- The UI exposes `√`, `xʸ`, `%` with the interaction rules in [Design § Interaction model](../../DESIGN.md#interaction-model); unary keys are no-ops when a second operand is empty.
- Advanced operations add no new endpoint, no client arithmetic, and no new dependency.
- Coverage target remains met; CI green.

## Tests

- `domain`: table-driven cases for the full matrix above, including `math.Pow` domain edges.
- `http`: success and 422 cases for each operation; arity errors.
- Web: unary-on-entry, unary-on-result, no-op behavior, power as pending operator; display/format checks for fractional and negative results.
- API client: MSW fixtures for `NON_FINITE_RESULT` and `NEGATIVE_SQUARE_ROOT`; copy mapping asserted.

## Decision

Include all three optional operations, define percentage as unary, and treat any non-finite result as a domain error.

## Reasons

- The brief lists them as optional, but they cost one domain function, one registry entry, and three keys each; they also prove that the contract handles unary as well as binary arity.
- Unary percentage is unambiguous: the key transforms the displayed value. Contextual percentage requires expression context the UI deliberately does not have.
- A single `NON_FINITE_RESULT` code keeps the error surface small while the UI copy explains "too large or undefined".
- Using `math.Pow`/`math.Sqrt` is the platform's well-tested implementation; hand-rolled powers would lose precision and add edge cases.

## Rejected alternatives

- **Skip the optional operations:** the brief explicitly invites them and the architecture is already generic over arity.
- **Contextual percentage** (`200 + 10% → 220`): ambiguous, inconsistent across calculators, and needs pending-operation introspection; rejected in favor of `x / 100`.
- **Binary "percent of"** (`percent(a, b)` = a% of b): a second meaning with no keypad affordance; deferred rather than guessed.
- **Complex results** for negative square roots/powers: out of scope; returning a real result the UI cannot display would be worse than a clear domain error.
- **Repeated multiplication for power:** wrong for fractional/negative exponents and slower to converge; `math.Pow` is the correct primitive.
- **Additional operations** (log, trig, reciprocal): not requested; scope stays closed.

## Open questions

None. Complex numbers and contextual percentage would each need a new spec and keypad design.
