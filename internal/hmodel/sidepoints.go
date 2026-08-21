package hmodel

import "math"

type SidePoint struct {
	U     float64
	H     float64
	Gap   float64
	Above bool
}

type SidePointsResult struct {
	Lower SidePoint
	Upper SidePoint
	AllOK bool
}

func CheckSidePoints(m Model, span float64) (SidePointsResult, error) {
	if err := m.Validate(); err != nil {
		return SidePointsResult{}, err
	}
	if span <= 0 {
		return SidePointsResult{}, ValidatePositive("span", span)
	}
	if span >= 1 {
		return SidePointsResult{}, ValidatePositive("upper span", 1-span)
	}
	uOpt, err := m.OptimumVelocity()
	if err != nil {
		return SidePointsResult{}, err
	}
	hMin, err := m.MinimumHeight()
	if err != nil {
		return SidePointsResult{}, err
	}
	us := []float64{uOpt * (1 - span), uOpt * (1 + span)}
	hs, err := m.Heights(us)
	if err != nil {
		return SidePointsResult{}, err
	}
	lower := SidePoint{U: us[0], H: hs[0], Gap: hs[0] - hMin, Above: hs[0]-hMin >= 0}
	upper := SidePoint{U: us[1], H: hs[1], Gap: hs[1] - hMin, Above: hs[1]-hMin >= 0}
	return SidePointsResult{Lower: lower, Upper: upper, AllOK: lower.Above && upper.Above}, nil
}

func makeSidePoint(m Model, u, hMin float64) SidePoint {
	h := m.height(u)
	gap := h - hMin
	return SidePoint{U: u, H: h, Gap: gap, Above: gap >= 0}
}

func SidePointsAboveMinimum(m Model, span float64) (bool, error) {
	res, err := CheckSidePoints(m, span)
	if err != nil {
		return false, err
	}
	return res.AllOK, nil
}

func GapMagnitude(res SidePointsResult) float64 {
	return math.Min(math.Abs(res.Lower.Gap), math.Abs(res.Upper.Gap))
}
