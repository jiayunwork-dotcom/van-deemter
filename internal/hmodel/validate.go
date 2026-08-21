package hmodel

import "fmt"

const (
	errFmtNegative = "coefficient %s must be positive, got %.6g"
	errFmtVelocity = "velocity u must be positive, got %.6g"
)

func (m Model) Validate() error {
	if m.A <= 0 {
		return fmt.Errorf(errFmtNegative, "A", m.A)
	}
	if m.B <= 0 {
		return fmt.Errorf(errFmtNegative, "B", m.B)
	}
	if m.C <= 0 {
		return fmt.Errorf(errFmtNegative, "C", m.C)
	}
	return nil
}

func ValidateCoefficientA(a float64) error {
	if a <= 0 {
		return fmt.Errorf(errFmtNegative, "A", a)
	}
	return nil
}

func ValidateCoefficientB(b float64) error {
	if b <= 0 {
		return fmt.Errorf(errFmtNegative, "B", b)
	}
	return nil
}

func ValidateCoefficientC(c float64) error {
	if c <= 0 {
		return fmt.Errorf(errFmtNegative, "C", c)
	}
	return nil
}

func ValidateVelocity(u float64) error {
	if u <= 0 {
		return fmt.Errorf(errFmtVelocity, u)
	}
	return nil
}

func ValidatePositive(label string, v float64) error {
	if v <= 0 {
		return fmt.Errorf("value %s must be positive, got %.6g", label, v)
	}
	return nil
}
