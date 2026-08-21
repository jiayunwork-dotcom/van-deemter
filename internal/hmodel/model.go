package hmodel

import "fmt"

type Model struct {
	A float64
	B float64
	C float64
}

func NewModel(a, b, c float64) (Model, error) {
	m := Model{A: a, B: b, C: c}
	if err := m.Validate(); err != nil {
		return Model{}, err
	}
	return m, nil
}

func (m Model) WithA(a float64) Model {
	m.A = a
	return m
}

func (m Model) WithB(b float64) Model {
	m.B = b
	return m
}

func (m Model) WithC(c float64) Model {
	m.C = c
	return m
}

func (m Model) Fields() (float64, float64, float64) {
	return m.A, m.B, m.C
}

func (m Model) String() string {
	return fmt.Sprintf("A=%.6g B=%.6g C=%.6g", m.A, m.B, m.C)
}
