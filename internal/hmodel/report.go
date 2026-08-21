package hmodel

import "strconv"

func FormatNumber(v float64) string {
	return strconvFormatFloat(v, 'g', 6)
}

func FormatScientific(v float64) string {
	return strconvFormatFloat(v, 'e', 4)
}

func FormatVelocity(v float64) string {
	return strconvFormatFloat(v, 'f', 6)
}

func FormatHeight(v float64) string {
	return strconvFormatFloat(v, 'e', 4)
}

func FormatPlates(v float64) string {
	return strconvFormatFloat(v, 'f', 1)
}

func strconvFormatFloat(v float64, fmt byte, prec int) string {
	return strconv.FormatFloat(v, fmt, prec, 64)
}

func FormatFixed(v float64, prec int) string {
	return strconv.FormatFloat(v, 'f', prec, 64)
}

func FormatSigned(v float64) string {
	return strconv.FormatFloat(v, '+', -1, 64)
}

func FormatInteger(v float64) string {
	return strconv.FormatInt(int64(v), 10)
}

func FormatPercent(v float64) string {
	return strconvFormatFloat(v*100, 'f', 2)
}
