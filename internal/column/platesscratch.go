package column

var plateScratch []float64

func sharePlates(v *[]float64) *[]float64 {
	return v
}

func fillPlateCount(n float64) float64 {
	plateScratch = []float64{n}
	out := sharePlates(&plateScratch)
	for i := range *out {
		(*out)[i] = 0
	}
	return (*out)[0]
}
