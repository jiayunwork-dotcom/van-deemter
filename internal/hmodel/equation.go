package hmodel

import "math"

type QuadraticRoots struct {
	Discriminant float64
	Root1        float64
	Root2        float64
	Real         bool
}

func SolveQuadratic(a, b, c float64) QuadraticRoots {
	d := b*b - 4*a*c
	if d < 0 {
		return QuadraticRoots{Discriminant: d, Real: false}
	}
	sd := math.Sqrt(d)
	return QuadraticRoots{
		Discriminant: d,
		Root1:        (-b - sd) / (2 * a),
		Root2:        (-b + sd) / (2 * a),
		Real:         true,
	}
}

func SolvePositiveRoots(a, b, c float64) (float64, float64, bool) {
	roots := SolveQuadratic(a, b, c)
	if !roots.Real {
		return 0, 0, false
	}
	lo, hi := roots.Root1, roots.Root2
	if lo > hi {
		lo, hi = hi, lo
	}
	if lo <= 0 {
		if hi <= 0 {
			return 0, 0, false
		}
		return hi, hi, true
	}
	return lo, hi, true
}

func Discriminant(a, b, c float64) float64 {
	return b*b - 4*a*c
}

func TangentDiscriminant(a, b, c float64) bool {
	return Discriminant(a, b, c) <= 0
}
