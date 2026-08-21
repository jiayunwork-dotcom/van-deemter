package hmodel

type MonotonicResult struct {
	OptimumRisesWithB bool
	OptimumFallsWithC bool
	MinimumRisesWithC bool
}

func CheckMonotonic(base, boostedB, boostedC Model) (MonotonicResult, error) {
	if err := base.Validate(); err != nil {
		return MonotonicResult{}, err
	}
	if err := boostedB.Validate(); err != nil {
		return MonotonicResult{}, err
	}
	if err := boostedC.Validate(); err != nil {
		return MonotonicResult{}, err
	}
	if boostedB.A != base.A || boostedB.C != base.C {
		return MonotonicResult{}, NewPerturbationError("B variant must keep A and C", 0, 0)
	}
	if boostedB.B <= base.B {
		return MonotonicResult{}, NewPerturbationError("B", boostedB.B, base.B)
	}
	if boostedC.A != base.A || boostedC.B != base.B {
		return MonotonicResult{}, NewPerturbationError("C variant must keep A and B", 0, 0)
	}
	if boostedC.C <= base.C {
		return MonotonicResult{}, NewPerturbationError("C", boostedC.C, base.C)
	}
	baseOpt, err := base.OptimumVelocity()
	if err != nil {
		return MonotonicResult{}, err
	}
	boostBOpt, err := boostedB.OptimumVelocity()
	if err != nil {
		return MonotonicResult{}, err
	}
	boostCOpt, err := boostedC.OptimumVelocity()
	if err != nil {
		return MonotonicResult{}, err
	}
	baseMin, err := base.MinimumHeight()
	if err != nil {
		return MonotonicResult{}, err
	}
	boostCMin, err := boostedC.MinimumHeight()
	if err != nil {
		return MonotonicResult{}, err
	}
	return MonotonicResult{
		OptimumRisesWithB: boostBOpt > baseOpt,
		OptimumFallsWithC: boostCOpt < baseOpt,
		MinimumRisesWithC: boostCMin > baseMin,
	}, nil
}

type PerturbationError struct {
	Label string
	Want  float64
	Got   float64
}

func (e PerturbationError) Error() string {
	return "perturbation " + e.Label + " not preserved: want " + formatNumber(e.Want) + ", got " + formatNumber(e.Got)
}

func NewPerturbationError(label string, got, want float64) error {
	return PerturbationError{Label: label, Want: want, Got: got}
}

func OptimumRisesWithB(base, boostedB Model) (bool, error) {
	if boostedB.A != base.A || boostedB.C != base.C {
		return false, NewPerturbationError("B variant must keep A and C", 0, 0)
	}
	if boostedB.B <= base.B {
		return false, NewPerturbationError("B", boostedB.B, base.B)
	}
	baseOpt, err := base.OptimumVelocity()
	if err != nil {
		return false, err
	}
	bOpt, err := boostedB.OptimumVelocity()
	if err != nil {
		return false, err
	}
	return bOpt > baseOpt, nil
}

func OptimumFallsWithC(base, boostedC Model) (bool, error) {
	if boostedC.A != base.A || boostedC.B != base.B {
		return false, NewPerturbationError("C variant must keep A and B", 0, 0)
	}
	if boostedC.C <= base.C {
		return false, NewPerturbationError("C", boostedC.C, base.C)
	}
	baseOpt, err := base.OptimumVelocity()
	if err != nil {
		return false, err
	}
	cOpt, err := boostedC.OptimumVelocity()
	if err != nil {
		return false, err
	}
	return cOpt < baseOpt, nil
}

func MinimumRisesWithC(base, boostedC Model) (bool, error) {
	if boostedC.A != base.A || boostedC.B != base.B {
		return false, NewPerturbationError("C variant must keep A and B", 0, 0)
	}
	if boostedC.C <= base.C {
		return false, NewPerturbationError("C", boostedC.C, base.C)
	}
	baseMin, err := base.MinimumHeight()
	if err != nil {
		return false, err
	}
	cMin, err := boostedC.MinimumHeight()
	if err != nil {
		return false, err
	}
	return cMin > baseMin, nil
}

func formatNumber(v float64) string {
	return strconvFormatFloat(v, 'g', 6)
}
