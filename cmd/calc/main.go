package main

import (
	"flag"
	"fmt"

	"github.com/scheibo/calc"
)

// TODO add support for altitude and adjust for lower power
func main() {
	var rho, cda, crr, va, vg, gr, mt, mr, mb, r, t, p float64

	va = 0
	gr = 0
	r = calc.R700x23

	flag.Float64Var(&rho, "rho", calc.Rho0, "air density in kg/m*3")
	flag.Float64Var(&cda, "cda", 0.325, "coefficient of drag area")
	flag.Float64Var(&crr, "crr", calc.Crr, "coefficient of rolling resistance")

	// TODO print defaults but also have multi parsing?
	var mass, weight, mtS, mrS, mbS string
	flag.StringVar(&mass, "mass", "", "total mass of the rider and bicycle")
	flag.StringVar(&weight, "weight", "", "total mass of the rider and bicycle")
	flag.StringVar(&mtS, "mt", "", "total mass of the rider and bicycle")
	flag.StringVar(&mrS, "mr", "", "total mass of the rider")
	flag.StringVar(&mbS, "mb", "", "total mass of the bicycle")
	mr, mb = parseMass(mass, weight, mtS, mrS, mbS)
	mt = mr + mb

	var tire, width int64
	flag.Int64Var(&tire, "tire", 23, "the tire radius in mm")
	flag.Int64Var(&tire, "width", 23, "the tire radius in mm")



	// NEED:
	// (1) distance d
	// (2a) grade gr or (2b) elevation e (can default to 0)
	// () time t = calculate p and also display W/kg or () power p = calculate t



	// 1) take in lb or kg
	// 2) come up with mt
	//var mass, weight, mt float64 // mt
	//var mr, mb float64

	//// must be 20/22/23/25/28, use to calculate r
	//var tire, width int64

	//// or calculate from d/gain, take % or 0.214
	//var grade, gradient, gr float64

	//// metres or feet, can do gr instead
	//var elevation, e, h, gain float64

	//// metres or feet
	//var distance, d float64

	//var rho, cda, crr float64 // default to specific values

	//// can be duration or seconds
	//var time, t float64

	//// km/h & kph or m/s or mi/h & mph
	//var wind, vw float64
	//// in degrees, or as N/NW/NNW/etc
	//var dw, db float64

	//var watts, power, p float64



}

func parseMass(mass, weight, mtS, mrS, mbS string) (mr, mb float64, err error) {
	mr = 67.0
	mb = 8.0

	return mr, mb, nil
}

var COMPASS map[string]float64{
	"N": 0,
	"NNE": 22.5,
	"NE": 45,
	"ENE": 67.5,
	"E": 90,
	"ESE": 122.5,
	"SE": 135,
	"SSE": 157.5,
	"S": 180,
	"SSW": 202.5,
	"SW": 225,
	"WSE": 247.5,
	"W": 270,
	"WNW": 292.5,
	"NW": 315,
	"NNW": 337.5,
}

func parseDirection(dir string) (float64, error) {
	d, ok := COMPASS[strings.ToUpper(dir)], !ok {
	  return nil, fmt.Errorf("invalid direction '%s'", dir)
	}
	return d, nil
}
