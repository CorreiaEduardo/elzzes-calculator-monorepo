// Package calculator is the pure calculator core: operation vocabulary, arity,
// arithmetic, and domain errors. It must not import transport, HTTP, logging, or
// configuration packages.
package calculator

import (
	"fmt"
	"math"
)

type Operation string

const (
	Add      Operation = "add"
	Subtract Operation = "subtract"
	Multiply Operation = "multiply"
	Divide   Operation = "divide"
	Power    Operation = "power"
	Sqrt     Operation = "sqrt"
	Percent  Operation = "percent"
)

type operationSpec struct {
	arity int
	apply func([]float64) (float64, error)
}

var operations = map[Operation]operationSpec{
	Add:      {arity: 2, apply: add},
	Subtract: {arity: 2, apply: subtract},
	Multiply: {arity: 2, apply: multiply},
	Divide:   {arity: 2, apply: divide},
	Power:    {arity: 2, apply: power},
	Sqrt:     {arity: 1, apply: sqrt},
	Percent:  {arity: 1, apply: percent},
}

func ParseOperation(raw string) (Operation, bool) {
	operation := Operation(raw)
	_, ok := operations[operation]
	return operation, ok
}

func (o Operation) Arity() int {
	return operations[o].arity
}

type Calculator struct{}

func New() *Calculator {
	return &Calculator{}
}

func (c *Calculator) Calculate(operation Operation, operands []float64) (float64, error) {
	spec, ok := operations[operation]
	if !ok {
		return 0, fmt.Errorf("unknown operation %q", operation)
	}
	if len(operands) != spec.arity {
		return 0, fmt.Errorf("%s expects %d operands, got %d", operation, spec.arity, len(operands))
	}

	result, err := spec.apply(operands)
	if err != nil {
		return 0, err
	}
	return finite(result)
}

func add(operands []float64) (float64, error) {
	return operands[0] + operands[1], nil
}

func subtract(operands []float64) (float64, error) {
	return operands[0] - operands[1], nil
}

func multiply(operands []float64) (float64, error) {
	return operands[0] * operands[1], nil
}

func divide(operands []float64) (float64, error) {
	if operands[1] == 0 {
		return 0, ErrDivisionByZero
	}
	return operands[0] / operands[1], nil
}

func power(operands []float64) (float64, error) {
	return math.Pow(operands[0], operands[1]), nil
}

func sqrt(operands []float64) (float64, error) {
	if operands[0] < 0 {
		return 0, ErrNegativeSquareRoot
	}
	return math.Sqrt(operands[0]), nil
}

func percent(operands []float64) (float64, error) {
	return operands[0] / 100, nil
}

func finite(result float64) (float64, error) {
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrNonFiniteResult
	}
	if result == 0 {
		return 0, nil
	}
	return result, nil
}
