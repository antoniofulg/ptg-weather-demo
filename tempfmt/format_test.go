package tempfmt

import "testing"

func TestFormat(t *testing.T) {
	tests := []struct {
		name string
		temp float64
		unit Unit
		want string
	}{
		{"UT-020 celsius", 21.34, Celsius, "21.3°C"},
		{"UT-021 fahrenheit conversion and round", 21.34, Fahrenheit, "70.4°F"},
		{"UT-022 unknown unit falls back to celsius", 21.34, Unit("k"), "21.3°C"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.temp, tt.unit); got != tt.want {
				t.Errorf("Format(%v, %q) = %q, want %q", tt.temp, tt.unit, got, tt.want)
			}
		})
	}
}
