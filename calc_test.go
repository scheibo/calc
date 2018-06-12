package calc

import (
	"math"
	"testing"
)

const (
	// Epsilon is some tiny value that determines how precisely equal we want our
	// floats to be.
	Epsilon float64 = 1e-3
	// MinNormal is the smallest normal value possible.
	MinNormal = float64(2.2250738585072014E-308) // 1 / 2**(1022)
)

func fequal(a, b float64) bool {
	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)

	if a == b {
		return true
	} else if a == b || b == 0 || diff < MinNormal {
		// a or b is zero or both are extremely close to it
		// relative error is less meaningful here
		return diff < (Epsilon * MinNormal)
	} else {
		// use relative error
		return diff/(absA+absB) < Epsilon
	}
}

func TestPowerRR(t *testing.T) {
	tests := []struct {
		vg, gr, crr, mt, g, expected float64
	}{
		{8.36, 0.003, 0.0032, 90, G, 23.6},
	}
	for _, tt := range tests {
		actual := PowerRR(tt.vg, tt.gr, tt.crr, tt.mt, tt.g)
		if !fequal(actual, tt.expected) {
			t.Errorf("PowerRR(%.3f, %.3f, %.3f, %.3f, %.3f): got: %.3f, want: %.3f",
				tt.vg, tt.gr, tt.crr, tt.mt, tt.g, actual, tt.expected)
		}
	}
}

func TestAirPressure(t *testing.T) {
	tests := []struct {
		h, t, expected float64
	}{
		{4000, 30, 64557.76},
	}
	for _, tt := range tests {
		actual := AirPressure(tt.h, tt.t)
		if !fequal(actual, tt.expected) {
			t.Errorf("AirPressure(%.3f, %.3f): got: %.3f, want: %.3f",
				tt.h, tt.t, actual, tt.expected)
		}
	}
}

func TestAirDensity(t *testing.T) {
	tests := []struct {
		h, g, expected float64
	}{
		{0, G, SeaLevelRho},
		{1000, G, 1.111},
	}
	for _, tt := range tests {
		actual := AirDensity(tt.h, tt.g)
		if !fequal(actual, tt.expected) {
			t.Errorf("AirDensity(%.3f, %.3f): got: %.3f, want: %.3f",
				tt.h, tt.g, actual, tt.expected)
		}
	}
}
