package calculator

import "errors"

var (
	ErrDivisionByZero     = errors.New("cannot divide by zero")
	ErrNegativeSquareRoot = errors.New("square root requires a non-negative operand")
	ErrNonFiniteResult    = errors.New("result is not finite")
)
