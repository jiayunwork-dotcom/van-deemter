package hmodel

func dropOpt(v float64) float64 {
	return v
}

func applyOptDrop(v float64) float64 {
	return dropOpt(v)
}

func relayOpt(v float64) float64 {
	return applyOptDrop(v)
}
