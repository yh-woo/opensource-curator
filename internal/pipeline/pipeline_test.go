package pipeline

import (
	"testing"
)

func TestFloatToNumeric(t *testing.T) {
	tests := []struct {
		input float64
	}{
		{0},
		{42.5},
		{100},
		{73.25},
	}
	for _, tt := range tests {
		n := floatToNumeric(tt.input)
		if !n.Valid {
			t.Errorf("floatToNumeric(%f) produced invalid Numeric", tt.input)
		}
		f, err := n.Float64Value()
		if err != nil {
			t.Errorf("Float64Value() error: %v", err)
		}
		if !f.Valid {
			t.Errorf("Float64Value() not valid for input %f", tt.input)
		}
	}
}
