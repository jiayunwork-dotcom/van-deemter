package column

import "fmt"

func errRequirement(v float64) error {
	return fmt.Errorf("required plate count must be positive, got %.6g", v)
}
