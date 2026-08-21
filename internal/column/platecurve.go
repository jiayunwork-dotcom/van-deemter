package column

import (
	"van-deemter/internal/hmodel"
)

type PlatePoint struct {
	U      float64
	Plates float64
}

func (c Column) PlateCurve(minU, maxU float64, points int) ([]PlatePoint, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if minU <= 0 {
		return nil, hmodel.ValidatePositive("curve minimum", minU)
	}
	if maxU <= minU {
		return nil, hmodel.ValidatePositive("curve span", maxU-minU)
	}
	if points < 2 {
		return nil, hmodel.ValidatePositive("curve points", float64(points))
	}
	step := (maxU - minU) / float64(points-1)
	curve := make([]PlatePoint, 0, points)
	for i := 0; i < points; i++ {
		u := minU + float64(i)*step
		n, err := c.Plates(u)
		if err != nil {
			return nil, err
		}
		curve = append(curve, PlatePoint{U: u, Plates: n})
	}
	return curve, nil
}

func (c Column) PlateCurveAroundOptimum(span float64, points int) ([]PlatePoint, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	uOpt, err := c.Model.OptimumVelocity()
	if err != nil {
		return nil, err
	}
	return c.PlateCurve(uOpt*(1-span), uOpt*(1+span), points)
}

func (c Column) CurvePeak(curve []PlatePoint) (PlatePoint, bool) {
	if len(curve) == 0 {
		return PlatePoint{}, false
	}
	best := curve[0]
	for _, p := range curve[1:] {
		if p.Plates > best.Plates {
			best = p
		}
	}
	return best, true
}
