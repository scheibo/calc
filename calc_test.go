package calc

import (
	"testing"
)

func TestPowerTOT(t *testing.T) {
	tests := []struct {
		rho, cda, crr, va, vg, gr, mt, r, vgi, vgf, ti, tf, g, ec, fw, i, expected float64
	}{
		{Rho0, DropsCdA, Crr, 5.55, 5.55, 0.08125, 75.0, R700x23, 5.55, 5.55, 0, 864.865, G, Ec, Fw, I, 389.9},
		{1.1921, TopsCdA, 0.008, 2.327, 4.293, 0.079, 85.0, R700x28, 0, 11.1111, 0, 3240, G, 0.95, Fw, 0.12, 334.98},
	}
	for _, tt := range tests {
		actual := PowerTOT(tt.rho, tt.cda, tt.crr, tt.va, tt.vg, tt.gr, tt.mt, tt.r, tt.vgi, tt.vgf, tt.ti, tt.tf, tt.g, tt.ec, tt.fw, tt.i)
		if !Eqf(actual, tt.expected) {
			t.Errorf("PowerTOT(%.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f): got: %.3f, want: %.3f",
				tt.rho, tt.cda, tt.crr, tt.va, tt.vg, tt.gr, tt.mt, tt.r, tt.vgi, tt.vgf, tt.ti, tt.tf, tt.g, tt.ec, tt.fw, tt.i, actual, tt.expected)
		}
	}
}

func TestVelocity(t *testing.T) {
	tests := []struct {
		p, rho, cda, crr, vw, dw, db, gr, mt, g, ec, fw, expected float64
	}{
		{389.9, Rho0, DropsCdA, Crr, 0, 0, 0, 0.08125, 75.0, G, Ec, Fw, 5.55},
		{333.175, 1.1921, TopsCdA, 0.008, 2.7778, 180, 45, 0.079, 85.0, G, 0.95, Fw, 4.293},
	}
	for _, tt := range tests {
		actual := Velocity(tt.p, tt.rho, tt.cda, tt.crr, tt.vw, tt.dw, tt.db, tt.gr, tt.mt, tt.g, tt.ec, tt.fw)
		if !Eqf(actual, tt.expected) {
			t.Errorf("Velocity(%.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f): got: %.3f, want: %.3f",
				tt.p, tt.rho, tt.cda, tt.crr, tt.vw, tt.dw, tt.db, tt.gr, tt.mt, tt.g, tt.ec, tt.fw, actual, tt.expected)
		}
	}
}

func TestTime(t *testing.T) {
	tests := []struct {
		p, d, rho, cda, crr, vw, dw, db, gr, mt, g, ec, fw, expected float64
	}{
		{389.9, 4800, Rho0, DropsCdA, Crr, 0, 0, 0, 0.08125, 75.0, G, Ec, Fw, 864.865},
		{333.175, 13910, 1.1921, TopsCdA, 0.008, 2.78, 180, 45, 0.079, 85.0, G, 0.95, Fw, 3240},
	}
	for _, tt := range tests {
		actual := Time(tt.p, tt.d, tt.rho, tt.cda, tt.crr, tt.vw, tt.dw, tt.db, tt.gr, tt.mt, tt.g, tt.ec, tt.fw)
		if !Eqf(actual, tt.expected) {
			t.Errorf("Time(%.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f): got: %.3f, want: %.3f",
				tt.p, tt.d, tt.rho, tt.cda, tt.crr, tt.vw, tt.dw, tt.db, tt.gr, tt.mt, tt.g, tt.ec, tt.fw, actual, tt.expected)
		}
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		p, t, rho, cda, crr, vw, dw, db, gr, mt, g, ec, fw, expected float64
	}{
		{389.9, 864.865, Rho0, DropsCdA, Crr, 0, 0, 0, 0.08125, 75.0, G, Ec, Fw, 4800},
		{333.175, 3240, 1.1921, TopsCdA, 0.008, 2.78, 180, 45, 0.079, 85.0, G, 0.95, Fw, 13910},
	}
	for _, tt := range tests {
		actual := Distance(tt.p, tt.t, tt.rho, tt.cda, tt.crr, tt.vw, tt.dw, tt.db, tt.gr, tt.mt, tt.g, tt.ec, tt.fw)
		if !Eqf(actual, tt.expected) {
			t.Errorf("Distance(%.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f): got: %.3f, want: %.3f",
				tt.p, tt.t, tt.rho, tt.cda, tt.crr, tt.vw, tt.dw, tt.db, tt.gr, tt.mt, tt.g, tt.ec, tt.fw, actual, tt.expected)
		}
	}
}

func TestPowerAT(t *testing.T) {
	tests := []struct {
		rho, cda, fw, va, vg, expected float64
	}{
		{Rho0, 0.2565, Fw, 10.91, 8.36, 158.8},
		{1.1921, TopsCdA, Fw, 2.327, 4.293, 5.6},
	}
	for _, tt := range tests {
		actual := PowerAT(tt.rho, tt.cda, tt.fw, tt.va, tt.vg)
		if !Eqf(actual, tt.expected) {
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
		{4.293, 0.079, 0.008, 85, G, 28.55},
	}
	for _, tt := range tests {
		actual := PowerRR(tt.vg, tt.gr, tt.crr, tt.mt, tt.g)
		if !Eqf(actual, tt.expected) {
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
		{4.293, 0.551},
	}
	for _, tt := range tests {
		actual := PowerWB(tt.vg)
		if !Eqf(actual, tt.expected) {
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
		{4.293, 85, G, 0.079, 281.93},
	}
	for _, tt := range tests {
		actual := PowerPE(tt.vg, tt.mt, tt.g, tt.gr)
		if !Eqf(actual, tt.expected) {
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
		{85, 0.12, R700x28, 0, 11.1111, 0, 3240, 1.64},
	}
	for _, tt := range tests {
		actual := PowerKE(tt.mt, tt.i, tt.r, tt.vgi, tt.vgf, tt.ti, tt.tf)
		if !Eqf(actual, tt.expected) {
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
		{4.293, 2.78, 180, 45, 2.327},
	}
	for _, tt := range tests {
		actual := AirVelocity(tt.vg, tt.vw, tt.dw, tt.db)
		if !Eqf(actual, tt.expected) {
			t.Errorf("AirVelocity(%.3f, %.3f, %.3f, %.3f): got: %.3f, want: %.3f",
				tt.vg, tt.vw, tt.dw, tt.db, actual, tt.expected)
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
		if !Eqf(actual, tt.expected) {
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
		if !Eqf(actual, tt.expected) {
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
		if !Eqf(actual, tt.expected) {
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
		if !Eqf(actual, tt.expected) {
			t.Errorf("AirPressure(%.3f, %.3f): got: %.3f, want: %.3f",
				tt.h, tt.t, actual, tt.expected)
		}
	}
}

func TestAirDensity(t *testing.T) {
	tests := []struct {
		h, g, expected float64
	}{
		{0, G, Rho0},
		{1000, G, 1.111},
	}
	for _, tt := range tests {
		actual := AirDensity(tt.h, tt.g)
		if !Eqf(actual, tt.expected) {
			t.Errorf("AirDensity(%.3f, %.3f): got: %.3f, want: %.3f",
				tt.h, tt.g, actual, tt.expected)
		}
	}
}

func TestAltitudeAdjust(t *testing.T) {
	tests := []struct {
		p, h, expected float64
	}{
		{300, 0, 301},
		{300, 2000, 270},
	}
	for _, tt := range tests {
		actual := AltitudeAdjust(tt.p, tt.h)
		if !Eqf(actual, tt.expected) {
			t.Errorf("AltitudeAdjust(%.3f, %.3f): got: %.3f, want: %.3f",
				tt.p, tt.h, actual, tt.expected)
		}
	}
}
