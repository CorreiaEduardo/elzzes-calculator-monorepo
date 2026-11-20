package http

import (
	"encoding/json"
	"errors"
	"math"
	nethttp "net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"elzzes-calculator-api/internal/calculator"
)

const calculatePath = "/api/v1/calculate"

type calculateSuccess struct {
	Operation string    `json:"operation"`
	Operands  []float64 `json:"operands"`
	Result    float64   `json:"result"`
}

func serveBody(h nethttp.Handler, method, path, contentType, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func serveCalculate(h nethttp.Handler, body string) *httptest.ResponseRecorder {
	return serveBody(h, nethttp.MethodPost, calculatePath, "application/json", body)
}

func decodeSuccess(t *testing.T, rec *httptest.ResponseRecorder) calculateSuccess {
	t.Helper()
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode response %q: %v", rec.Body.String(), err)
	}
	if len(raw) != 3 {
		t.Fatalf("response keys = %v, want operation/operands/result", raw)
	}
	var got calculateSuccess
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response %q: %v", rec.Body.String(), err)
	}
	return got
}

type stubCalculator struct {
	result    float64
	err       error
	calls     int
	operation calculator.Operation
	operands  []float64
}

func (s *stubCalculator) Calculate(operation calculator.Operation, operands []float64) (float64, error) {
	s.calls++
	s.operation = operation
	s.operands = operands
	return s.result, s.err
}

type recordingLogger struct {
	errors []string
}

func (l *recordingLogger) Info(string, ...any) {}

func (l *recordingLogger) Error(msg string, _ ...any) {
	l.errors = append(l.errors, msg)
}

func TestCalculateSuccess(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantOp       string
		wantOperands []float64
		want         float64
		tolerance    float64
	}{
		{name: "add", body: `{"operation":"add","operands":[12.5,7.5]}`, wantOp: "add", wantOperands: []float64{12.5, 7.5}, want: 20},
		{name: "add negative", body: `{"operation":"add","operands":[-12.5,7.5]}`, wantOp: "add", wantOperands: []float64{-12.5, 7.5}, want: -5},
		{name: "subtract", body: `{"operation":"subtract","operands":[12.5,7.5]}`, wantOp: "subtract", wantOperands: []float64{12.5, 7.5}, want: 5},
		{name: "multiply", body: `{"operation":"multiply","operands":[12.5,7.5]}`, wantOp: "multiply", wantOperands: []float64{12.5, 7.5}, want: 93.75},
		{name: "divide", body: `{"operation":"divide","operands":[10,4]}`, wantOp: "divide", wantOperands: []float64{10, 4}, want: 2.5},
		{name: "divide repeating", body: `{"operation":"divide","operands":[1,3]}`, wantOp: "divide", wantOperands: []float64{1, 3}, want: 1.0 / 3.0, tolerance: 1e-12},
		{name: "power", body: `{"operation":"power","operands":[2,10]}`, wantOp: "power", wantOperands: []float64{2, 10}, want: 1024},
		{name: "power zero to zero", body: `{"operation":"power","operands":[0,0]}`, wantOp: "power", wantOperands: []float64{0, 0}, want: 1},
		{name: "power negative base", body: `{"operation":"power","operands":[-2,3]}`, wantOp: "power", wantOperands: []float64{-2, 3}, want: -8},
		{name: "sqrt", body: `{"operation":"sqrt","operands":[144]}`, wantOp: "sqrt", wantOperands: []float64{144}, want: 12},
		{name: "sqrt zero", body: `{"operation":"sqrt","operands":[0]}`, wantOp: "sqrt", wantOperands: []float64{0}, want: 0},
		{name: "percent", body: `{"operation":"percent","operands":[12.5]}`, wantOp: "percent", wantOperands: []float64{12.5}, want: 0.125},
		{name: "percent negative", body: `{"operation":"percent","operands":[-12.5]}`, wantOp: "percent", wantOperands: []float64{-12.5}, want: -0.125},
		{name: "large operand", body: `{"operation":"add","operands":[1.7976931348623157e308,0]}`, wantOp: "add", wantOperands: []float64{math.MaxFloat64, 0}, want: math.MaxFloat64},
		{name: "tiny operand", body: `{"operation":"add","operands":[5e-324,0]}`, wantOp: "add", wantOperands: []float64{math.SmallestNonzeroFloat64, 0}, want: math.SmallestNonzeroFloat64},
		{name: "underflow", body: `{"operation":"divide","operands":[5e-324,10]}`, wantOp: "divide", wantOperands: []float64{math.SmallestNonzeroFloat64, 10}, want: 0},
	}

	handler := NewHandler(fakeLogger{}, calculator.New())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serveCalculate(handler, tt.body)
			if rec.Code != nethttp.StatusOK {
				t.Fatalf("status = %d, want 200; body %q", rec.Code, rec.Body.String())
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}
			if strings.Contains(rec.Body.String(), "NaN") || strings.Contains(rec.Body.String(), "Inf") {
				t.Errorf("body contains a non-finite literal: %q", rec.Body.String())
			}

			got := decodeSuccess(t, rec)
			if got.Operation != tt.wantOp {
				t.Errorf("operation = %q, want %q", got.Operation, tt.wantOp)
			}
			if !slices.Equal(got.Operands, tt.wantOperands) {
				t.Errorf("operands = %v, want %v", got.Operands, tt.wantOperands)
			}
			if !closeEnough(got.Result, tt.want, tt.tolerance) {
				t.Errorf("result = %v, want %v (tolerance %v)", got.Result, tt.want, tt.tolerance)
			}
		})
	}
}

func TestCalculateEchoesNegativeZeroOperand(t *testing.T) {
	rec := serveCalculate(NewHandler(fakeLogger{}, calculator.New()), `{"operation":"percent","operands":[-0]}`)
	if rec.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want 200; body %q", rec.Code, rec.Body.String())
	}

	got := decodeSuccess(t, rec)
	if len(got.Operands) != 1 || !math.Signbit(got.Operands[0]) {
		t.Errorf("operands = %v, want one -0 value", got.Operands)
	}
	if got.Result != 0 || math.Signbit(got.Result) {
		t.Errorf("result = %v (signbit %v), want +0", got.Result, math.Signbit(got.Result))
	}
}

func TestCalculateDomainErrors(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		code    ErrorCode
		message string
	}{
		{name: "divide by zero", body: `{"operation":"divide","operands":[10,0]}`, code: CodeDivisionByZero, message: "cannot divide by zero"},
		{name: "divide by negative zero", body: `{"operation":"divide","operands":[10,-0]}`, code: CodeDivisionByZero, message: "cannot divide by zero"},
		{name: "zero divided by zero", body: `{"operation":"divide","operands":[0,0]}`, code: CodeDivisionByZero, message: "cannot divide by zero"},
		{name: "add overflow", body: `{"operation":"add","operands":[1.7976931348623157e308,1.7976931348623157e308]}`, code: CodeNonFiniteResult, message: "result is not finite"},
		{name: "subtract overflow", body: `{"operation":"subtract","operands":[-1.7976931348623157e308,1.7976931348623157e308]}`, code: CodeNonFiniteResult, message: "result is not finite"},
		{name: "multiply overflow", body: `{"operation":"multiply","operands":[1.7976931348623157e308,2]}`, code: CodeNonFiniteResult, message: "result is not finite"},
		{name: "divide overflow", body: `{"operation":"divide","operands":[1.7976931348623157e308,5e-324]}`, code: CodeNonFiniteResult, message: "result is not finite"},
		{name: "power overflow", body: `{"operation":"power","operands":[10,400]}`, code: CodeNonFiniteResult, message: "result is not finite"},
		{name: "zero to negative exponent", body: `{"operation":"power","operands":[0,-1]}`, code: CodeNonFiniteResult, message: "result is not finite"},
		{name: "negative base fractional exponent", body: `{"operation":"power","operands":[-8,0.5]}`, code: CodeNonFiniteResult, message: "result is not finite"},
		{name: "sqrt negative", body: `{"operation":"sqrt","operands":[-1]}`, code: CodeNegativeSquareRoot, message: "square root requires a non-negative operand"},
	}

	handler := NewHandler(fakeLogger{}, calculator.New())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serveCalculate(handler, tt.body)
			if rec.Code != nethttp.StatusUnprocessableEntity {
				t.Fatalf("status = %d, want 422; body %q", rec.Code, rec.Body.String())
			}
			envelope := decodeEnvelope(t, rec)
			if envelope.Error.Code != tt.code || envelope.Error.Message != tt.message {
				t.Errorf("envelope = %+v, want %s/%q", envelope.Error, tt.code, tt.message)
			}
		})
	}
}

func TestCalculateMalformedRequests(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		code    ErrorCode
		message string
	}{
		{name: "empty body", body: "", code: CodeInvalidJSON, message: "request body is not valid JSON"},
		{name: "whitespace body", body: "   ", code: CodeInvalidJSON, message: "request body is not valid JSON"},
		{name: "truncated json", body: `{`, code: CodeInvalidJSON, message: "request body is not valid JSON"},
		{name: "not json", body: "not-json", code: CodeInvalidJSON, message: "request body is not valid JSON"},
		{name: "array top level", body: `[]`, code: CodeInvalidJSON, message: "request body is not valid JSON"},
		{name: "string top level", body: `"add"`, code: CodeInvalidJSON, message: "request body is not valid JSON"},
		{name: "number top level", body: `5`, code: CodeInvalidJSON, message: "request body is not valid JSON"},
		{name: "bool top level", body: `true`, code: CodeInvalidJSON, message: "request body is not valid JSON"},
		{name: "null top level", body: `null`, code: CodeInvalidJSON, message: "request body is not valid JSON"},
		{name: "non-string operation", body: `{"operation":1,"operands":[1,2]}`, code: CodeUnknownOperation, message: "operation must be a string"},
		{name: "trailing value", body: `{"operation":"add","operands":[1,2]} {"x":1}`, code: CodeInvalidJSON, message: "request body is not valid JSON"},
		{name: "trailing comma", body: `{"operation":"add","operands":[1,2],}`, code: CodeInvalidJSON, message: "request body is not valid JSON"},
		{name: "unknown field", body: `{"operation":"add","operands":[1,2],"precision":2}`, code: CodeUnknownField, message: `unknown field "precision"`},
		{name: "multiple unknown fields sorted", body: `{"operation":"add","operands":[1,2],"precision":2,"alpha":1}`, code: CodeUnknownField, message: `unknown field "alpha"`},
		{name: "empty object", body: `{}`, code: CodeUnknownOperation, message: "operation is required"},
		{name: "missing operation", body: `{"operands":[1,2]}`, code: CodeUnknownOperation, message: "operation is required"},
		{name: "null operation", body: `{"operation":null,"operands":[1,2]}`, code: CodeUnknownOperation, message: "operation must be a string"},
		{name: "empty operation", body: `{"operation":"","operands":[1,2]}`, code: CodeUnknownOperation, message: `unknown operation ""`},
		{name: "unknown operation", body: `{"operation":"modulo","operands":[1,2]}`, code: CodeUnknownOperation, message: `unknown operation "modulo"`},
		{name: "add no operands", body: `{"operation":"add","operands":[]}`, code: CodeInvalidOperandCount, message: "add expects 2 operands"},
		{name: "add one operand", body: `{"operation":"add","operands":[1]}`, code: CodeInvalidOperandCount, message: "add expects 2 operands"},
		{name: "add three operands", body: `{"operation":"add","operands":[1,2,3]}`, code: CodeInvalidOperandCount, message: "add expects 2 operands"},
		{name: "power one operand", body: `{"operation":"power","operands":[1]}`, code: CodeInvalidOperandCount, message: "power expects 2 operands"},
		{name: "sqrt no operands", body: `{"operation":"sqrt","operands":[]}`, code: CodeInvalidOperandCount, message: "sqrt expects 1 operand"},
		{name: "sqrt two operands", body: `{"operation":"sqrt","operands":[1,2]}`, code: CodeInvalidOperandCount, message: "sqrt expects 1 operand"},
		{name: "percent two operands", body: `{"operation":"percent","operands":[1,2]}`, code: CodeInvalidOperandCount, message: "percent expects 1 operand"},
		{name: "missing operands", body: `{"operation":"add"}`, code: CodeInvalidOperandCount, message: "add expects 2 operands"},
		{name: "null operands", body: `{"operation":"add","operands":null}`, code: CodeInvalidOperand, message: "operands must be an array of numbers"},
		{name: "number operands", body: `{"operation":"add","operands":5}`, code: CodeInvalidOperand, message: "operands must be an array of numbers"},
		{name: "object operands", body: `{"operation":"add","operands":{}}`, code: CodeInvalidOperand, message: "operands must be an array of numbers"},
		{name: "string operands", body: `{"operation":"add","operands":"[1,2]"}`, code: CodeInvalidOperand, message: "operands must be an array of numbers"},
		{name: "string operand", body: `{"operation":"add","operands":["1",2]}`, code: CodeInvalidOperand, message: "operands[0] must be a JSON number"},
		{name: "second invalid operand", body: `{"operation":"add","operands":[1,"2"]}`, code: CodeInvalidOperand, message: "operands[1] must be a JSON number"},
		{name: "null operand", body: `{"operation":"add","operands":[null,2]}`, code: CodeInvalidOperand, message: "operands[0] must be a JSON number"},
		{name: "bool operand", body: `{"operation":"add","operands":[true,2]}`, code: CodeInvalidOperand, message: "operands[0] must be a JSON number"},
		{name: "object operand", body: `{"operation":"add","operands":[{"a":1},2]}`, code: CodeInvalidOperand, message: "operands[0] must be a JSON number"},
		{name: "array operand", body: `{"operation":"add","operands":[[1],2]}`, code: CodeInvalidOperand, message: "operands[0] must be a JSON number"},
		{name: "out of range operand", body: `{"operation":"add","operands":[1e999,2]}`, code: CodeInvalidOperand, message: "operands[0] must be a JSON number"},
	}

	handler := NewHandler(fakeLogger{}, calculator.New())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serveCalculate(handler, tt.body)
			if rec.Code != nethttp.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body %q", rec.Code, rec.Body.String())
			}
			envelope := decodeEnvelope(t, rec)
			if envelope.Error.Code != tt.code || envelope.Error.Message != tt.message {
				t.Errorf("envelope = %+v, want %s/%q", envelope.Error, tt.code, tt.message)
			}
		})
	}
}

func TestCalculateDuplicateKeysLastValueWins(t *testing.T) {
	handler := NewHandler(fakeLogger{}, calculator.New())

	t.Run("operands", func(t *testing.T) {
		rec := serveCalculate(handler, `{"operation":"add","operands":[1,2],"operands":[3,4]}`)
		if rec.Code != nethttp.StatusOK {
			t.Fatalf("status = %d, want 200; body %q", rec.Code, rec.Body.String())
		}
		got := decodeSuccess(t, rec)
		if got.Result != 7 || !slices.Equal(got.Operands, []float64{3, 4}) {
			t.Errorf("response = %+v, want last operands [3 4] and result 7", got)
		}
	})

	t.Run("operation", func(t *testing.T) {
		rec := serveCalculate(handler, `{"operation":"add","operation":"multiply","operands":[3,4]}`)
		if rec.Code != nethttp.StatusOK {
			t.Fatalf("status = %d, want 200; body %q", rec.Code, rec.Body.String())
		}
		got := decodeSuccess(t, rec)
		if got.Operation != "multiply" || got.Result != 12 {
			t.Errorf("response = %+v, want last operation multiply and result 12", got)
		}
	})
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func (failingReader) Close() error {
	return nil
}

func TestCalculateReadError(t *testing.T) {
	logger := &recordingLogger{}
	handler := NewHandler(logger, calculator.New())

	req := httptest.NewRequest(nethttp.MethodPost, calculatePath, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = failingReader{}
	req.ContentLength = -1
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body %q", rec.Code, rec.Body.String())
	}
	envelope := decodeEnvelope(t, rec)
	if envelope.Error.Code != CodeInvalidJSON || envelope.Error.Message != "request body is not valid JSON" {
		t.Errorf("envelope = %+v, want INVALID_JSON", envelope.Error)
	}
	if len(logger.errors) != 1 || logger.errors[0] != "calculate request read failed" {
		t.Errorf("error logs = %v, want one calculate request read failed entry", logger.errors)
	}
}

func TestCalculateContentType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		wantStatus  int
	}{
		{name: "missing", contentType: "", wantStatus: nethttp.StatusUnsupportedMediaType},
		{name: "text plain", contentType: "text/plain", wantStatus: nethttp.StatusUnsupportedMediaType},
		{name: "form encoded", contentType: "application/x-www-form-urlencoded", wantStatus: nethttp.StatusUnsupportedMediaType},
		{name: "malformed", contentType: "application/json; charset", wantStatus: nethttp.StatusUnsupportedMediaType},
		{name: "case insensitive", contentType: "Application/JSON", wantStatus: nethttp.StatusOK},
		{name: "charset parameter", contentType: "application/json; charset=utf-8", wantStatus: nethttp.StatusOK},
	}

	handler := NewHandler(fakeLogger{}, calculator.New())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serveBody(handler, nethttp.MethodPost, calculatePath, tt.contentType, `{"operation":"add","operands":[1,2]}`)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body %q", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantStatus == nethttp.StatusUnsupportedMediaType {
				envelope := decodeEnvelope(t, rec)
				if envelope.Error.Code != CodeInvalidContentType || envelope.Error.Message != "Content-Type must be application/json" {
					t.Errorf("envelope = %+v, want INVALID_CONTENT_TYPE", envelope.Error)
				}
			}
		})
	}
}

func TestCalculateBodyLimit(t *testing.T) {
	const base = `{"operation":"add","operands":[1,2]}`
	handler := NewHandler(fakeLogger{}, calculator.New())

	t.Run("exactly four kibibytes", func(t *testing.T) {
		body := base + strings.Repeat(" ", 4096-len(base))
		rec := serveCalculate(handler, body)
		if rec.Code != nethttp.StatusOK {
			t.Fatalf("status = %d, want 200; body %q", rec.Code, rec.Body.String())
		}
	})

	t.Run("declared length over limit", func(t *testing.T) {
		body := base + strings.Repeat(" ", 4097-len(base))
		rec := serveCalculate(handler, body)
		if rec.Code != nethttp.StatusRequestEntityTooLarge {
			t.Fatalf("status = %d, want 413; body %q", rec.Code, rec.Body.String())
		}
		envelope := decodeEnvelope(t, rec)
		if envelope.Error.Code != CodePayloadTooLarge || envelope.Error.Message != "request body exceeds 4 KiB" {
			t.Errorf("envelope = %+v, want PAYLOAD_TOO_LARGE", envelope.Error)
		}
	})

	t.Run("unknown length over limit", func(t *testing.T) {
		body := base + strings.Repeat(" ", 5000)
		req := httptest.NewRequest(nethttp.MethodPost, calculatePath, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.ContentLength = -1
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != nethttp.StatusRequestEntityTooLarge {
			t.Fatalf("status = %d, want 413; body %q", rec.Code, rec.Body.String())
		}
	})
}

func TestCalculateMethodNotAllowed(t *testing.T) {
	handler := NewHandler(fakeLogger{}, calculator.New())
	for _, method := range []string{nethttp.MethodGet, nethttp.MethodPut, nethttp.MethodOptions} {
		t.Run(method, func(t *testing.T) {
			rec := serveBody(handler, method, calculatePath, "", "")
			if rec.Code != nethttp.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want 405", rec.Code)
			}
			if got := rec.Header().Get("Allow"); got != "POST" {
				t.Errorf("Allow = %q, want POST", got)
			}
			envelope := decodeEnvelope(t, rec)
			wantMessage := "method " + method + " is not allowed"
			if envelope.Error.Code != CodeMethodNotAllowed || envelope.Error.Message != wantMessage {
				t.Errorf("envelope = %+v, want METHOD_NOT_ALLOWED/%q", envelope.Error, wantMessage)
			}
		})
	}
}

func TestCalculateUnknownRoute(t *testing.T) {
	rec := serveBody(NewHandler(fakeLogger{}, calculator.New()), nethttp.MethodPost, "/api/v1/unknown", "application/json", "{}")
	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	envelope := decodeEnvelope(t, rec)
	if envelope.Error.Code != CodeNotFound || envelope.Error.Message != "not found" {
		t.Errorf("envelope = %+v, want NOT_FOUND/not found", envelope.Error)
	}
}

func TestCalculateUnexpectedError(t *testing.T) {
	logger := &recordingLogger{}
	calc := &stubCalculator{err: errors.New("boom")}
	rec := serveCalculate(NewHandler(logger, calc), `{"operation":"add","operands":[1,2]}`)

	if rec.Code != nethttp.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "boom") {
		t.Errorf("body leaks internal error: %q", rec.Body.String())
	}
	envelope := decodeEnvelope(t, rec)
	if envelope.Error.Code != CodeInternalError || envelope.Error.Message != "internal error" {
		t.Errorf("envelope = %+v, want INTERNAL_ERROR/internal error", envelope.Error)
	}
	if len(logger.errors) != 1 || logger.errors[0] != "calculate failed" {
		t.Errorf("error logs = %v, want one calculate failed entry", logger.errors)
	}
}

func TestCalculatePassesValidatedInputToPort(t *testing.T) {
	calc := &stubCalculator{result: 42}
	rec := serveCalculate(NewHandler(fakeLogger{}, calc), `{"operation":"add","operands":[12.5,7.5]}`)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want 200; body %q", rec.Code, rec.Body.String())
	}
	if calc.calls != 1 {
		t.Errorf("calculator calls = %d, want 1", calc.calls)
	}
	if calc.operation != calculator.Add {
		t.Errorf("operation = %q, want add", calc.operation)
	}
	if !slices.Equal(calc.operands, []float64{12.5, 7.5}) {
		t.Errorf("operands = %v, want [12.5 7.5]", calc.operands)
	}

	invalid := &stubCalculator{result: 1}
	serveCalculate(NewHandler(fakeLogger{}, invalid), `{"operation":"modulo","operands":[1,2]}`)
	if invalid.calls != 0 {
		t.Errorf("calculator calls for invalid request = %d, want 0", invalid.calls)
	}
}

func closeEnough(got, want, tolerance float64) bool {
	if got == want {
		return true
	}
	if tolerance == 0 {
		return false
	}
	scale := math.Max(math.Abs(got), math.Abs(want))
	if scale == 0 {
		return true
	}
	return math.Abs(got-want)/scale <= tolerance
}
