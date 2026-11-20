package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestParseOperation(t *testing.T) {
	valid := []Operation{Add, Subtract, Multiply, Divide, Power, Sqrt, Percent}
	for _, operation := range valid {
		got, ok := ParseOperation(string(operation))
		if !ok || got != operation {
			t.Errorf("ParseOperation(%q) = %q, %v; want %q, true", operation, got, ok, operation)
		}
	}

	invalid := []string{"", "Add", "ADD", " add", "add ", "modulo", "sqrt2"}
	for _, raw := range invalid {
		if got, ok := ParseOperation(raw); ok {
			t.Errorf("ParseOperation(%q) = %q, true; want false", raw, got)
		}
	}
}

func TestOperationArity(t *testing.T) {
	tests := map[Operation]int{
		Add:      2,
		Subtract: 2,
		Multiply: 2,
		Divide:   2,
		Power:    2,
		Sqrt:     1,
		Percent:  1,
	}
	for operation, want := range tests {
		if got := operation.Arity(); got != want {
			t.Errorf("%s.Arity() = %d, want %d", operation, got, want)
		}
	}
	if got := Operation("modulo").Arity(); got != 0 {
		t.Errorf("unknown operation arity = %d, want 0", got)
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		operation Operation
		operands  []float64
		want      float64
		tolerance float64
	}{
		{name: "add", operation: Add, operands: []float64{12.5, 7.5}, want: 20},
		{name: "add negative", operation: Add, operands: []float64{-1.5, 1.5}, want: 0},
		{name: "add fraction artifact", operation: Add, operands: []float64{0.1, 0.2}, want: 0.30000000000000004, tolerance: 1e-12},
		{name: "subtract", operation: Subtract, operands: []float64{12.5, 7.5}, want: 5},
		{name: "subtract zero", operation: Subtract, operands: []float64{0, 0}, want: 0},
		{name: "multiply", operation: Multiply, operands: []float64{12.5, 7.5}, want: 93.75},
		{name: "multiply negative", operation: Multiply, operands: []float64{-4, 2.5}, want: -10},
		{name: "divide", operation: Divide, operands: []float64{10, 4}, want: 2.5},
		{name: "divide repeating", operation: Divide, operands: []float64{1, 3}, want: 1.0 / 3.0, tolerance: 1e-12},
		{name: "divide negative", operation: Divide, operands: []float64{-10, 4}, want: -2.5},
		{name: "power", operation: Power, operands: []float64{2, 10}, want: 1024},
		{name: "power fractional exponent", operation: Power, operands: []float64{9, 0.5}, want: 3},
		{name: "power negative base", operation: Power, operands: []float64{-2, 3}, want: -8},
		{name: "power negative exponent", operation: Power, operands: []float64{2, -3}, want: 0.125},
		{name: "power zero to zero", operation: Power, operands: []float64{0, 0}, want: 1},
		{name: "power large finite", operation: Power, operands: []float64{2, 1023}, want: math.Pow(2, 1023)},
		{name: "sqrt", operation: Sqrt, operands: []float64{144}, want: 12},
		{name: "sqrt zero", operation: Sqrt, operands: []float64{0}, want: 0},
		{name: "sqrt two", operation: Sqrt, operands: []float64{2}, want: math.Sqrt(2), tolerance: 1e-12},
		{name: "sqrt negative zero", operation: Sqrt, operands: []float64{negativeZero()}, want: 0},
		{name: "percent", operation: Percent, operands: []float64{12.5}, want: 0.125},
		{name: "percent negative", operation: Percent, operands: []float64{-12.5}, want: -0.125},
		{name: "percent zero", operation: Percent, operands: []float64{0}, want: 0},
		{name: "percent negative zero", operation: Percent, operands: []float64{negativeZero()}, want: 0},
	}

	calculator := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Calculate(tt.operation, tt.operands)
			if err != nil {
				t.Fatalf("Calculate() error = %v", err)
			}
			if !nearlyEqual(got, tt.want, tt.tolerance) {
				t.Errorf("Calculate() = %v, want %v (tolerance %v)", got, tt.want, tt.tolerance)
			}
			if got == 0 && math.Signbit(got) {
				t.Errorf("Calculate() = -0, want normalized +0")
			}
		})
	}
}

func TestCalculateNegativeZeroNormalized(t *testing.T) {
	calculator := New()
	tests := []struct {
		name      string
		operation Operation
		operands  []float64
	}{
		{name: "sqrt", operation: Sqrt, operands: []float64{negativeZero()}},
		{name: "percent", operation: Percent, operands: []float64{negativeZero()}},
		{name: "multiply", operation: Multiply, operands: []float64{negativeZero(), 5}},
		{name: "add", operation: Add, operands: []float64{negativeZero(), negativeZero()}},
		{name: "divide", operation: Divide, operands: []float64{0, -5}},
		{name: "power", operation: Power, operands: []float64{negativeZero(), 3}},
		{name: "subtract", operation: Subtract, operands: []float64{negativeZero(), 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Calculate(tt.operation, tt.operands)
			if err != nil {
				t.Fatalf("Calculate() error = %v", err)
			}
			if got != 0 || math.Signbit(got) {
				t.Errorf("Calculate() = %v (signbit %v), want +0", got, math.Signbit(got))
			}
		})
	}
}

func TestCalculateDomainErrors(t *testing.T) {
	calculator := New()
	tests := []struct {
		name      string
		operation Operation
		operands  []float64
		wantErr   error
	}{
		{name: "divide by zero", operation: Divide, operands: []float64{10, 0}, wantErr: ErrDivisionByZero},
		{name: "divide by negative zero", operation: Divide, operands: []float64{10, negativeZero()}, wantErr: ErrDivisionByZero},
		{name: "zero divided by zero", operation: Divide, operands: []float64{0, 0}, wantErr: ErrDivisionByZero},
		{name: "add overflow", operation: Add, operands: []float64{math.MaxFloat64, math.MaxFloat64}, wantErr: ErrNonFiniteResult},
		{name: "subtract overflow", operation: Subtract, operands: []float64{-math.MaxFloat64, math.MaxFloat64}, wantErr: ErrNonFiniteResult},
		{name: "multiply overflow", operation: Multiply, operands: []float64{math.MaxFloat64, 2}, wantErr: ErrNonFiniteResult},
		{name: "divide overflow", operation: Divide, operands: []float64{math.MaxFloat64, math.SmallestNonzeroFloat64}, wantErr: ErrNonFiniteResult},
		{name: "power overflow", operation: Power, operands: []float64{10, 400}, wantErr: ErrNonFiniteResult},
		{name: "zero to negative exponent", operation: Power, operands: []float64{0, -1}, wantErr: ErrNonFiniteResult},
		{name: "negative base fractional exponent", operation: Power, operands: []float64{-8, 0.5}, wantErr: ErrNonFiniteResult},
		{name: "sqrt negative", operation: Sqrt, operands: []float64{-1}, wantErr: ErrNegativeSquareRoot},
		{name: "sqrt negative max", operation: Sqrt, operands: []float64{-math.MaxFloat64}, wantErr: ErrNegativeSquareRoot},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Calculate(tt.operation, tt.operands)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Calculate() error = %v, want %v", err, tt.wantErr)
			}
			if got != 0 {
				t.Errorf("Calculate() result = %v on error, want 0", got)
			}
		})
	}
}

func TestCalculateUnderflow(t *testing.T) {
	calculator := New()

	got, err := calculator.Calculate(Divide, []float64{math.SmallestNonzeroFloat64, 10})
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}
	if got != 0 || math.Signbit(got) {
		t.Errorf("Calculate() = %v (signbit %v), want +0", got, math.Signbit(got))
	}

	got, err = calculator.Calculate(Divide, []float64{1e-308, 1e10})
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}
	if !nearlyEqual(got, 1e-318, 1e-12) {
		t.Errorf("Calculate(1e-308 / 1e10) = %v, want 1e-318", got)
	}
}

func TestCalculateExtremeOperands(t *testing.T) {
	calculator := New()

	got, err := calculator.Calculate(Multiply, []float64{math.MaxFloat64, 1})
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}
	if got != math.MaxFloat64 {
		t.Errorf("Calculate(MaxFloat64 * 1) = %v, want MaxFloat64", got)
	}

	got, err = calculator.Calculate(Add, []float64{math.SmallestNonzeroFloat64, 0})
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}
	if got != math.SmallestNonzeroFloat64 {
		t.Errorf("Calculate(5e-324 + 0) = %v, want 5e-324", got)
	}
}

func TestCalculateRejectsUnknownOperation(t *testing.T) {
	if _, err := New().Calculate(Operation("modulo"), []float64{1, 2}); err == nil {
		t.Fatal("Calculate() error = nil, want unknown operation error")
	}
}

func TestCalculateRejectsWrongArity(t *testing.T) {
	calculator := New()
	tests := []struct {
		operation Operation
		operands  []float64
	}{
		{operation: Add},
		{operation: Add, operands: []float64{1}},
		{operation: Add, operands: []float64{1, 2, 3}},
		{operation: Sqrt},
		{operation: Sqrt, operands: []float64{1, 2}},
		{operation: Percent},
	}

	for _, tt := range tests {
		if _, err := calculator.Calculate(tt.operation, tt.operands); err == nil {
			t.Errorf("Calculate(%s, %v) error = nil, want arity error", tt.operation, tt.operands)
		}
	}
}

func negativeZero() float64 {
	return math.Copysign(0, -1)
}

func nearlyEqual(got, want, tolerance float64) bool {
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
