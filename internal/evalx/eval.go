package evalx

import (
	"van-deemter/internal/hmodel"
)

type EvalResult struct {
	Name             string
	Velocity         float64
	Height           float64
	Plates           float64
	HasRequirement   bool
	MaxAllowedHeight float64
	MeetsRequirement bool
	HasRetention     bool
	Resolution       float64
}

func Eval(c Case, u float64) (EvalResult, error) {
	if err := hmodel.ValidateVelocity(u); err != nil {
		return EvalResult{}, err
	}
	m, col, err := c.Resolve()
	if err != nil {
		return EvalResult{}, err
	}
	h, err := m.Height(u)
	if err != nil {
		return EvalResult{}, err
	}
	n, err := col.Plates(u)
	if err != nil {
		return EvalResult{}, err
	}
	res := EvalResult{
		Name:     c.NameOrDefault(),
		Velocity: u,
		Height:   h,
		Plates:   n,
	}
	if c.HasRequirement() {
		req, err := c.Requirement()
		if err != nil {
			return EvalResult{}, err
		}
		sufficient, err := col.SufficientVelocity(u, req)
		if err != nil {
			return EvalResult{}, err
		}
		res.HasRequirement = true
		res.MaxAllowedHeight = sufficient.MaxAllowedHeight
		res.MeetsRequirement = sufficient.MeetsRequirement
	}
	if c.HasRetention() {
		r, present, err := c.RetentionOrNone()
		if err != nil {
			return EvalResult{}, err
		}
		if present {
			rs, err := col.ResolutionAt(u, r)
			if err != nil {
				return EvalResult{}, err
			}
			res.HasRetention = true
			res.Resolution = rs
		}
	}
	return res, nil
}

func EvalSummary(c Case, u float64) (EvalResult, error) {
	res, err := Eval(c, u)
	if err != nil {
		return EvalResult{}, err
	}
	res.Name = c.NameOrDefault()
	return res, nil
}

func EvalForRequirement(c Case, u, nReq float64) (EvalResult, error) {
	c.NReq = nReq
	return Eval(c, u)
}
