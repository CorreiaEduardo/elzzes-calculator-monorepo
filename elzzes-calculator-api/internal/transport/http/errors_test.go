package http

import (
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorCodeValues(t *testing.T) {
	want := map[ErrorCode]string{
		CodeInvalidJSON:         "INVALID_JSON",
		CodeUnknownField:        "UNKNOWN_FIELD",
		CodeUnknownOperation:    "UNKNOWN_OPERATION",
		CodeInvalidOperandCount: "INVALID_OPERAND_COUNT",
		CodeInvalidOperand:      "INVALID_OPERAND",
		CodePayloadTooLarge:     "PAYLOAD_TOO_LARGE",
		CodeInvalidContentType:  "INVALID_CONTENT_TYPE",
		CodeDivisionByZero:      "DIVISION_BY_ZERO",
		CodeNegativeSquareRoot:  "NEGATIVE_SQUARE_ROOT",
		CodeNonFiniteResult:     "NON_FINITE_RESULT",
		CodeNotFound:            "NOT_FOUND",
		CodeMethodNotAllowed:    "METHOD_NOT_ALLOWED",
		CodeInternalError:       "INTERNAL_ERROR",
	}

	if len(want) != 13 {
		t.Fatalf("got %d codes, want 13", len(want))
	}
	for code, value := range want {
		if string(code) != value {
			t.Errorf("code %q = %q, want %q", code, string(code), value)
		}
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		status  int
		code    ErrorCode
		message string
	}{
		{status: nethttp.StatusBadRequest, code: CodeInvalidJSON, message: "request body is not valid JSON"},
		{status: nethttp.StatusNotFound, code: CodeNotFound, message: "not found"},
		{status: nethttp.StatusMethodNotAllowed, code: CodeMethodNotAllowed, message: "method GET is not allowed"},
		{status: nethttp.StatusInternalServerError, code: CodeInternalError, message: "internal error"},
	}

	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeError(rec, tt.status, tt.code, tt.message)

			if rec.Code != tt.status {
				t.Errorf("status = %d, want %d", rec.Code, tt.status)
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}

			var raw map[string]map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
				t.Fatalf("decode envelope: %v", err)
			}
			if len(raw) != 1 {
				t.Fatalf("top-level keys = %d, want 1", len(raw))
			}
			body, ok := raw["error"]
			if !ok {
				t.Fatalf("missing error key in %v", raw)
			}
			if len(body) != 2 || body["code"] != string(tt.code) || body["message"] != tt.message {
				t.Errorf("error = %v, want code %q message %q", body, tt.code, tt.message)
			}
		})
	}
}
