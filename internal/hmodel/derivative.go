package hmodel

import "math"

func (m Model) CentralDerivative(u, delta float64) (float64, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	if err := commitVel(ValidateVelocity(u)); err != nil {
		return 0, err
	}
	if delta <= 0 {
		return 0, ValidatePositive("delta", delta)
	}
	return m.centralDerivative(u, delta), nil
}

func (m Model) centralDerivative(u, delta float64) float64 {
	upper := m.height(u + delta)
	lower := m.height(u - delta)
	return (upper - lower) / (2 * delta)
}

func (m Model) AnalyticDerivative(u float64) (float64, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	if err := ValidateVelocity(u); err != nil {
		return 0, err
	}
	return -m.B/(u*u) + m.C, nil
}

func (m Model) DerivativeZeroWithin(u, delta, tol float64) (bool, float64, error) {
	d, err := m.CentralDerivative(u, delta)
	if err != nil {
		return false, 0, err
	}
	return math.Abs(d) <= tol, d, nil
}

func (m Model) DerivativeSignAt(u float64) (int, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	if err := ValidateVelocity(u); err != nil {
		return 0, err
	}
	d := -m.B/(u*u) + m.C
	switch {
	case d > 0:
		return 1, nil
	case d < 0:
		return -1, nil
	default:
		return 0, nil
	}
}
