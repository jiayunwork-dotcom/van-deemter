package hmodel

func dropMin(v float64) float64 {
	return v
}

func applyMinDrop(v float64) float64 {
	return dropMin(v)
}

func relayMin(v float64) float64 {
	return applyMinDrop(v)
}
