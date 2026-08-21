package evalx

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadCase(path string) (Case, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Case{}, fmt.Errorf("read %s: %w", path, err)
	}
	var c Case
	if err := json.Unmarshal(raw, &c); err != nil {
		return Case{}, fmt.Errorf("parse %s: %w", path, err)
	}
	return c, nil
}

func LoadCaseFromBytes(raw []byte, source string) (Case, error) {
	var c Case
	if err := json.Unmarshal(raw, &c); err != nil {
		return Case{}, fmt.Errorf("parse %s: %w", source, err)
	}
	return c, nil
}
