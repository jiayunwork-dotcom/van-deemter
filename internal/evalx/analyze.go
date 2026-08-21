package evalx

import (
	"strings"

	"van-deemter/internal/hmodel"
)

type Analysis struct {
	Eval     EvalResult
	Scan     ScanResult
	Budget   hmodel.Budget
	Range    RangeInfo
	Inspect  InspectInfo
}

type RangeInfo struct {
	Feasible bool
	Lower    float64
	Upper    float64
}

type InspectInfo struct {
	DominantTerm string
	TermRatio    float64
	Efficiency   float64
}

func Analyze(c Case, u float64) (Analysis, error) {
	if err := hmodel.ValidateVelocity(u); err != nil {
		return Analysis{}, err
	}
	eval, err := Eval(c, u)
	if err != nil {
		return Analysis{}, err
	}
	scan, err := Scan(c)
	if err != nil {
		return Analysis{}, err
	}
	m, col, err := c.Resolve()
	if err != nil {
		return Analysis{}, err
	}
	budget, err := m.HeightBudget(u)
	if err != nil {
		return Analysis{}, err
	}
	analysis := Analysis{Eval: eval, Scan: scan, Budget: budget}
	if c.HasRequirement() {
		req, err := c.Requirement()
		if err != nil {
			return Analysis{}, err
		}
		rng, err := col.VelocityRangeForRequirement(req)
		if err != nil {
			return Analysis{}, err
		}
		analysis.Range = RangeInfo{Feasible: rng.Feasible, Lower: rng.Lower, Upper: rng.Upper}
	}
	term, err := m.DominantTerm(u)
	if err != nil {
		return Analysis{}, err
	}
	ratio, err := m.TermRatio(u)
	if err != nil {
		return Analysis{}, err
	}
	nOpt, _, err := col.PlatesAtOptimum()
	if err != nil {
		return Analysis{}, err
	}
	efficiency := 0.0
	if nOpt > 0 {
		efficiency = eval.Plates / nOpt
	}
	analysis.Inspect = InspectInfo{DominantTerm: term, TermRatio: ratio, Efficiency: efficiency}
	return analysis, nil
}

func RenderAnalysis(a Analysis) string {
	var b strings.Builder
	b.WriteString(RenderEval(a.Eval))
	b.WriteString("\n")
	b.WriteString(RenderScan(a.Scan))
	b.WriteString("\n")
	b.WriteString("eddy    : " + hmodel.FormatScientific(a.Budget.Eddy) + " m\n")
	b.WriteString("axial   : " + hmodel.FormatScientific(a.Budget.Axial) + " m\n")
	b.WriteString("resist  : " + hmodel.FormatScientific(a.Budget.Resist) + " m\n")
	b.WriteString("dominant: " + a.Inspect.DominantTerm + "\n")
	b.WriteString("eff     : " + hmodel.FormatPercent(a.Inspect.Efficiency) + "\n")
	if a.Range.Feasible {
		b.WriteString("u_range : " + hmodel.FormatVelocity(a.Range.Lower) + " .. " + hmodel.FormatVelocity(a.Range.Upper) + " m/s\n")
	} else {
		b.WriteString("u_range : none (requirement infeasible)\n")
	}
	return strings.TrimRight(b.String(), "\n")
}
