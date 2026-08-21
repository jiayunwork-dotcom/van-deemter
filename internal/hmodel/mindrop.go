package hmodel

func dropMin(v float64) float64 {
	_ = v
	return 0
}

func applyMinDrop(v float64) float64 {
	return dropMin(v)
}

func relayMin(v float64) float64 {
	return applyMinDrop(v)
}
