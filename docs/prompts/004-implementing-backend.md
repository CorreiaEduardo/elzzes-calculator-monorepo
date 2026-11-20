> First I just asked to the AI to implement the specs 1 and 2 following the principles, architecture definition, decisions and conventions. But then I realised that i needed to change one core understanding, so I prompted:

The convention that transport/http imports only middleware and the standard library is too restrictive for our use case. The HTTP adapter needs to translate HTTP requests into calls to the calculator core. That naturally requires it to know about the calculator's public contract... HTTP handler should consume behavior not needing to depend on a concrete implementation.

Important: don't create a redundant Operation type inside the HTTP package. The calculator owns the meaning of an operation. Also, the transport maps domain errors to HTTP semantics.

So, instead of: transport/http imports only middleware and the standard library. I'd establish this: transport/http may depend on the calculator's public API, middleware, and the Go standard library. The calculator core must not depend on transport or HTTP concerns.

For this calculator, it is the simplest, most explicit, and most idiomatic Go boundary. The transport should depend on the core's public types because it is an adapter for that core, not an independently isolated application that needs to pretend it doesn't know what it's adapting.

Update specs or any relevant .md files IF NEEDED to correct this understanding.

> Then, doing my code review I realized some points that worth fixing, so I did this:
Lets change null and wrong-type handling, "operation": null gives UNKNOWN_OPERATION, but "operands": null gives INVALID_JSON, and the code special-cases null for one field only. Wrong types such as "operation": 1 or "operands": 5 return INVALID_JSON with the message “request body is not valid JSON”, even though the body is valid JSON. the message is misleading, and the existing codes (INVALID_OPERAND, INVALID_OPERAND_COUNT) may fit better.

Also, duplicate keys silently last-wins. {"operation":"add","operands":[1,2],"operands":[3,4]} computes with [3,4] and no test covers this. Document the behavior as intended/out of scope and add a test for it.

In parseCalculateRequest, any non-MaxBytesError failure from ReadAll becomes INVALID_JSON. That includes a client disconnecting or hitting ReadTimeout? a log line would help debugging, just in case.

parseCalculateRequest does too much. It’s about 110 lines covering content-type, size, decoding, unknown-field checks, operation parsing and operand parsing. Splitting it into readBody, parseOperation and parseOperands would read better. Making requestFailure implement error would let the helpers return plain error and the handler use errors.As.

Operand errors don’t say which operand. “operand must be a JSON number” would be more useful as operands[1] must be a JSON number.