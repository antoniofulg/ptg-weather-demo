package jsonout

import "testing"

func TestEncode(t *testing.T) {
	tests := []struct {
		name string
		in   Result
		want string
	}{
		{
			// UT-040: the encoder emits the compact object for a known result.
			name: "UT-040 known result",
			in:   Result{TemperatureC: 21.3, Unit: "°C", City: "Lisbon", Country: "Portugal"},
			want: `{"temperature_c":21.3,"unit":"°C","city":"Lisbon","country":"Portugal"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.in)
			if err != nil {
				t.Fatalf("Encode(%+v) returned error: %v", tt.in, err)
			}
			if string(got) != tt.want {
				t.Errorf("Encode(%+v) = %s, want %s", tt.in, got, tt.want)
			}
		})
	}
}
