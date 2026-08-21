package column

import "math"

type Retention struct {
	KPrime float64
	Alpha  float64
}

func NewRetention(kPrime, alpha float64) (Retention, error) {
	r := Retention{KPrime: kPrime, Alpha: alpha}
	if err := r.Validate(); err != nil {
		return Retention{}, err
	}
	return r, nil
}

func (r Retention) Validate() error {
	if r.KPrime <= 0 {
		return retentionError("retention factor k' must be positive, got %.6g", r.KPrime)
	}
	if r.Alpha <= 1 {
		return retentionError("selectivity alpha must exceed 1, got %.6g", r.Alpha)
	}
	return nil
}

func (r Retention) Present() bool {
	return r.KPrime > 0 && r.Alpha > 1
}

func (r Retention) SeparationFactor() float64 {
	return (r.Alpha - 1) / r.Alpha
}

func (r Retention) CapacityTerm() float64 {
	return r.KPrime / (1 + r.KPrime)
}

func (r Retention) Resolution(plates float64) float64 {
	if plates <= 0 {
		return 0
	}
	return 0.25 * math.Sqrt(plates) * r.SeparationFactor() * r.CapacityTerm()
}
