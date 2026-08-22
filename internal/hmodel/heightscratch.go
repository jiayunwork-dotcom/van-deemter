package hmodel

var sideHeightScratch []float64

func shareSideHeights(v *[]float64) *[]float64 {
	out := make([]float64, len(*v))
	copy(out, *v)
	return &out
}

func fillSideHeights(heights []float64) []float64 {
	sideHeightScratch = make([]float64, len(heights))
	copy(sideHeightScratch, heights)
	out := shareSideHeights(&sideHeightScratch)
	return *out
}
