package column

import (
	"fmt"

	"van-deemter/internal/hmodel"
)

type Column struct {
	Length float64
	Model  hmodel.Model
}

func NewColumn(length float64, m hmodel.Model) (Column, error) {
	if err := hmodel.ValidatePositive("length", length); err != nil {
		return Column{}, err
	}
	if err := m.Validate(); err != nil {
		return Column{}, err
	}
	return Column{Length: length, Model: m}, nil
}

func (c Column) Validate() error {
	if c.Length <= 0 {
		return fmt.Errorf("column length must be positive, got %.6g", c.Length)
	}
	return c.Model.Validate()
}

func (c Column) WithLength(length float64) Column {
	c.Length = length
	return c
}

func (c Column) Fields() (float64, float64, float64, float64) {
	a, b, cc := c.Model.Fields()
	return c.Length, a, b, cc
}

func (c Column) String() string {
	return fmt.Sprintf("L=%.6g %s", c.Length, c.Model.String())
}
