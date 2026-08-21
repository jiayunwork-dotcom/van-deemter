package hmodel

import "math"

type CurveMetrics struct {
	MinU float64
	MaxU float64
}

func RangeAroundOptimum(m Model, lowerFactor, upperFactor float64) (CurveMetrics, error) {
	if err := m.Validate(); err != nil {
		return CurveMetrics{}, err
	}
	if lowerFactor <= 0 {
		return CurveMetrics{}, ValidatePositive("lower factor", lowerFactor)
	}
	if upperFactor <= lowerFactor {
		return CurveMetrics{}, ValidatePositive("upper span", upperFactor-lowerFactor)
	}
	uOpt := m.optimumVelocity()
	return CurveMetrics{MinU: lowerFactor * uOpt, MaxU: upperFactor * uOpt}, nil
}

func BuildCurve(m Model, minU, maxU float64, points int) ([]GridPoint, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if minU <= 0 {
		return nil, ValidatePositive("curve minimum", minU)
	}
	if maxU <= minU {
		return nil, ValidatePositive("curve span", maxU-minU)
	}
	if points < 2 {
		return nil, ValidatePositive("curve points", float64(points))
	}
	step := (maxU - minU) / float64(points-1)
	curve := make([]GridPoint, 0, points)
	for i := 0; i < points; i++ {
		u := minU + float64(i)*step
		curve = append(curve, GridPoint{U: u, H: m.height(u)})
	}
	return curve, nil
}

func CurveAt(velocities []float64, m Model) ([]GridPoint, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	points := make([]GridPoint, 0, len(velocities))
	for _, u := range velocities {
		if err := ValidateVelocity(u); err != nil {
			return nil, err
		}
		points = append(points, GridPoint{U: u, H: m.height(u)})
	}
	return points, nil
}

func CurveMinimumHeight(curve []GridPoint) (float64, bool) {
	best, ok := GridMinimum(curve)
	if !ok {
		return 0, false
	}
	return best.H, true
}

func CurveMinimaTolerance(hCurve, hMin, tol float64) bool {
	return math.Abs(hCurve-hMin) <= tol
}
