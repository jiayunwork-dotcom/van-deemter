package column

import (
	"van-deemter/internal/hmodel"
)

type VelocityRange struct {
	Feasible bool
	Lower    float64
	Upper    float64
	MinHeight float64
}

func (c Column) VelocityRangeForRequirement(nReq float64) (VelocityRange, error) {
	if err := c.Validate(); err != nil {
		return VelocityRange{}, err
	}
	if nReq <= 0 {
		return VelocityRange{}, errRequirement(nReq)
	}
	hMax := c.Length / nReq
	hMin, err := c.Model.MinimumHeight()
	if err != nil {
		return VelocityRange{}, err
	}
	if hMax < hMin {
		return VelocityRange{Feasible: false, MinHeight: hMin}, nil
	}
	lo, hi, err := c.Model.InverseVelocities(hMax)
	if err != nil {
		return VelocityRange{}, err
	}
	return VelocityRange{Feasible: true, Lower: lo, Upper: hi, MinHeight: hMin}, nil
}

func (r VelocityRange) Contains(u float64) bool {
	return r.Feasible && u >= r.Lower && u <= r.Upper
}

func (r VelocityRange) Width() float64 {
	if !r.Feasible {
		return 0
	}
	return r.Upper - r.Lower
}

func (r VelocityRange) MarginAt(u float64) float64 {
	if !r.Feasible {
		return 0
	}
	if u < r.Lower {
		return r.Lower - u
	}
	if u > r.Upper {
		return u - r.Upper
	}
	return 0
}

func (c Column) VelocityRangeByHeight(hMax float64) (VelocityRange, error) {
	if err := c.Validate(); err != nil {
		return VelocityRange{}, err
	}
	if err := hmodel.ValidatePositive("target plate height", hMax); err != nil {
		return VelocityRange{}, err
	}
	hMin, err := c.Model.MinimumHeight()
	if err != nil {
		return VelocityRange{}, err
	}
	if hMax < hMin {
		return VelocityRange{Feasible: false, MinHeight: hMin}, nil
	}
	lo, hi, err := c.Model.InverseVelocities(hMax)
	if err != nil {
		return VelocityRange{}, err
	}
	return VelocityRange{Feasible: true, Lower: lo, Upper: hi, MinHeight: hMin}, nil
}
