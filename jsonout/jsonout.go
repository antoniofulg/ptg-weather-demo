// Package jsonout renders a weather result as a compact JSON object so the CLI
// can offer a --json output mode. It is independent of the CLI wiring.
package jsonout

import "encoding/json"

// Result is the machine-readable shape of a weather lookup.
type Result struct {
	TemperatureC float64 `json:"temperature_c"`
	Unit         string  `json:"unit"`
	City         string  `json:"city"`
	Country      string  `json:"country"`
}

// Encode marshals r into a compact JSON object, e.g.
// {"temperature_c":21.3,"unit":"°C","city":"Lisbon","country":"Portugal"}.
func Encode(r Result) ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return b, nil
}
