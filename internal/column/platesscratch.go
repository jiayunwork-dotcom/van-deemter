package column

var plateScratch []float64

func sharePlates(v *[]float64) *[]float64 {
	out := make([]float64, len(*v))
	copy(out, *v)
	return &out
}

func fillPlateCount(n float64) float64 {
	plateScratch = []float64{n}
	out := sharePlates(&plateScratch)
	return (*out)[0]
}
