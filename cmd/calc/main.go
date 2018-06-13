package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/scheibo/calc"
)

// RADII maps from tire width to radius
var RADII = map[int64]float64{
	20: calc.R700x20,
	22: calc.R700x22,
	23: calc.R700x23,
	25: calc.R700x25,
	28: calc.R700x28,
}

// COMPASS maps from cardinal direction to degrees
var COMPASS = map[string]float64{
	"N":   0,
	"NNE": 22.5,
	"NE":  45,
	"ENE": 67.5,
	"E":   90,
	"ESE": 122.5,
	"SE":  135,
	"SSE": 157.5,
	"S":   180,
	"SSW": 202.5,
	"SW":  225,
	"WSE": 247.5,
	"W":   270,
	"WNW": 292.5,
	"NW":  315,
	"NNW": 337.5,
}

func verify(s string, x float64) {
	if x < 0 {
		// print and exit "%s must be non negative"
		_ = s
	}
}

// TODO add support for altitude and adjust for lower power
func main() {
	// TODO calc.Power/calc.Time method should calculate va and take vw/dw/db instead
	var rho, cda, crr, vg, vw, dw, db, e, gr, mt, mr, mb, r, t, d, p float64
	var dwS, dbS string
	var tire int64
	var err error

	flag.Float64Var(&rho, "rho", calc.Rho0, "air density in kg/m*3")
	flag.Float64Var(&cda, "cda", 0.325, "coefficient of drag area")
	flag.Float64Var(&crr, "crr", calc.Crr, "coefficient of rolling resistance")

	flag.Float64Var(&mr, "mr", 67.0, "total mass of the rider in kg")
	flag.Float64Var(&mb, "mb", 8.0, "total mass of the bicycle in kg")

	flag.Int64Var(&tire, "tire", 23, "the tire width in mm")

	flag.Float64Var(&vw, "vw", 0, "the wind speed")
	flag.StringVar(&dwS, "dw", "N", "the cardinal direction the wind travelling (*not* its origin)")
	flag.StringVar(&dbS, "db", "N", "the cardinal direction the bicycle is travelling")

	flag.Float64Var(&e, "e", 0, "total elevation gained in m")
	flag.Float64Var(&gr, "gr", 0, "average grade")

	flag.Float64Var(&d, "d", -1, "distance travelled in m")
	flag.Float64Var(&p, "p", -1, "power in watts")
	flag.Float64Var(&p, "t", -1, "duration in s")

	flag.Parse()

	verify("rho", rho)
	verify("cda", cda)
	verify("crr", crr)

	verify("mr", mr)
	verify("mb", mb)
	mt = mr + mb

	r, err = tireRadius(tire)
	if err != nil {
		// print and exit
	}

	verify("vw", vw)
	if vw > 0 {
		dw, err = parseDirection(dwS)
		if err != nil {
			// print and exit
		}
		db, err = parseDirection(dbS)
		if err != nil {
			// print and exist
		}
	}

	verify("gr", gr)
	// error correct in case grade was passed in as a %
	if gr > 1 {
		gr = gr / 100
	}

	// if both are specified, make sure they agree
	if e > 0 && gr > 0 &&
		((d*gr != e) || (e/d != gr)) {
		// print and exit
	}

	if d <= 0 {
		// print and exit
	}

	if p != -1 {
		verify("p", p)
		if t != -1 {
			// print and exit
		}
		// TODO calculate time (vg, then use d to get t)
		_ = rho * cda * crr * vg * vw * dw * db * e * gr * mt * r * d * p
	}

	if t != -1 {
		verify("t", t)
		if p != -1 {
			// print and exit
		}
		vg = calc.Vg(d, time.Duration(t)*time.Second)
		// TODO calculate power!, display components and W/kg
		_ = rho * cda * crr * vg * vw * dw * db * e * gr * mt * r * d * t
	}
}

func tireRadius(tire int64) (float64, error) {
	r, ok := RADII[tire]
	if !ok {
		return 0, fmt.Errorf("invalid tire width '%d'", tire)
	}
	return r, nil
}

func parseDirection(dir string) (float64, error) {
	d, ok := COMPASS[strings.ToUpper(dir)]
	if !ok {
		return 0, fmt.Errorf("invalid direction '%s'", dir)
	}
	return d, nil
}
