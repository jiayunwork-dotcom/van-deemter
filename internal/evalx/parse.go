package evalx

import (
	"fmt"

	"van-deemter/internal/column"
	"van-deemter/internal/hmodel"
)

func (c Case) Model() (hmodel.Model, error) {
	return hmodel.NewModel(c.A, c.B, c.C)
}

func (c Case) Column() (column.Column, error) {
	m, err := c.Model()
	if err != nil {
		return column.Column{}, err
	}
	return column.NewColumn(c.Length, m)
}

func (c Case) RetentionOrNone() (column.Retention, bool, error) {
	if c.KPrime == 0 && c.Alpha == 0 {
		return column.Retention{}, false, nil
	}
	if c.KPrime <= 0 {
		return column.Retention{}, false, fmt.Errorf("selectivity alpha given but retention factor k' missing")
	}
	if c.Alpha <= 1 {
		return column.Retention{}, false, fmt.Errorf("retention factor k' given but selectivity alpha missing or not exceeding 1")
	}
	r, err := column.NewRetention(c.KPrime, c.Alpha)
	if err != nil {
		return column.Retention{}, false, err
	}
	return r, true, nil
}

func (c Case) Resolve() (hmodel.Model, column.Column, error) {
	m, err := c.Model()
	if err != nil {
		return hmodel.Model{}, column.Column{}, err
	}
	col, err := column.NewColumn(c.Length, m)
	if err != nil {
		return hmodel.Model{}, column.Column{}, err
	}
	return m, col, nil
}

func (c Case) Requirement() (float64, error) {
	if c.NReq <= 0 {
		return 0, fmt.Errorf("required plate count must be positive, got %.6g", c.NReq)
	}
	return c.NReq, nil
}
