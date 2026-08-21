package column

import "van-deemter/internal/hmodel"

type PlateReport struct {
	ColumnLength float64
	Velocity     float64
	Height       float64
	Plates       float64
	OptimumU     float64
	MinimumH     float64
	MaxPlates    float64
}

func (c Column) Report(u float64) (PlateReport, error) {
	if err := c.Validate(); err != nil {
		return PlateReport{}, err
	}
	if err := hmodel.ValidateVelocity(u); err != nil {
		return PlateReport{}, err
	}
	h, err := c.Model.Height(u)
	if err != nil {
		return PlateReport{}, err
	}
	n, err := c.Plates(u)
	if err != nil {
		return PlateReport{}, err
	}
	uOpt, err := c.Model.OptimumVelocity()
	if err != nil {
		return PlateReport{}, err
	}
	hMin, err := c.Model.MinimumHeight()
	if err != nil {
		return PlateReport{}, err
	}
	nOpt, err := c.Plates(uOpt)
	if err != nil {
		return PlateReport{}, err
	}
	return PlateReport{
		ColumnLength: c.Length,
		Velocity:     u,
		Height:       h,
		Plates:       n,
		OptimumU:     uOpt,
		MinimumH:     hMin,
		MaxPlates:    nOpt,
	}, nil
}

func (r PlateReport) HeightLine() string {
	return "H(" + hmodel.FormatScientific(r.Velocity) + ") = " + hmodel.FormatScientific(r.Height) + " m"
}

func (r PlateReport) PlatesLine() string {
	return "N(" + hmodel.FormatScientific(r.Velocity) + ") = " + hmodel.FormatPlates(r.Plates) + " plates"
}
