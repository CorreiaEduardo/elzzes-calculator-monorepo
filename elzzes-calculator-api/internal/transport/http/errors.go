package http

import nethttp "net/http"

// ErrorCode mirrors the ErrorCode enum in docs/openapi.yaml, which is canonical.
type ErrorCode string

const (
	CodeInvalidJSON         ErrorCode = "INVALID_JSON"
	CodeUnknownField        ErrorCode = "UNKNOWN_FIELD"
	CodeUnknownOperation    ErrorCode = "UNKNOWN_OPERATION"
	CodeInvalidOperandCount ErrorCode = "INVALID_OPERAND_COUNT"
	CodeInvalidOperand      ErrorCode = "INVALID_OPERAND"
	CodePayloadTooLarge     ErrorCode = "PAYLOAD_TOO_LARGE"
	CodeInvalidContentType  ErrorCode = "INVALID_CONTENT_TYPE"
	CodeDivisionByZero      ErrorCode = "DIVISION_BY_ZERO"
	CodeNegativeSquareRoot  ErrorCode = "NEGATIVE_SQUARE_ROOT"
	CodeNonFiniteResult     ErrorCode = "NON_FINITE_RESULT"
	CodeNotFound            ErrorCode = "NOT_FOUND"
	CodeMethodNotAllowed    ErrorCode = "METHOD_NOT_ALLOWED"
	CodeInternalError       ErrorCode = "INTERNAL_ERROR"
)

type errorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type errorEnvelope struct {
	Error errorDetail `json:"error"`
}

func writeError(w nethttp.ResponseWriter, status int, code ErrorCode, message string) {
	writeJSON(w, status, errorEnvelope{Error: errorDetail{Code: code, Message: message}})
}
