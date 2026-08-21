package hmodel

import "fmt"

func (m Model) InverseVelocities(targetH float64) (float64, float64, error) {
	if err := m.Validate(); err != nil {
		return 0, 0, err
	}
	if err := ValidatePositive("target plate height", targetH); err != nil {
		return 0, 0, err
	}
	hMin, err := m.MinimumHeight()
	if err != nil {
		return 0, 0, err
	}
	if targetH < hMin {
		return 0, 0, fmt.Errorf("target H %.6g below H_min %.6g; no velocity reaches it", targetH, hMin)
	}
	excess := targetH - m.A
	lo, hi, ok := SolvePositiveRoots(m.C, -excess, m.B)
	if !ok {
		return 0, 0, fmt.Errorf("target H %.6g gives no positive velocity roots", targetH)
	}
	return lo, hi, nil
}

func (m Model) InverseVelocitiesAroundOptimum(targetH float64) (VelocityPair, error) {
	lo, hi, err := m.InverseVelocities(targetH)
	if err != nil {
		return VelocityPair{}, err
	}
	uOpt, err := m.OptimumVelocity()
	if err != nil {
		return VelocityPair{}, err
	}
	return VelocityPair{Lower: lo, Upper: hi, Optimum: uOpt}, nil
}

type VelocityPair struct {
	Lower   float64
	Upper   float64
	Optimum float64
}

func (p VelocityPair) Contains(u float64) bool {
	return u >= p.Lower && u <= p.Upper
}

func (p VelocityPair) Width() float64 {
	return p.Upper - p.Lower
}

func (p VelocityPair) SpanRatio() float64 {
	if p.Lower <= 0 {
		return 0
	}
	return p.Upper / p.Lower
}

func (m Model) VelocityForTargetHeight(targetH float64) (VelocityPair, error) {
	return m.InverseVelocitiesAroundOptimum(targetH)
}
