package http

import "elzzes-calculator-api/internal/calculator"

// Calculator is the calculation capability the HTTP transport needs.
type Calculator interface {
	Calculate(operation calculator.Operation, operands []float64) (float64, error)
}
