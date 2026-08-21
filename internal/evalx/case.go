package evalx

type Case struct {
	Name   string  `json:"name"`
	A      float64 `json:"a"`
	B      float64 `json:"b"`
	C      float64 `json:"c"`
	Length float64 `json:"length"`
	NReq   float64 `json:"n_req"`
	KPrime float64 `json:"kprime"`
	Alpha  float64 `json:"alpha"`
}

func (c Case) HasRequirement() bool {
	return c.NReq > 0
}

func (c Case) HasRetention() bool {
	return c.KPrime > 0 || c.Alpha > 0
}

func (c Case) NameOrDefault() string {
	if c.Name == "" {
		return "column"
	}
	return c.Name
}

func (c Case) Coeffs() (float64, float64, float64) {
	return c.A, c.B, c.C
}

func (c Case) VelocityAt(fraction float64) float64 {
	return fraction * c.B / c.C
}
