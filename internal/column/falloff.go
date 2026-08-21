package column

import "van-deemter/internal/hmodel"

type FalloffPoint struct {
	U     float64
	Plates float64
	Below bool
}

type FalloffResult struct {
	OptimumPlates float64
	Lower         FalloffPoint
	Upper         FalloffPoint
	AllBelow      bool
}

func (c Column) FalloffAroundOptimum(span float64) (FalloffResult, error) {
	if err := c.Validate(); err != nil {
		return FalloffResult{}, err
	}
	if span <= 0 {
		return FalloffResult{}, hmodel.ValidatePositive("span", span)
	}
	if span >= 1 {
		return FalloffResult{}, hmodel.ValidatePositive("upper span", 1-span)
	}
	uOpt, err := c.Model.OptimumVelocity()
	if err != nil {
		return FalloffResult{}, err
	}
	nOpt, err := c.Plates(uOpt)
	if err != nil {
		return FalloffResult{}, err
	}
	lowerU := uOpt * (1 - span)
	upperU := uOpt * (1 + span)
	lowerN, err := c.Plates(lowerU)
	if err != nil {
		return FalloffResult{}, err
	}
	upperN, err := c.Plates(upperU)
	if err != nil {
		return FalloffResult{}, err
	}
	return FalloffResult{
		OptimumPlates: nOpt,
		Lower:         FalloffPoint{U: lowerU, Plates: lowerN, Below: lowerN < nOpt},
		Upper:         FalloffPoint{U: upperU, Plates: upperN, Below: upperN < nOpt},
		AllBelow:      lowerN < nOpt && upperN < nOpt,
	}, nil
}

func (c Column) PlatesAboveRequirement(u, nReq float64) (bool, error) {
	n, err := c.Plates(u)
	if err != nil {
		return false, err
	}
	return n >= nReq, nil
}

func (c Column) RatioToOptimum(u float64) (float64, error) {
	if err := c.Validate(); err != nil {
		return 0, err
	}
	if err := hmodel.ValidateVelocity(u); err != nil {
		return 0, err
	}
	n, err := c.Plates(u)
	if err != nil {
		return 0, err
	}
	nOpt, _, err := c.PlatesAtOptimum()
	if err != nil {
		return 0, err
	}
	if nOpt == 0 {
		return 0, nil
	}
	return n / nOpt, nil
}
