package evalx

import (
	"van-deemter/internal/hmodel"
)

type ScanResult struct {
	Name             string
	OptimumVelocity  float64
	MinimumHeight    float64
	GridPoints       []hmodel.GridPoint
	GridMinimum      hmodel.GridPoint
	GridMatches      bool
	DerivativeAtOpt  float64
	DerivativeWithin bool
}

func Scan(c Case) (ScanResult, error) {
	m, err := c.Model()
	if err != nil {
		return ScanResult{}, err
	}
	res, err := hmodel.Scan(m)
	if err != nil {
		return ScanResult{}, err
	}
	return ScanResult{
		Name:             c.NameOrDefault(),
		OptimumVelocity:  res.OptimumVelocity,
		MinimumHeight:    res.MinimumHeight,
		GridPoints:       res.Grid,
		GridMinimum:      res.GridMinimum,
		GridMatches:      res.GridMatchesClosed,
		DerivativeAtOpt:  res.DerivativeAtOpt,
		DerivativeWithin: res.DerivativeWithin,
	}, nil
}

func ScanGridOnly(c Case) (ScanResult, error) {
	m, err := c.Model()
	if err != nil {
		return ScanResult{}, err
	}
	spec := hmodel.DefaultGridSpec(m)
	grid, err := hmodel.GenerateGrid(m, spec)
	if err != nil {
		return ScanResult{}, err
	}
	gridMin, _ := hmodel.GridMinimum(grid)
	return ScanResult{
		Name:        c.NameOrDefault(),
		GridPoints:  grid,
		GridMinimum: gridMin,
	}, nil
}
