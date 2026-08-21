package hmodel

type ScanResult struct {
	OptimumVelocity   float64
	MinimumHeight     float64
	Grid              []GridPoint
	GridMinimum       GridPoint
	GridFound         bool
	GridMatchesClosed bool
	DerivativeAtOpt   float64
	DerivativeWithin  bool
}

func Scan(m Model) (ScanResult, error) {
	recordScan("grid")
	if err := m.Validate(); err != nil {
		return ScanResult{}, err
	}
	uOpt, err := m.OptimumVelocity()
	if err != nil {
		return ScanResult{}, err
	}
	hMin, err := m.MinimumHeight()
	if err != nil {
		return ScanResult{}, err
	}
	spec := DefaultGridSpec(m)
	grid, err := GenerateGrid(m, spec)
	if err != nil {
		return ScanResult{}, err
	}
	gridMin, found := GridMinimum(grid)
	delta := 0.001 * uOpt
	deriv, err := m.CentralDerivative(uOpt, delta)
	if err != nil {
		return ScanResult{}, err
	}
	scale := m.B/(uOpt*uOpt) + m.C
	tol := 1e-6 * scale
	within := deriv <= tol && deriv >= -tol
	matches := false
	if found {
		matches = gridMin.H >= hMin*(1-1e-6)
	}
	return ScanResult{
		OptimumVelocity:   uOpt,
		MinimumHeight:     hMin,
		Grid:              grid,
		GridMinimum:       gridMin,
		GridFound:         found,
		GridMatchesClosed: matches,
		DerivativeAtOpt:   deriv,
		DerivativeWithin:  within,
	}, nil
}

func ScanWithTolerance(m Model, derivTol float64) (ScanResult, error) {
	res, err := Scan(m)
	if err != nil {
		return ScanResult{}, err
	}
	res.DerivativeWithin = res.DerivativeAtOpt <= derivTol && res.DerivativeAtOpt >= -derivTol
	return res, nil
}

func ScanFromVelocity(m Model, u float64) (ScanResult, error) {
	if err := ValidateVelocity(u); err != nil {
		return ScanResult{}, err
	}
	res, err := Scan(m)
	if err != nil {
		return ScanResult{}, err
	}
	return res, nil
}
