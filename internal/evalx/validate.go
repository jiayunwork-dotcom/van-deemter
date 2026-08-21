package evalx

import (
	"fmt"
	"strings"

	"van-deemter/internal/hmodel"
)

func ValidateCase(c Case) error {
	errs := ValidateCaseAll(c)
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errs, "; "))
}

func ValidateCaseAll(c Case) []string {
	var errs []string
	if c.A <= 0 {
		errs = append(errs, fmt.Sprintf("coefficient A must be positive, got %.6g", c.A))
	}
	if c.B <= 0 {
		errs = append(errs, fmt.Sprintf("coefficient B must be positive, got %.6g", c.B))
	}
	if c.C <= 0 {
		errs = append(errs, fmt.Sprintf("coefficient C must be positive, got %.6g", c.C))
	}
	if c.Length <= 0 {
		errs = append(errs, fmt.Sprintf("column length must be positive, got %.6g", c.Length))
	}
	if c.NReq < 0 {
		errs = append(errs, fmt.Sprintf("required plate count must be positive, got %.6g", c.NReq))
	}
	if c.HasRetention() {
		if c.KPrime <= 0 {
			errs = append(errs, "retention factor k' must be positive when given")
		}
		if c.Alpha <= 1 {
			errs = append(errs, "selectivity alpha must exceed 1 when given")
		}
	}
	return errs
}

func ValidateVelocityForCase(c Case, u float64) error {
	if err := hmodel.ValidateVelocity(u); err != nil {
		return err
	}
	return ValidateCase(c)
}

func (c Case) ValidateModelOnly() error {
	_, err := c.Model()
	return err
}

func (c Case) ValidateColumnOnly() error {
	_, err := c.Column()
	return err
}
