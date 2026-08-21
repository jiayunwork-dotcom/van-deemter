package column

import (
	"van-deemter/internal/hmodel"
)

type InspectResult struct {
	Velocity      float64
	Height        float64
	Plates        float64
	OptimumU      float64
	MinimumH      float64
	MaxPlates     float64
	Range         VelocityRange
	DominantTerm  string
	TermRatio     float64
}

func (c Column) Inspect(u float64) (InspectResult, error) {
	if err := c.Validate(); err != nil {
		return InspectResult{}, err
	}
	if err := hmodel.ValidateVelocity(u); err != nil {
		return InspectResult{}, err
	}
	h, err := c.Model.Height(u)
	if err != nil {
		return InspectResult{}, err
	}
	n, err := c.Plates(u)
	if err != nil {
		return InspectResult{}, err
	}
	uOpt, err := c.Model.OptimumVelocity()
	if err != nil {
		return InspectResult{}, err
	}
	hMin, err := c.Model.MinimumHeight()
	if err != nil {
		return InspectResult{}, err
	}
	nOpt, err := c.Plates(uOpt)
	if err != nil {
		return InspectResult{}, err
	}
	term, err := c.Model.DominantTerm(u)
	if err != nil {
		return InspectResult{}, err
	}
	ratio, err := c.Model.TermRatio(u)
	if err != nil {
		return InspectResult{}, err
	}
	return InspectResult{
		Velocity:     u,
		Height:       h,
		Plates:       n,
		OptimumU:     uOpt,
		MinimumH:     hMin,
		MaxPlates:    nOpt,
		DominantTerm: term,
		TermRatio:    ratio,
	}, nil
}

func (c Column) InspectWithRequirement(u, nReq float64) (InspectResult, error) {
	res, err := c.Inspect(u)
	if err != nil {
		return InspectResult{}, err
	}
	rng, err := c.VelocityRangeForRequirement(nReq)
	if err != nil {
		return InspectResult{}, err
	}
	res.Range = rng
	return res, nil
}

func (r InspectResult) PlateEfficiency() float64 {
	if r.MaxPlates <= 0 {
		return 0
	}
	return r.Plates / r.MaxPlates
}
