package calc

import (
	"math"
)

// G is the acceleration of gravity in metres per second squared.
const G = 9.80665

// Ec is the drive chain efficiency factor that accounts for the
// frictional losses in power.
const Ec = 0.976

// Fw is the factor associated with wheel rotation that represents the
// incremental drag area of the spokes in squared metres.
const Fw = 0.0044

// I is the moment of inertia of the two wheels on the bicycle in
// kilograms per metre squared.
const I = 0.14

// Tire radii in metres associated with standard road cycling tire widths.
const (
	R700x20 = 0.331
	R700x22 = 0.333
	R700x23 = 0.334
	R700x25 = 0.336
	R700x28 = 0.339
)

// Crr is the typical coefficient of rolling resistance for bicycle tires on
// a smooth asphalt road.
const Crr = 0.004

// Typical measurements of aerodynamic drag (CdA) for various cycling positions.
const (
	TopsCdA     = 0.400
	HoodsCdA    = 0.350 // 0.3240
	DropsCdA    = 0.310 // 0.3070, 0.3019
	RoadAeroCdA = 0.290 // 0.2914, 0.2662
	TTAeroCdA   = 0.270 // 0.2680, 0.2427 (0.2323 w/ Aero Helmet)
)

// Typical measurments coefficients of drag (Cd) for various cycling positions.
const (
	DropsCd = 0.88
	AeroCd  = 0.70
)

// Rho0 is the 'standard' dry, sea level air density at T0 in kg/m*3
const Rho0 = 1.225

// T0 is the sea level standard temperature in Kelvin.
const T0 = 288.15

// P0 is the air pressure at sea level (1 atm) in Pa.
const P0 = 101325

// M is the molar mass of air in kg/mol.
const M = 0.0289644

// R is the universal gas constant in J/(mol*K)
const R = 8.31432

// K is the value of Kelvin corresponding to 0 Celsius.
const K = 273.15

// L is the temperatue lapse rate in the troposphere in K/m.
const L = 0.0065

// Components is the amount of power in watts each component of the model requires.
type Components struct {
	AT float64
	RR float64
	WB float64
	PE float64
	KE float64
}

// Ptot calculates the total power required, equal to the net total power
// of Pat, Prr, Pwb, Ppe, and Pke divided by the drive chain efficiency ec.
func Ptot(rho, cda, crr, va, vg, gr, mt, r, vgi, vgf, ti, tf, g, ec, fw, i float64) float64 {
	comp := Pcomp(rho, cda, crr, va, vg, gr, mt, r, vgi, vgf, ti, tf, g, ec, fw, i)
	return comp.AT + comp.RR + comp.WB + comp.PE + comp.KE
}

// PowerTOT is an alias for the Ptot function.
var PowerTOT = Ptot

// Psimp calculates a simplified version of the total power required, equal to
// the net total power of Pat, Prr, Pwb, Ppe, divided by the drive chain efficiency ec,
// but without contributions from Pke.
func Psimp(rho, cda, crr, va, vg, gr, mt, g, ec, fw float64) float64 {
	// NOTE: (tf - ti) must not equal 0 so we use tf = 1
	comp := Pcomp(rho, cda, crr, va, vg, gr, mt, 0, 0, 0, 0, 1, g, ec, fw, 0)
	return comp.AT + comp.RR + comp.WB + comp.PE // comp.KE = 0
}

// Power is an alias for the Psimp function.
var Power = Psimp

// Vg calculates the velocity of the bicycle relative to the ground in m/s based
// on the net total power p given rho, cda, crr, vw, dw, db, gr, mt, g, ec and fw.
// NOTE: this method is only valid for velocities between 0 and 100 m/s.
func Vg(p, rho, cda, crr, vw, dw, db, gr, mt, g, ec, fw float64) float64 {
	// epsilon is some small value that determines when we will stop the search
	const epsilon = 1e-6
	// max is the maxmium number of iterations of the search
	const max = 100

	vgl, vgm, vgh := 0.0, 50.0, 100.0
	for j := 0; j < max; j++ {
		pm := Psimp(rho, cda, crr, Va(vgm, vw, dw, db), vgm, gr, mt, g, ec, fw)
		if Eqf(pm, p, epsilon) {
			break
		}

		if pm > p {
			vgh = vgm
		} else {
			vgl = vgm
		}

		vgm = (vgh + vgl) / 2.0
	}

	return vgm
}

// GroundVelocity is an alias for the Vg function.
var GroundVelocity = Vg

// Velocity is an alias for the Vg function.
var Velocity = Vg

// T calculates the duration in seconds of a performance over distance d in metres
// with net total power p given rho, cda, crr, vw, dw, db, gr, mt, g, ec and fw.
// NOTE: this method is only valid for velocities between 0 and 100 m/s.
func T(p, d, rho, cda, crr, vw, dw, db, gr, mt, g, ec, fw float64) float64 {
	return d / Vg(p, rho, cda, crr, vw, dw, db, gr, mt, g, ec, fw)
}

// Time is an alias for the T function.
var Time = T

// D calculates the distance in metres of a performance over duration t in seconds
// with net total power p given rho, cda, crr, vw, dw, db, gr, mt, g, ec and fw.
// NOTE: this method is only valid for velocities between 0 and 100 m/s.
func D(p, t, rho, cda, crr, vw, dw, db, gr, mt, g, ec, fw float64) float64 {
	return t * Vg(p, rho, cda, crr, vw, dw, db, gr, mt, g, ec, fw)
}

// Distance is an alias for the D function.
var Distance = D

// Pcomp calculates the total power required, broken down by the components of
// Pat, Prr, Pwb, Ppe, and Pke, each divided by the drive chain efficiency ec.
func Pcomp(rho, cda, crr, va, vg, gr, mt, r, vgi, vgf, ti, tf, g, ec, fw, i float64) Components {
	return Components{
		AT: Pat(rho, cda, fw, va, vg) / ec,
		RR: Prr(vg, gr, crr, mt, g) / ec,
		WB: Pwb(vg) / ec,
		PE: Ppe(vg, mt, g, gr) / ec,
		KE: Pke(mt, i, r, vgi, vgf, ti, tf) / ec,
	}
}

// PowerCOMP is an alias for the Pcomp function
var PowerCOMP = Pcomp

// Pat calculates the power to overcome the force due to total aerodynamic
// drag given the air density rho, the coefficient of drag multiplied by the
// drag area cda, the factor associated with wheel rotation that represents
// the incremental drag area of the spokes fw, the air velocity of bicycle va
// (the ground velocity of the bicycle added to the component of wind velocity
// tangent to the direction of travel of the bicycle) and vg, the ground
// velocity of the bicycle.
func Pat(rho, cda, fw, va, vg float64) float64 {
	return 0.5 * rho * (cda + fw) * (math.Pow(va, 2)) * vg
}

// PowerAT is an alias for the Pat function.
var PowerAT = Pat

// Prr calculates the power to overcome the force due to rolling resistance
// given the velocity of the rider relative to the ground vg, the road
// gradient gr (rise/run), the coefficient of rolling resistance crr, the
// total mass of the rider and bike mt, and the acceleration of gravity g.
func Prr(vg, gr, crr, mt, g float64) float64 {
	return vg * math.Cos(math.Atan(gr)) * crr * mt * g
}

// PowerRR is an alias for the Prr function.
var PowerRR = Prr

// Pwb calculates the power to overcome the frictional losses associated with
// the bicycle wheel bearings given the ground velocity of the bicycle vg.
func Pwb(vg float64) float64 {
	return vg * (91 + 8.7*vg) * 0.001
}

// PowerWB is an alias for the Pwb function.
var PowerWB = Pwb

// Ppe calculates the power associated with changes in potential energy, i.e.
// the force required to overcome gravity. This is calculated from the ground
// velocity of the bicyle vg, the total mass of the rider and the bike mt, the
// acceleration of gravity g and the road gradient gr (rise/run).
func Ppe(vg, mt, g, gr float64) float64 {
	return vg * mt * g * math.Sin(math.Atan(gr))
}

// PowerPE is an alias for the Ppe function.
var PowerPE = Ppe

// Pke calculates the power related to changes in kinetic energy given the total
// mass of the rider and the bike mt, the moment of inertia of the two wheels i,
// the outside radius of the tire r, the initial and final ground velocities of
// the bicycle vgi and vgf and the initial and final times ti and tf.
func Pke(mt, i, r, vgi, vgf, ti, tf float64) float64 {
	return 0.5 * (mt + i/math.Pow(r, 2)) * (math.Pow(vgf, 2) - math.Pow(vgi, 2)) / (tf - ti)
}

// PowerKE is an alias for the Pke function.
var PowerKE = Pke

// Va returns the air velocity of the bicycle given the velocity of the bicycle
// relative to the ground vg, the absolute wind velocity vw, the wind direction
// dw and the direction of travel of the bicycle db, both in degrees.
func Va(vg, vw, dw, db float64) float64 {
	return vg + (vw * math.Cos((dw*math.Pi/180)-(db*math.Pi/180)))
}

// AirVelocity is an alias for the Va function.
var AirVelocity = Va

// Yaw calculates the yaw angle of the bike and rider relative to the wind given
// the air velocity va, the absolute wind velocity vw, the wind direction dw and
// the direction of travel of the bicycle db, both in degrees.
func Yaw(va, vw, dw, db float64) float64 {
	return math.Atan(vw*math.Sin((dw*math.Pi/180)-(db*math.Pi/180))/va) * 180 / math.Pi
}

// CalculateDropsCdA calculates the estimated typical aerodynamic drag of a rider
// in the drops with a height of h and a mass of m, assuming zero yaw.
func CalculateDropsCdA(h, m float64) float64 {
	return DropsCd * DropsA(h, m)
}

// CalculateAeroCdA calculates the estimated typical aerodynamic drag of a rider in
// the aerobars with a height of h and a mass of m, assuming zero yaw.
func CalculateAeroCdA(h, m float64) float64 {
	return AeroCd * AeroA(h, m)
}

// DropsA calculates the estimated typical front area of a rider in the drops
// with a height of h and a mass of m, assuming zero yaw.
func DropsA(h, m float64) float64 {
	return A(h, m, 0.0276, 0.1647)
}

// AeroA calculates the estimated typical front area of a rider in the aerobars
// with a height of h and a mass of m, assuming zero yaw.
func AeroA(h, m float64) float64 {
	return A(h, m, 0.0293, 0.0604)
}

// A calculates the estimated typical frontal area of a rider of height h and
// mass m given constants a and b which vary depending on position, assuming
// zero yaw.
func A(h, m, a, b float64) float64 {
	return a*math.Pow(h, 0.725)*math.Pow(m, 0.425) + b
}

// AirPressure implements the barometric formula for calculating the air pressure at
// altitude h metres high where the temperature at the altitude is t in Celsius.
func AirPressure(h, t float64) float64 {
	return P0 * math.Exp((-G*M*h)/(R*(t+K)))
}

// Rho calculates the air density at a given altitude h in metres and the acceleration
// due to gravity g using the ideal gas law.
func Rho(h, g float64) float64 {
	t := T0 - (L * h)
	p := P0 * math.Pow((1-((L*h)/T0)), ((g*M)/(R*L)))
	return (p * M) / (R * t)
}

// AirDensity is an alias for the Rho function.
var AirDensity = Rho

// AltitudeAdjust calculates the equivalent sustainable power at altitude h metres
// compared to a seal level power of p based on the formula derived from Clark et al
// "The effect of acute simulated moderate altitude on power, performances and pacing
// strategies in well-trained cyclists".
func AltitudeAdjust(p, h float64) float64 {
	x := h / 1000
	return p * ((-0.0092 * math.Pow(x, 2)) - (0.0323 * x) + 1)
}

func Eqf(a, b float64, eps ...float64) bool {
	e := 1e-3
	if len(eps) > 0 {
		e = eps[0]
	}
	// min is the smallest normal value possible
	const min = float64(2.2250738585072014E-308) // 1 / 2**(1022)

	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)

	if a == b {
		return true
	} else if a == b || b == 0 || diff < min {
		// a or b is zero or both are extremely close to it relative error is less meaningful here
		return diff < (e * min)
	} else {
		// use relative error
		return diff/(absA+absB) < e
	}
}
