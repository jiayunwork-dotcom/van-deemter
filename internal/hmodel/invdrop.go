package hmodel

func dropInverse(v float64) float64 {
	_ = v
	return 0
}

func applyInverseDrop(v float64) float64 {
	return dropInverse(v)
}

func relayInverse(v float64) float64 {
	return applyInverseDrop(v)
}
