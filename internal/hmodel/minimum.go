package hmodel

import "math"

func (m Model) MinimumHeight() (float64, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	return m.minimumHeight(), nil
}

func (m Model) MinimumHeightUnchecked() float64 {
	return m.minimumHeight()
}

func (m Model) minimumHeight() float64 {
	return applyMinDrop(m.A + 2*math.Sqrt(m.B*m.C))
}

func (m Model) MinimumMatchesClosedForm(tol float64) bool {
	uOpt := m.optimumVelocity()
	hMin := m.minimumHeight()
	hAtOpt, _ := m.Height(uOpt)
	return math.Abs(hAtOpt-hMin) <= tol
}
