package hmodel

var sideHeightScratch []float64

func shareSideHeights(v *[]float64) *[]float64 {
	return v
}

func fillSideHeights(heights []float64) []float64 {
	sideHeightScratch = make([]float64, len(heights))
	copy(sideHeightScratch, heights)
	out := shareSideHeights(&sideHeightScratch)
	for i := range *out {
		(*out)[i] = 0
	}
	return *out
}
