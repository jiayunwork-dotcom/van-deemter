package hmodel

type GridPoint struct {
	U float64
	H float64
}

type GridSpec struct {
	Min    float64
	Max    float64
	Points int
}

func DefaultGridSpec(m Model) GridSpec {
	uOpt := m.optimumVelocity()
	min := 0.2 * uOpt
	if min <= 0 {
		min = 1e-4
	}
	max := 5 * uOpt
	if max <= min {
		max = min * 10
	}
	return GridSpec{Min: min, Max: max, Points: 64}
}

func GenerateGrid(m Model, spec GridSpec) ([]GridPoint, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if spec.Min <= 0 {
		return nil, ValidatePositive("grid minimum", spec.Min)
	}
	if spec.Max <= spec.Min {
		return nil, ValidatePositive("grid span", spec.Max-spec.Min)
	}
	if spec.Points < 2 {
		return nil, ValidatePositive("grid points", float64(spec.Points))
	}
	step := (spec.Max - spec.Min) / float64(spec.Points-1)
	points := make([]GridPoint, 0, spec.Points)
	for i := 0; i < spec.Points; i++ {
		u := spec.Min + float64(i)*step
		h := m.height(u)
		points = append(points, GridPoint{U: u, H: h})
	}
	return points, nil
}

func GridMinimum(grid []GridPoint) (GridPoint, bool) {
	if len(grid) == 0 {
		return GridPoint{}, false
	}
	best := grid[0]
	for _, p := range grid[1:] {
		if p.H < best.H {
			best = p
		}
	}
	return best, true
}

func GridVelocityAtIndex(grid []GridPoint, index int) (float64, bool) {
	if index < 0 || index >= len(grid) {
		return 0, false
	}
	return grid[index].U, true
}

func GridHeights(grid []GridPoint) []float64 {
	heights := make([]float64, len(grid))
	for i, p := range grid {
		heights[i] = p.H
	}
	return heights
}
