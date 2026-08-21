package evalx

import (
	"math"
	"strings"
	"testing"
)

const (
	testA = 2.4e-5
	testB = 1.6e-8
	testC = 4.8e-4
	testL = 0.25
)

func packedCase() Case {
	return Case{
		Name:   "packed-column",
		A:      testA,
		B:      testB,
		C:      testC,
		Length: testL,
		NReq:   8000,
		KPrime: 3.0,
		Alpha:  1.1,
	}
}

func relDiff(got, want float64) float64 {
	if want == 0 {
		return math.Abs(got)
	}
	return math.Abs(got-want) / math.Abs(want)
}

func TestEvalConsistencyFromExample(t *testing.T) {
	c, err := LoadCase("../../example/packed.json")
	if err != nil {
		t.Fatalf("LoadCase error: %v", err)
	}
	if c.A != testA || c.B != testB || c.C != testC || c.Length != testL {
		t.Errorf("case coeffs (%g,%g,%g,L=%g) differ from packed.json (%g,%g,%g,L=%g)",
			c.A, c.B, c.C, c.Length, testA, testB, testC, testL)
	}
	m, col, err := c.Resolve()
	if err != nil {
		t.Fatalf("Resolve error: %v", err)
	}
	uOpt, err := m.OptimumVelocity()
	if err != nil {
		t.Fatalf("OptimumVelocity error: %v", err)
	}
	wantOpt := math.Sqrt(c.B / c.C)
	if relDiff(uOpt, wantOpt) > 1e-12 {
		t.Errorf("u_opt = %g, want %g", uOpt, wantOpt)
	}
	hMin, err := m.MinimumHeight()
	if err != nil {
		t.Fatalf("MinimumHeight error: %v", err)
	}
	wantMin := c.A + 2*math.Sqrt(c.B*c.C)
	if relDiff(hMin, wantMin) > 1e-12 {
		t.Errorf("H_min = %g, want %g", hMin, wantMin)
	}
	u := 0.02
	res, err := Eval(c, u)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	wantH := c.A + c.B/u + c.C*u
	if relDiff(res.Height, wantH) > 1e-9 {
		t.Errorf("H = %g, want %g", res.Height, wantH)
	}
	if relDiff(res.Plates, c.Length/wantH) > 1e-9 {
		t.Errorf("N = %g, want %g", res.Plates, c.Length/wantH)
	}
	hAtOpt, err := m.Height(uOpt)
	if err != nil {
		t.Fatalf("Height at optimum error: %v", err)
	}
	if relDiff(hAtOpt, hMin) > 1e-9 {
		t.Errorf("H(u_opt)=%g must equal H_min=%g", hAtOpt, hMin)
	}
	nOpt, err := col.Plates(uOpt)
	if err != nil {
		t.Fatalf("Plates at optimum error: %v", err)
	}
	if !(nOpt > res.Plates) {
		t.Errorf("N(u_opt)=%g must exceed N(u)=%g", nOpt, res.Plates)
	}
	scale := c.B/(uOpt*uOpt) + c.C
	delta := 0.001 * uOpt
	deriv, err := m.CentralDerivative(uOpt, delta)
	if err != nil {
		t.Fatalf("CentralDerivative error: %v", err)
	}
	if math.Abs(deriv) > 1e-6*scale {
		t.Errorf("central derivative at u_opt = %g, want |.| <= %g", deriv, 1e-6*scale)
	}
}

func TestScanFromExample(t *testing.T) {
	c, err := LoadCase("../../example/packed.json")
	if err != nil {
		t.Fatalf("LoadCase error: %v", err)
	}
	res, err := Scan(c)
	if err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	wantOpt := math.Sqrt(c.B / c.C)
	if relDiff(res.OptimumVelocity, wantOpt) > 1e-12 {
		t.Errorf("scan u_opt = %g, want %g", res.OptimumVelocity, wantOpt)
	}
	wantMin := c.A + 2*math.Sqrt(c.B*c.C)
	if relDiff(res.MinimumHeight, wantMin) > 1e-12 {
		t.Errorf("scan H_min = %g, want %g", res.MinimumHeight, wantMin)
	}
	if !res.GridMatches {
		t.Errorf("grid minimum %g must match closed-form H_min %g", res.GridMinimum.H, res.MinimumHeight)
	}
	if !res.DerivativeWithin {
		t.Errorf("derivative at u_opt = %g, expected within tolerance", res.DerivativeAtOpt)
	}
}

func TestEvalReportsRequirement(t *testing.T) {
	c := packedCase()
	res, err := Eval(c, 0.02)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	wantMax := c.Length / c.NReq
	if relDiff(res.MaxAllowedHeight, wantMax) > 1e-12 {
		t.Errorf("H_max = %g, want L/N_req = %g", res.MaxAllowedHeight, wantMax)
	}
	if res.MeetsRequirement {
		t.Errorf("u=0.02 must not meet N_req=%g (H=%g > H_max=%g)", c.NReq, res.Height, res.MaxAllowedHeight)
	}
	resOpt, err := Eval(c, math.Sqrt(c.B/c.C))
	if err != nil {
		t.Fatalf("Eval at optimum error: %v", err)
	}
	if !resOpt.MeetsRequirement {
		t.Errorf("u_opt must meet N_req=%g (H=%g <= H_max=%g)", c.NReq, resOpt.Height, resOpt.MaxAllowedHeight)
	}
}

func TestEvalOmitsRsWithoutRetention(t *testing.T) {
	c := packedCase()
	c.KPrime = 0
	c.Alpha = 0
	res, err := Eval(c, 0.02)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if res.HasRetention {
		t.Error("resolution must not be reported when kprime and alpha are absent")
	}
	if res.Resolution != 0 {
		t.Errorf("resolution = %g, want 0 when absent", res.Resolution)
	}
}

func TestEvalResolutionWithRetention(t *testing.T) {
	c := packedCase()
	res, err := Eval(c, 0.02)
	if err != nil {
		t.Fatalf("Eval error: %v", err)
	}
	if !res.HasRetention {
		t.Fatal("resolution expected when kprime and alpha are present")
	}
	n := res.Plates
	want := 0.25 * math.Sqrt(n) * (0.1 / 1.1) * (3.0 / 4.0)
	if relDiff(res.Resolution, want) > 1e-12 {
		t.Errorf("Rs = %g, want %g", res.Resolution, want)
	}
}

func TestRetentionPartialInput(t *testing.T) {
	c := packedCase()
	c.Alpha = 0
	if _, err := Eval(c, 0.02); err == nil {
		t.Error("kprime without alpha must error")
	}
	c2 := packedCase()
	c2.KPrime = 0
	if _, err := Eval(c2, 0.02); err == nil {
		t.Error("alpha without kprime must error")
	}
}

func TestMissingFileReportsError(t *testing.T) {
	if _, err := LoadCase("../../example/does-not-exist.json"); err == nil {
		t.Error("LoadCase on missing file expected error, got nil")
	}
}

func TestMalformedJSONReportsError(t *testing.T) {
	raw := []byte(`{"a": 2.4e-5, "b": not-a-number}`)
	if _, err := LoadCaseFromBytes(raw, "malformed"); err == nil {
		t.Error("malformed JSON expected error, got nil")
	}
}

func TestNegativeFieldsInJSON(t *testing.T) {
	c := packedCase()
	c.C = -4.8e-4
	if _, err := Eval(c, 0.02); err == nil {
		t.Error("negative C in case must error")
	}
	if _, _, err := c.Resolve(); err == nil {
		t.Error("negative C must fail Resolve")
	}
	c2 := packedCase()
	c2.Length = -0.1
	if _, err := Eval(c2, 0.02); err == nil {
		t.Error("negative length must error")
	}
	c3 := packedCase()
	if _, err := Eval(c3, 0); err == nil || !strings.Contains(err.Error(), "velocity") {
		t.Errorf("u=0 error %q must name the velocity boundary", err)
	}
}

func TestValidateCaseCollectsErrors(t *testing.T) {
	c := packedCase()
	if err := ValidateCase(c); err != nil {
		t.Fatalf("valid case must pass ValidateCase, got %v", err)
	}
	bad := Case{A: -1, B: 0, C: -2, Length: -3, NReq: -4, KPrime: 2, Alpha: 0.5}
	errs := ValidateCaseAll(bad)
	if len(errs) < 6 {
		t.Errorf("expected errors for every bad field, got %d: %v", len(errs), errs)
	}
	if err := ValidateCase(bad); err == nil {
		t.Error("invalid case must fail ValidateCase")
	}
}

func TestAnalyzeCombinesEvalAndRange(t *testing.T) {
	c := packedCase()
	analysis, err := Analyze(c, 0.02)
	if err != nil {
		t.Fatalf("Analyze error: %v", err)
	}
	if !analysis.Range.Feasible {
		t.Error("N_req=8000 must be feasible in analysis")
	}
	if analysis.Eval.Plates <= 0 || analysis.Scan.OptimumVelocity <= 0 {
		t.Errorf("analysis must carry eval/scan numbers (plates=%g, u_opt=%g)", analysis.Eval.Plates, analysis.Scan.OptimumVelocity)
	}
	if analysis.Inspect.DominantTerm == "" {
		t.Error("analysis must report the dominant H term")
	}
	if analysis.Budget.Total <= 0 {
		t.Error("analysis must carry an H budget")
	}
}
