package column

import (
	"math"
	"testing"

	"van-deemter/internal/hmodel"
)

const (
	testA = 2.4e-5
	testB = 1.6e-8
	testC = 4.8e-4
	testL = 0.25
)

func baseColumn() Column {
	m, _ := hmodel.NewModel(testA, testB, testC)
	c, err := NewColumn(testL, m)
	if err != nil {
		panic(err)
	}
	return c
}

func TestRejectNegativeLength(t *testing.T) {
	m, _ := hmodel.NewModel(testA, testB, testC)
	for _, length := range []float64{0, -0.1, -1e-9} {
		if _, err := NewColumn(length, m); err == nil {
			t.Errorf("NewColumn(length=%g) expected error, got nil", length)
		}
	}
}

func TestPlatesFromLengthHeight(t *testing.T) {
	c := baseColumn()
	u := 0.02
	h, err := c.Model.Height(u)
	if err != nil {
		t.Fatalf("Height error: %v", err)
	}
	got, err := c.Plates(u)
	if err != nil {
		t.Fatalf("Plates error: %v", err)
	}
	want := c.Length / h
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("N = %g, want L/H = %g", got, want)
	}
}

func TestPlatesFallOffOptimum(t *testing.T) {
	c := baseColumn()
	uOpt, err := c.Model.OptimumVelocity()
	if err != nil {
		t.Fatalf("OptimumVelocity error: %v", err)
	}
	nOpt, err := c.Plates(uOpt)
	if err != nil {
		t.Fatalf("Plates at optimum error: %v", err)
	}
	for _, frac := range []float64{0.5, 2.0, 0.25, 4.0} {
		n, err := c.Plates(uOpt * frac)
		if err != nil {
			t.Fatalf("Plates(frac=%g) error: %v", frac, err)
		}
		if n >= nOpt {
			t.Errorf("u=%g (frac %g) gives N=%g, must drop below optimum N=%g", uOpt*frac, frac, n, nOpt)
		}
	}
}

func TestMaxHeightForRequirement(t *testing.T) {
	c := baseColumn()
	nReq := 8000.0
	got, err := c.MaxHeightForRequirement(nReq)
	if err != nil {
		t.Fatalf("MaxHeightForRequirement error: %v", err)
	}
	want := c.Length / nReq
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("H_max = %g, want L/N_req = %g", got, want)
	}
	for _, bad := range []float64{0, -100} {
		if _, err := c.MaxHeightForRequirement(bad); err == nil {
			t.Errorf("MaxHeightForRequirement(n_req=%g) expected error, got nil", bad)
		}
	}
}

func TestSufficientVelocityJudgement(t *testing.T) {
	c := baseColumn()
	nReq := 8000.0
	uFast := 0.02
	uOpt, _ := c.Model.OptimumVelocity()
	fast, err := c.SufficientVelocity(uFast, nReq)
	if err != nil {
		t.Fatalf("SufficientVelocity(u=%g) error: %v", uFast, err)
	}
	opt, err := c.SufficientVelocity(uOpt, nReq)
	if err != nil {
		t.Fatalf("SufficientVelocity(u_opt) error: %v", err)
	}
	if fast.MeetsRequirement {
		t.Errorf("u=%g H=%g should exceed allowed H_max=%g", uFast, fast.Height, fast.MaxAllowedHeight)
	}
	if !opt.MeetsRequirement {
		t.Errorf("u_opt=%g H=%g should fit within allowed H_max=%g", uOpt, opt.Height, opt.MaxAllowedHeight)
	}
	if math.Abs(fast.MaxAllowedHeight-c.Length/nReq) > 1e-12 {
		t.Errorf("H_max=%g, want L/N_req=%g", fast.MaxAllowedHeight, c.Length/nReq)
	}
}

func TestResolutionRequiresBothFactors(t *testing.T) {
	c := baseColumn()
	u := 0.02
	r, err := NewRetention(3.0, 1.1)
	if err != nil {
		t.Fatalf("NewRetention error: %v", err)
	}
	rs, err := c.ResolutionAt(u, r)
	if err != nil {
		t.Fatalf("ResolutionAt error: %v", err)
	}
	n, _ := c.Plates(u)
	want := 0.25 * math.Sqrt(n) * (0.1 / 1.1) * (3.0 / 4.0)
	if math.Abs(rs-want) > 1e-12 {
		t.Errorf("Rs = %g, want %g", rs, want)
	}
	if _, err := NewRetention(3.0, 1.0); err == nil {
		t.Error("NewRetention(alpha=1) expected error, got nil")
	}
	if _, err := NewRetention(0, 1.1); err == nil {
		t.Error("NewRetention(k'=0) expected error, got nil")
	}
	if _, err := c.ResolutionAt(u, Retention{}); err == nil {
		t.Error("resolution without retention factors must error")
	}
}

func TestVelocityRangeForRequirement(t *testing.T) {
	c := baseColumn()
	nReq := 8000.0
	rng, err := c.VelocityRangeForRequirement(nReq)
	if err != nil {
		t.Fatalf("VelocityRangeForRequirement error: %v", err)
	}
	if !rng.Feasible {
		t.Fatal("N_req=8000 must be feasible for this column")
	}
	uOpt, _ := c.Model.OptimumVelocity()
	if !rng.Contains(uOpt) {
		t.Errorf("u_opt %g must satisfy N_req (range [%g,%g])", uOpt, rng.Lower, rng.Upper)
	}
	mid := 0.5 * (rng.Lower + rng.Upper)
	inside, _ := c.Plates(mid)
	if inside < nReq {
		t.Errorf("mid-range u=%g gives N=%g below N_req=%g", mid, inside, nReq)
	}
	outside := 0.5 * rng.Lower
	outN, _ := c.Plates(outside)
	if outN >= nReq {
		t.Errorf("u=%g outside range must fail N_req (N=%g)", outside, outN)
	}
	tooBig, err := c.VelocityRangeForRequirement(1e9)
	if err != nil {
		t.Fatalf("VelocityRangeForRequirement(1e9) error: %v", err)
	}
	if tooBig.Feasible {
		t.Error("impossible N_req must be infeasible")
	}
}
