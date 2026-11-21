package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	nethttp "net/http"
	"sort"
	"strconv"
	"strings"

	"elzzes-calculator-api/internal/calculator"
)

const maxBodyBytes int64 = 4 << 10

type calculateResponse struct {
	Operation calculator.Operation `json:"operation"`
	Operands  []float64            `json:"operands"`
	Result    float64              `json:"result"`
}

type requestFailure struct {
	status  int
	code    ErrorCode
	message string
	cause   error
}

func (f *requestFailure) Error() string {
	return f.message
}

func (f *requestFailure) write(w nethttp.ResponseWriter) {
	writeError(w, f.status, f.code, f.message)
}

func calculateHandler(logger Logger, calc Calculator) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, req *nethttp.Request) {
		operation, operands, err := parseCalculateRequest(w, req)
		if err != nil {
			var failure *requestFailure
			if !errors.As(err, &failure) {
				logger.Error("calculate request failed", "error", err)
				writeError(w, nethttp.StatusInternalServerError, CodeInternalError, msgInternalError)
				return
			}
			if failure.cause != nil {
				logger.Error("calculate request read failed", "error", failure.cause)
			}
			failure.write(w)
			return
		}

		result, err := calc.Calculate(operation, operands)
		if err != nil {
			if writeCalculatorError(w, err) {
				return
			}
			logger.Error("calculate failed", "operation", operation, "error", err)
			writeError(w, nethttp.StatusInternalServerError, CodeInternalError, msgInternalError)
			return
		}

		writeJSON(w, nethttp.StatusOK, calculateResponse{
			Operation: operation,
			Operands:  operands,
			Result:    result,
		})
	})
}

func parseCalculateRequest(w nethttp.ResponseWriter, req *nethttp.Request) (calculator.Operation, []float64, error) {
	if err := checkContentType(req); err != nil {
		return "", nil, err
	}

	body, err := readBody(w, req)
	if err != nil {
		return "", nil, err
	}

	fields, err := decodeFields(body)
	if err != nil {
		return "", nil, err
	}

	operation, err := parseOperation(fields)
	if err != nil {
		return "", nil, err
	}

	operands, err := parseOperands(fields, operation)
	if err != nil {
		return "", nil, err
	}

	return operation, operands, nil
}

func checkContentType(req *nethttp.Request) error {
	mediaType, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil || !strings.EqualFold(mediaType, "application/json") {
		return &requestFailure{
			status:  nethttp.StatusUnsupportedMediaType,
			code:    CodeInvalidContentType,
			message: "Content-Type must be application/json",
		}
	}
	return nil
}

func readBody(w nethttp.ResponseWriter, req *nethttp.Request) ([]byte, error) {
	if req.ContentLength > maxBodyBytes {
		return nil, payloadTooLargeFailure()
	}

	body := req.Body
	if body == nil {
		body = nethttp.NoBody
	}

	raw, err := io.ReadAll(nethttp.MaxBytesReader(w, body, maxBodyBytes))
	if err != nil {
		var maxErr *nethttp.MaxBytesError
		if errors.As(err, &maxErr) {
			return nil, payloadTooLargeFailure()
		}
		failure := invalidJSONFailure()
		failure.cause = err
		return nil, failure
	}
	return raw, nil
}

func decodeFields(body []byte) (map[string]json.RawMessage, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil || fields == nil {
		return nil, invalidJSONFailure()
	}

	unknown := make([]string, 0)
	for key := range fields {
		if key != "operation" && key != "operands" {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return nil, &requestFailure{
			status:  nethttp.StatusBadRequest,
			code:    CodeUnknownField,
			message: fmt.Sprintf("unknown field %q", unknown[0]),
		}
	}
	return fields, nil
}

func parseOperation(fields map[string]json.RawMessage) (calculator.Operation, error) {
	raw, ok := fields["operation"]
	if !ok {
		return "", &requestFailure{
			status:  nethttp.StatusBadRequest,
			code:    CodeUnknownOperation,
			message: "operation is required",
		}
	}

	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return "", operationTypeFailure()
	}
	rawOperation, ok := decoded.(string)
	if !ok {
		return "", operationTypeFailure()
	}

	operation, ok := calculator.ParseOperation(rawOperation)
	if !ok {
		return "", &requestFailure{
			status:  nethttp.StatusBadRequest,
			code:    CodeUnknownOperation,
			message: fmt.Sprintf("unknown operation %q", rawOperation),
		}
	}
	return operation, nil
}

func parseOperands(fields map[string]json.RawMessage, operation calculator.Operation) ([]float64, error) {
	raw, ok := fields["operands"]
	if !ok {
		return nil, operandCountFailure(operation)
	}

	var decoded any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&decoded); err != nil {
		return nil, operandsTypeFailure()
	}

	values, ok := decoded.([]any)
	if !ok {
		return nil, operandsTypeFailure()
	}
	if len(values) != operation.Arity() {
		return nil, operandCountFailure(operation)
	}

	operands := make([]float64, len(values))
	for i, rawValue := range values {
		number, ok := rawValue.(json.Number)
		if !ok {
			return nil, invalidOperandFailure(i)
		}
		value, err := strconv.ParseFloat(number.String(), 64)
		if err != nil || math.IsInf(value, 0) || math.IsNaN(value) {
			return nil, invalidOperandFailure(i)
		}
		operands[i] = value
	}
	return operands, nil
}

func operationTypeFailure() *requestFailure {
	return &requestFailure{
		status:  nethttp.StatusBadRequest,
		code:    CodeUnknownOperation,
		message: "operation must be a string",
	}
}

func operandsTypeFailure() *requestFailure {
	return &requestFailure{
		status:  nethttp.StatusBadRequest,
		code:    CodeInvalidOperand,
		message: "operands must be an array of numbers",
	}
}

func invalidOperandFailure(index int) *requestFailure {
	return &requestFailure{
		status:  nethttp.StatusBadRequest,
		code:    CodeInvalidOperand,
		message: fmt.Sprintf("operands[%d] must be a JSON number", index),
	}
}

func operandCountFailure(operation calculator.Operation) *requestFailure {
	arity := operation.Arity()
	message := fmt.Sprintf("%s expects %d operands", operation, arity)
	if arity == 1 {
		message = fmt.Sprintf("%s expects 1 operand", operation)
	}
	return &requestFailure{
		status:  nethttp.StatusBadRequest,
		code:    CodeInvalidOperandCount,
		message: message,
	}
}

func invalidJSONFailure() *requestFailure {
	return &requestFailure{
		status:  nethttp.StatusBadRequest,
		code:    CodeInvalidJSON,
		message: "request body is not valid JSON",
	}
}

func payloadTooLargeFailure() *requestFailure {
	return &requestFailure{
		status:  nethttp.StatusRequestEntityTooLarge,
		code:    CodePayloadTooLarge,
		message: "request body exceeds 4 KiB",
	}
}
