package column

import "fmt"

func (c Column) MaxHeightForRequirement(nReq float64) (float64, error) {
	if err := c.Validate(); err != nil {
		return 0, err
	}
	if nReq <= 0 {
		return 0, fmt.Errorf("required plate count must be positive, got %.6g", nReq)
	}
	return c.Length / nReq, nil
}

func (c Column) PlateRequirement(nReq float64) (PlateRequirementResult, error) {
	if err := c.Validate(); err != nil {
		return PlateRequirementResult{}, err
	}
	if nReq <= 0 {
		return PlateRequirementResult{}, fmt.Errorf("required plate count must be positive, got %.6g", nReq)
	}
	maxH := c.Length / nReq
	return PlateRequirementResult{NReq: nReq, MaxAllowedHeight: maxH}, nil
}

type PlateRequirementResult struct {
	NReq            float64
	MaxAllowedHeight float64
}

func (r PlateRequirementResult) MaxHeight() float64 {
	return r.MaxAllowedHeight
}

func (r PlateRequirementResult) RequiredCount() float64 {
	return r.NReq
}
