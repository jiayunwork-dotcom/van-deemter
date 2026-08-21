package hmodel

type Range struct {
	Lower   float64
	Upper   float64
	Optimum float64
}

func (m Model) RecommendedRange() (Range, error) {
	if err := m.Validate(); err != nil {
		return Range{}, err
	}
	uOpt, err := m.OptimumVelocity()
	if err != nil {
		return Range{}, err
	}
	return Range{Lower: uOpt / 5, Upper: uOpt * 5, Optimum: uOpt}, nil
}

func (m Model) BalancedVelocity() (float64, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	return m.optimumVelocity(), nil
}

func (m Model) DominantTerm(u float64) (string, error) {
	if err := m.Validate(); err != nil {
		return "", err
	}
	if err := ValidateVelocity(u); err != nil {
		return "", err
	}
	diffusion := m.B / u
	resistance := m.C * u
	if diffusion > resistance {
		return "longitudinal diffusion", nil
	}
	if resistance > diffusion {
		return "mass-transfer resistance", nil
	}
	return "balanced", nil
}

func (m Model) TermRatio(u float64) (float64, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	if err := ValidateVelocity(u); err != nil {
		return 0, err
	}
	resistance := m.C * u
	diffusion := m.B / u
	if diffusion <= 0 {
		return 0, nil
	}
	return resistance / diffusion, nil
}

func (m Model) HeightBudget(u float64) (Budget, error) {
	if err := m.Validate(); err != nil {
		return Budget{}, err
	}
	if err := ValidateVelocity(u); err != nil {
		return Budget{}, err
	}
	h := m.height(u)
	return Budget{
		Eddy:   m.A,
		Axial:  m.B / u,
		Resist: m.C * u,
		Total:  h,
	}, nil
}

type Budget struct {
	Eddy   float64
	Axial  float64
	Resist float64
	Total  float64
}

func (b Budget) Fractions() (float64, float64, float64) {
	if b.Total <= 0 {
		return 0, 0, 0
	}
	return b.Eddy / b.Total, b.Axial / b.Total, b.Resist / b.Total
}

func (b Budget) ResistShare() float64 {
	_, _, resist := b.Fractions()
	return resist
}

func (b Budget) AxialShare() float64 {
	_, axial, _ := b.Fractions()
	return axial
}
