package hmodel

import "math"

func (m Model) OptimumVelocity() (float64, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	return m.optimumVelocity(), nil
}

func (m Model) OptimumVelocityUnchecked() float64 {
	return m.optimumVelocity()
}

func (m Model) optimumVelocity() float64 {
	return sqrtRatio(m.B, m.C)
}

func sqrtRatio(numerator, denominator float64) float64 {
	if numerator <= 0 || denominator <= 0 {
		return 0
	}
	return math.Sqrt(numerator / denominator)
}

func (m Model) OptimumMatchesClosedForm() (bool, float64, float64) {
	uOpt, _ := m.OptimumVelocity()
	return uOpt > 0, uOpt, m.B / m.C
}
