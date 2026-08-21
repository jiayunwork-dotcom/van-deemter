package column

import "van-deemter/internal/hmodel"

type SufficiencyResult struct {
	Velocity          float64
	Height            float64
	Plates            float64
	MaxAllowedHeight  float64
	MeetsRequirement  bool
}

func (c Column) SufficientVelocity(u, nReq float64) (SufficiencyResult, error) {
	if err := c.Validate(); err != nil {
		return SufficiencyResult{}, err
	}
	if err := hmodel.ValidateVelocity(u); err != nil {
		return SufficiencyResult{}, err
	}
	if nReq <= 0 {
		return SufficiencyResult{}, errRequirement(nReq)
	}
	h, err := c.Model.Height(u)
	if err != nil {
		return SufficiencyResult{}, err
	}
	maxH := c.Length / nReq
	return SufficiencyResult{
		Velocity:         u,
		Height:           h,
		Plates:           c.Length / h,
		MaxAllowedHeight: maxH,
		MeetsRequirement: h <= maxH,
	}, nil
}

func (c Column) SufficientAtOptimum(nReq float64) (SufficiencyResult, error) {
	if err := c.Validate(); err != nil {
		return SufficiencyResult{}, err
	}
	uOpt, err := c.Model.OptimumVelocity()
	if err != nil {
		return SufficiencyResult{}, err
	}
	return c.SufficientVelocity(uOpt, nReq)
}

func (c Column) MarginForRequirement(u, nReq float64) (float64, error) {
	res, err := c.SufficientVelocity(u, nReq)
	if err != nil {
		return 0, err
	}
	return res.MaxAllowedHeight - res.Height, nil
}
