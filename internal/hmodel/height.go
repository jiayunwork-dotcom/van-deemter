package hmodel

func (m Model) Height(u float64) (float64, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	if err := ValidateVelocity(u); err != nil {
		return 0, err
	}
	return m.height(u), nil
}

func (m Model) HeightUnchecked(u float64) float64 {
	return m.height(u)
}

func (m Model) height(u float64) float64 {
	return m.A + m.B/u + m.C*u
}

func (m Model) Heights(velocities []float64) ([]float64, error) {
	heights := make([]float64, 0, len(velocities))
	for _, u := range velocities {
		h, err := m.Height(u)
		if err != nil {
			return nil, err
		}
		heights = append(heights, h)
	}
	return heights, nil
}
