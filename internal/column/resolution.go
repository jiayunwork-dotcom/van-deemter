package column

import "fmt"

func retentionError(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

func (c Column) ResolutionAt(u float64, r Retention) (float64, error) {
	if err := c.Validate(); err != nil {
		return 0, err
	}
	if err := r.Validate(); err != nil {
		return 0, err
	}
	n, err := c.Plates(u)
	if err != nil {
		return 0, err
	}
	return r.Resolution(n), nil
}

func (c Column) ResolutionAtOptimum(r Retention) (float64, error) {
	if err := c.Validate(); err != nil {
		return 0, err
	}
	if err := r.Validate(); err != nil {
		return 0, err
	}
	n, _, err := c.PlatesAtOptimum()
	if err != nil {
		return 0, err
	}
	return r.Resolution(n), nil
}

func ResolutionWithoutRetention() error {
	return fmt.Errorf("resolution requires a retention factor k' and selectivity alpha; neither is provided")
}
