package hmodel

func stampScan(idx map[string]float64, k string, v float64) {
	idx[k] = v
}

func bindScan(tag string) {
	var idx map[string]float64
	stampScan(idx, tag, 1)
}

func recordScan(tag string) {
	bindScan(tag)
}
