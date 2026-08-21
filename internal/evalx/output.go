package evalx

import "fmt"

type OutputLine struct {
	Label string
	Value string
}

func (l OutputLine) Render() string {
	return fmt.Sprintf("%-10s: %s", l.Label, l.Value)
}

func Line(label, value string) OutputLine {
	return OutputLine{Label: label, Value: value}
}

func (r EvalResult) Lines() []OutputLine {
	lines := []OutputLine{
		Line("case", r.Name),
		Line("u", fmt.Sprintf("%s m/s", formatVelocity(r.Velocity))),
		Line("H", fmt.Sprintf("%s m", formatSci(r.Height))),
		Line("N", fmt.Sprintf("%s plates", formatPlates(r.Plates))),
	}
	if r.HasRequirement {
		lines = append(lines,
			Line("H_max", fmt.Sprintf("%s m", formatSci(r.MaxAllowedHeight))),
			Line("meets", fmt.Sprintf("%v", r.MeetsRequirement)),
		)
	}
	if r.HasRetention {
		lines = append(lines, Line("Rs", fmt.Sprintf("%s", formatNumber(r.Resolution))))
	}
	return lines
}

func formatVelocity(v float64) string {
	return fmt.Sprintf("%f", v)
}

func formatSci(v float64) string {
	return fmt.Sprintf("%e", v)
}

func formatPlates(v float64) string {
	return fmt.Sprintf("%.1f", v)
}

func formatNumber(v float64) string {
	return fmt.Sprintf("%g", v)
}
