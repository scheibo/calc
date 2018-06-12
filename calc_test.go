package calc

import (
	"math"
	"testing"
	"time"
)

const (
	// Epsilon is some tiny value that determines how precisely equal we want
	// our floats to be.
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

// TODO
// https://www.cyclingpowerlab.com/PowerModels.aspx
// http://anonymous.coward.free.fr/wattage/olh.html
// Power, PowerTOT

func TestPowerAT(t *testing.T) {
	tests := []struct {
		rho, cda, fw, va, vg, expected float64
	}{
		{SeaLevelRho, 0.2565, Fw, 10.91, 8.36, 158.8},
	}
	for _, tt := range tests {
		actual := PowerAT(tt.rho, tt.cda, tt.fw, tt.va, tt.vg)
		if !fequal(actual, tt.expected) {
			t.Errorf("PowerAT(%.3f, %.3f, %.3f, %.3f, %.3f): got: %.3f, want: %.3f",
				tt.rho, tt.cda, tt.fw, tt.va, tt.vg, actual, tt.expected)
		}
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

func TestPowerWB(t *testing.T) {
	tests := []struct {
		vg, expected float64
	}{
		{8.36, 1.37},
	}
	for _, tt := range tests {
		actual := PowerWB(tt.vg)
		if !fequal(actual, tt.expected) {
			t.Errorf("PowerWB(%.3f): got: %.3f, want: %.3f",
				tt.vg, actual, tt.expected)
		}
	}
}

func TestPowerPE(t *testing.T) {
	tests := []struct {
		vg, mt, g, gr, expected float64
	}{
		{8.36, 90, G, 0.003, 22.1},
	}
	for _, tt := range tests {
		actual := PowerPE(tt.vg, tt.mt, tt.g, tt.gr)
		if !fequal(actual, tt.expected) {
			t.Errorf("PowerPE(%.3f, %.3f, %.3f, %.3f): got: %.3f, want: %.3f",
				tt.vg, tt.mt, tt.g, tt.gr, actual, tt.expected)
		}
	}
}

func TestPowerKE(t *testing.T) {
	tests := []struct {
		mt, i, r, vgi, vgf, ti, tf, expected float64
	}{
		{90, I, 0.311, 8.28, 8.45, 0, 56.42, 2.305},
	}
	for _, tt := range tests {
		actual := PowerKE(tt.mt, tt.i, tt.r, tt.vgi, tt.vgf, tt.ti, tt.tf)
		if !fequal(actual, tt.expected) {
			t.Errorf("PowerKE(%.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f): got: %.3f, want: %.3f",
				tt.mt, tt.i, tt.r, tt.vgi, tt.vgf, tt.ti, tt.tf, actual, tt.expected)
		}
	}
}

func TestAirVelocity(t *testing.T) {
	tests := []struct {
		vg, vw, dw, db, expected float64
	}{
		{8.36, 2.94, 310, 340, 10.91},
	}
	for _, tt := range tests {
		actual := AirVelocity(tt.vg, tt.vw, tt.dw, tt.db)
		if !fequal(actual, tt.expected) {
			t.Errorf("AirVelocity(%.3f, %.3f, %.3f, %.3f): got: %.3f, want: %.3f",
				tt.vg, tt.vw, tt.dw, tt.db, actual, tt.expected)
		}
	}
}

func TestGroundVelocity(t *testing.T) {
	tests := []struct {
		d        float64
		t        time.Duration
		expected float64
	}{
		{4800, 18 * time.Minute, 4.444},
	}
	for _, tt := range tests {
		actual := GroundVelocity(tt.d, tt.t)
		if !fequal(actual, tt.expected) {
			t.Errorf("GroundVelocity(%.3f, %v): got: %.3f, want: %.3f",
				tt.d, tt.t, actual, tt.expected)
		}
	}
}

func TestYaw(t *testing.T) {
	tests := []struct {
		va, vw, dw, db, expected float64
	}{
		{10.91, 2.94, 310, 340, -7.67},
	}
	for _, tt := range tests {
		actual := Yaw(tt.va, tt.vw, tt.dw, tt.db)
		if !fequal(actual, tt.expected) {
			t.Errorf("Yaw(%.3f, %.3f, %.3f, %.3f): got: %.3f, want: %.3f",
				tt.va, tt.vw, tt.dw, tt.db, actual, tt.expected)
		}
	}
}

func TestCalculateDropsCdA(t *testing.T) {
	tests := []struct {
		h, m, expected float64
	}{
		{1.75, 69, 0.3653},
		{1.60, 53, 0.3295},
	}
	for _, tt := range tests {
		actual := CalculateDropsCdA(tt.h, tt.m)
		if !fequal(actual, tt.expected) {
			t.Errorf("CalculateDropsCdA(%.3f, %.3f): got: %.3f, want: %.3f",
				tt.h, tt.m, actual, tt.expected)
		}
	}
}

func TestCalculateAeroCdA(t *testing.T) {
	tests := []struct {
		h, m, expected float64
	}{
		{1.75, 69, 0.2284},
		{1.60, 53, 0.1982},
	}
	for _, tt := range tests {
		actual := CalculateAeroCdA(tt.h, tt.m)
		if !fequal(actual, tt.expected) {
			t.Errorf("CalculateAeroCdA(%.3f, %.3f): got: %.3f, want: %.3f",
				tt.h, tt.m, actual, tt.expected)
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

func TestAltitudeAdjust(t *testing.T) {
	tests := []struct {
		p, h, expected float64
	}{
		{300, 0, 300},
		{300, 2000, 270},
	}
	for _, tt := range tests {
		actual := AltitudeAdjust(tt.p, tt.h)
		if !fequal(actual, tt.expected) {
			t.Errorf("AltitudeAdjust(%.3f, %.3f): got: %.3f, want: %.3f",
				tt.p, tt.h, actual, tt.expected)
		}
	}
}
