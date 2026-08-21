package hmodel

import (
	"math"
	"strings"
	"testing"
)

const (
	testA = 2.4e-5
	testB = 1.6e-8
	testC = 4.8e-4
)

func baseModel() Model {
	m, err := NewModel(testA, testB, testC)
	if err != nil {
		panic(err)
	}
	return m
}

func relDiff(got, want float64) float64 {
	if want == 0 {
		return math.Abs(got)
	}
	return math.Abs(got-want) / math.Abs(want)
}

func TestRejectNonPositiveVelocity(t *testing.T) {
	m := baseModel()
	for _, u := range []float64{0, -0.02, -1e-9} {
		if _, err := m.Height(u); err == nil {
			t.Errorf("Height(u=%g) expected error, got nil", u)
		}
		if _, err := m.CentralDerivative(u, 1e-4); err == nil {
			t.Errorf("CentralDerivative(u=%g) expected error, got nil", u)
		}
	}
	if _, err := m.Height(0); err == nil || !strings.Contains(err.Error(), "velocity") {
		t.Errorf("Height(0) error %q does not name the velocity boundary", err)
	}
}

func TestRejectNonPositiveCoefficient(t *testing.T) {
	cases := []struct {
		label string
		a, b, c float64
	}{
		{"A negative", -1, testB, testC},
		{"B negative", testA, -1, testC},
		{"C negative", testA, testB, -1},
		{"A zero", 0, testB, testC},
		{"B zero", testA, 0, testC},
		{"C zero", testA, testB, 0},
	}
	for _, tc := range cases {
		if _, err := NewModel(tc.a, tc.b, tc.c); err == nil {
			t.Errorf("%s: expected error for (%g,%g,%g)", tc.label, tc.a, tc.b, tc.c)
		}
	}
}

func TestOptimumClosedForm(t *testing.T) {
	m := baseModel()
	got, err := m.OptimumVelocity()
	if err != nil {
		t.Fatalf("OptimumVelocity error: %v", err)
	}
	want := math.Sqrt(testB / testC)
	if relDiff(got, want) > 1e-12 {
		t.Errorf("u_opt = %g, want %g (sqrt(B/C))", got, want)
	}
}

func TestMinimumClosedForm(t *testing.T) {
	m := baseModel()
	got, err := m.MinimumHeight()
	if err != nil {
		t.Fatalf("MinimumHeight error: %v", err)
	}
	want := testA + 2*math.Sqrt(testB*testC)
	if relDiff(got, want) > 1e-12 {
		t.Errorf("H_min = %g, want %g (A+2*sqrt(BC))", got, want)
	}
	uOpt, _ := m.OptimumVelocity()
	hAtOpt, err := m.Height(uOpt)
	if err != nil {
		t.Fatalf("Height at u_opt error: %v", err)
	}
	if relDiff(hAtOpt, want) > 1e-9 {
		t.Errorf("H(u_opt) = %g, want H_min = %g", hAtOpt, want)
	}
}

func TestDerivativeAtOptimum(t *testing.T) {
	m := baseModel()
	uOpt, _ := m.OptimumVelocity()
	scale := testB/(uOpt*uOpt) + testC
	tol := 1e-6 * scale
	for _, delta := range []float64{0.001 * uOpt, 0.0005 * uOpt, 0.0002 * uOpt} {
		d, err := m.CentralDerivative(uOpt, delta)
		if err != nil {
			t.Fatalf("CentralDerivative(delta=%g) error: %v", delta, err)
		}
		if math.Abs(d) > tol {
			t.Errorf("central derivative at u_opt = %g, want |.| <= %g (delta=%g)", d, tol, delta)
		}
	}
}

func TestOptimumRisesWithB(t *testing.T) {
	m := baseModel()
	boosted := m.WithB(testB * 2)
	ok, err := OptimumRisesWithB(m, boosted)
	if err != nil {
		t.Fatalf("OptimumRisesWithB error: %v", err)
	}
	if !ok {
		uBase, _ := m.OptimumVelocity()
		uBoost, _ := boosted.OptimumVelocity()
		t.Errorf("u_opt should rise when B alone increases: base=%g boosted=%g", uBase, uBoost)
	}
}

func TestOptimumFallsWithC(t *testing.T) {
	m := baseModel()
	boosted := m.WithC(testC * 2)
	ok, err := OptimumFallsWithC(m, boosted)
	if err != nil {
		t.Fatalf("OptimumFallsWithC error: %v", err)
	}
	if !ok {
		uBase, _ := m.OptimumVelocity()
		uBoost, _ := boosted.OptimumVelocity()
		t.Errorf("u_opt should fall when C alone increases: base=%g boosted=%g", uBase, uBoost)
	}
}

func TestMinimumRisesWithC(t *testing.T) {
	m := baseModel()
	boosted := m.WithC(testC * 2)
	ok, err := MinimumRisesWithC(m, boosted)
	if err != nil {
		t.Fatalf("MinimumRisesWithC error: %v", err)
	}
	if !ok {
		hBase, _ := m.MinimumHeight()
		hBoost, _ := boosted.MinimumHeight()
		t.Errorf("H_min should rise when C alone increases: base=%g boosted=%g", hBase, hBoost)
	}
}

func TestSidePointsAboveMinimum(t *testing.T) {
	m := baseModel()
	hMin, _ := m.MinimumHeight()
	for _, span := range []float64{0.2, 0.5, 0.8} {
		res, err := CheckSidePoints(m, span)
		if err != nil {
			t.Fatalf("CheckSidePoints(span=%g) error: %v", span, err)
		}
		if !res.AllOK {
			t.Errorf("span=%g: H below u_opt = %g (H_min=%g, gap=%g), above = %g (gap=%g)",
				span, res.Lower.H, hMin, res.Lower.Gap, res.Upper.H, res.Upper.Gap)
		}
		lowerRel := relDiff(res.Lower.H, hMin)
		upperRel := relDiff(res.Upper.H, hMin)
		if lowerRel < 1e-6 || upperRel < 1e-6 {
			t.Errorf("span=%g: side points must be strictly above H_min (lower rel=%g, upper rel=%g)",
				span, lowerRel, upperRel)
		}
	}
}

func TestGridScanMatchesClosedForm(t *testing.T) {
	m := baseModel()
	res, err := Scan(m)
	if err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	uOpt, _ := m.OptimumVelocity()
	hMin, _ := m.MinimumHeight()
	if relDiff(res.OptimumVelocity, uOpt) > 1e-12 {
		t.Errorf("scan u_opt = %g, want %g", res.OptimumVelocity, uOpt)
	}
	if relDiff(res.MinimumHeight, hMin) > 1e-12 {
		t.Errorf("scan H_min = %g, want %g", res.MinimumHeight, hMin)
	}
	if !res.GridFound {
		t.Fatal("grid scan found no minimum")
	}
	if relDiff(res.GridMinimum.H, hMin) > 1e-3 {
		t.Errorf("grid minimum H = %g, closed-form H_min = %g", res.GridMinimum.H, hMin)
	}
	if !res.DerivativeWithin {
		t.Errorf("central derivative at u_opt = %g, expected within tolerance", res.DerivativeAtOpt)
	}
}

func TestInverseVelocitiesEncloseTarget(t *testing.T) {
	m := baseModel()
	uOpt, _ := m.OptimumVelocity()
	hMin, _ := m.MinimumHeight()
	for _, mult := range []float64{1.1, 1.5, 2.0, 4.0} {
		target := hMin * mult
		lo, hi, err := m.InverseVelocities(target)
		if err != nil {
			t.Fatalf("InverseVelocities(target=%.6g) error: %v", target, err)
		}
		if !(lo < uOpt && uOpt < hi) {
			t.Errorf("target=%.6g: optimum %g not inside [%g,%g]", target, uOpt, lo, hi)
		}
		inside, _ := m.Height(0.5 * (lo + hi))
		outsideLow, _ := m.Height(0.5 * lo)
		outsideHigh, _ := m.Height(2 * hi)
		if inside > target || outsideLow < target || outsideHigh < target {
			t.Errorf("target=%.6g: H inside interval %g, outside %g/%g", target, inside, outsideLow, outsideHigh)
		}
	}
	if _, _, err := m.InverseVelocities(hMin * 0.5); err == nil {
		t.Error("target below H_min must error")
	}
}

func TestInverseUniqueAtMinimum(t *testing.T) {
	m := baseModel()
	uOpt, _ := m.OptimumVelocity()
	hMin, _ := m.MinimumHeight()
	lo, hi, err := m.InverseVelocities(hMin)
	if err != nil {
		t.Fatalf("InverseVelocities(H_min) error: %v", err)
	}
	if relDiff(lo, uOpt) > 1e-6 || relDiff(hi, uOpt) > 1e-6 {
		t.Errorf("roots at H_min (%g,%g) must both equal u_opt %g", lo, hi, uOpt)
	}
}
