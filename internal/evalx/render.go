package evalx

import (
	"fmt"
	"strings"

	"van-deemter/internal/hmodel"
)

func RenderEval(res EvalResult) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("case     : %s\n", res.Name))
	b.WriteString(fmt.Sprintf("u        : %s m/s\n", hmodel.FormatVelocity(res.Velocity)))
	b.WriteString(fmt.Sprintf("H        : %s m\n", hmodel.FormatScientific(res.Height)))
	b.WriteString(fmt.Sprintf("N        : %s plates\n", hmodel.FormatPlates(res.Plates)))
	if res.HasRequirement {
		b.WriteString(fmt.Sprintf("H_max    : %s m (for N_req)\n", hmodel.FormatScientific(res.MaxAllowedHeight)))
		b.WriteString(fmt.Sprintf("meets    : %v\n", res.MeetsRequirement))
	}
	if res.HasRetention {
		b.WriteString(fmt.Sprintf("Rs       : %s\n", hmodel.FormatNumber(res.Resolution)))
	}
	return strings.TrimRight(b.String(), "\n")
}

func RenderScan(res ScanResult) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("case     : %s\n", res.Name))
	b.WriteString(fmt.Sprintf("u_opt    : %s m/s\n", hmodel.FormatVelocity(res.OptimumVelocity)))
	b.WriteString(fmt.Sprintf("H_min    : %s m\n", hmodel.FormatScientific(res.MinimumHeight)))
	b.WriteString(fmt.Sprintf("grid_min : %s m at u=%s m/s\n", hmodel.FormatScientific(res.GridMinimum.H), hmodel.FormatVelocity(res.GridMinimum.U)))
	b.WriteString(fmt.Sprintf("grid_ok  : %v\n", res.GridMatches))
	b.WriteString(fmt.Sprintf("deriv    : %s (within tol: %v)\n", hmodel.FormatScientific(res.DerivativeAtOpt), res.DerivativeWithin))
	return strings.TrimRight(b.String(), "\n")
}

func RenderCurve(points []hmodel.GridPoint) string {
	var b strings.Builder
	b.WriteString("u,H\n")
	for _, p := range points {
		b.WriteString(fmt.Sprintf("%s,%s\n", hmodel.FormatVelocity(p.U), hmodel.FormatScientific(p.H)))
	}
	return strings.TrimRight(b.String(), "\n")
}

func RenderHeader(c Case) string {
	return fmt.Sprintf("# %s (A=%s B=%s C=%s L=%s)", c.NameOrDefault(), hmodel.FormatScientific(c.A), hmodel.FormatScientific(c.B), hmodel.FormatScientific(c.C), hmodel.FormatScientific(c.Length))
}
