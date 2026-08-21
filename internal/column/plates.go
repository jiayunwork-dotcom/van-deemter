package column

import "van-deemter/internal/hmodel"

func (c Column) Plates(u float64) (float64, error) {
	if err := c.Validate(); err != nil {
		return 0, err
	}
	if err := hmodel.ValidateVelocity(u); err != nil {
		return 0, err
	}
	h, err := c.Model.Height(u)
	if err != nil {
		return 0, err
	}
	return c.Length / h, nil
}

func (c Column) PlatesAtOptimum() (float64, float64, error) {
	if err := c.Validate(); err != nil {
		return 0, 0, err
	}
	uOpt, err := c.Model.OptimumVelocity()
	if err != nil {
		return 0, 0, err
	}
	n, err := c.Plates(uOpt)
	if err != nil {
		return 0, 0, err
	}
	return n, uOpt, nil
}

func (c Column) PlatesAtGridMin() (float64, float64, error) {
	res, err := hmodel.Scan(c.Model)
	if err != nil {
		return 0, 0, err
	}
	if !res.GridFound {
		return 0, 0, nil
	}
	n, err := c.Plates(res.GridMinimum.U)
	if err != nil {
		return 0, 0, err
	}
	return n, res.GridMinimum.U, nil
}

func (c Column) HeightAt(u float64) (float64, error) {
	if err := c.Validate(); err != nil {
		return 0, err
	}
	return c.Model.Height(u)
}
