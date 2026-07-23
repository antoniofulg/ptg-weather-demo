// Package tempfmt renders a Celsius temperature in a requested unit.
package tempfmt

import "fmt"

// Unit selects the output temperature scale.
type Unit string

const (
	// Celsius renders the temperature in degrees Celsius.
	Celsius Unit = "c"
	// Fahrenheit renders the temperature in degrees Fahrenheit.
	Fahrenheit Unit = "f"
)

// Format renders a Celsius temperature in the requested unit to one decimal
// place with the correct symbol, e.g. Format(21.34, Celsius) == "21.3°C" and
// Format(21.34, Fahrenheit) == "70.4°F" (using f = c*9/5 + 32). An unknown unit
// falls back to Celsius.
func Format(tempC float64, unit Unit) string {
	if unit == Fahrenheit {
		return fmt.Sprintf("%.1f°F", tempC*9/5+32)
	}
	return fmt.Sprintf("%.1f°C", tempC)
}
