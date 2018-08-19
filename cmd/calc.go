// calc provides a CLI for calculating either the power required or the time
// achievable for a given performance.
package main

import (
	"flag"
	"fmt"
	"os"
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

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	flag.PrintDefaults()
	os.Exit(1)
}

func verify(s string, x float64) {
	if x < 0 {
		exit(fmt.Errorf("%s must be non negative but was %f", s, x))
	}
}

// BUG(kjs): add support for altitude and adjust for lower power
func main() {
	var rho, cda, crr, vw, dw, db, e, gr, mt, mr, mb, r, t, d, p float64
	var dwS, dbS string
	var tire int64
	var err error
	var dur time.Duration

	flag.Float64Var(&rho, "rho", calc.Rho0, "air density in kg/m*3")
	flag.Float64Var(&cda, "cda", 0.325, "coefficient of drag area")
	flag.Float64Var(&crr, "crr", calc.Crr, "coefficient of rolling resistance")

	flag.Float64Var(&mr, "mr", 67.0, "total mass of the rider in kg")
	flag.Float64Var(&mb, "mb", 8.0, "total mass of the bicycle in kg")

	flag.Int64Var(&tire, "tire", 23, "the tire width in mm")

	flag.Float64Var(&vw, "vw", 0, "the wind speed in m/s")
	flag.StringVar(&dwS, "dw", "N", "the cardinal direction the wind originates from")
	flag.StringVar(&dbS, "db", "N", "the cardinal direction the bicycle is travelling")

	flag.Float64Var(&e, "e", 0, "total elevation gained in m")
	flag.Float64Var(&gr, "gr", 0, "average grade")

	flag.Float64Var(&d, "d", -1, "distance travelled in m")
	flag.Float64Var(&p, "p", -1, "power in watts")
	flag.DurationVar(&dur, "t", -1, "duration in minutes and seconds ('12m34s')")

	flag.Parse()

	verify("rho", rho)
	verify("cda", cda)
	verify("crr", crr)

	verify("mr", mr)
	verify("mb", mb)
	mt = mr + mb

	r, err = tireRadius(tire)
	if err != nil {
		exit(err)
	}

	verify("vw", vw)
	if vw > 0 {
		dw, err = parseDirection(dwS)
		if err != nil {
			exit(err)
		}
		db, err = parseDirection(dbS)
		if err != nil {
			exit(err)
		}
	}

	// error correct in case grade was passed in as a %
	if gr > 1 || gr < -1 {
		gr = gr / 100
	}

	if d <= 0 {
		exit(fmt.Errorf("d must be positive but was %f", d))
	}

	if e > 0 {
		// if both are specified, make sure they agree
		if gr > 0 && ((d*gr != e) || (e/d != gr)) {
			exit(fmt.Errorf("specified both e=%f and gr=%f but they do not agree", e, gr))
		}
		gr = e / d
	}

	if p != -1 {
		verify("p", p)
		if dur != -1 {
			exit(fmt.Errorf("t and p %.2f can't both be provided", p))
		}

		t = calc.T(p, d, rho, cda, crr, vw, dw, db, gr, mt, calc.G, calc.Ec, calc.Fw)
		dur = time.Duration(t) * time.Second
		wkg := p / mr

		fmt.Printf("%.2f km @ %.2f%% @ %.2f W (%.2f W/kg) = %s\n", d/1000, gr*100, p, wkg, fmtDuration(dur))
	} else if t != -1 {
		verify("t", float64(dur))
		t = float64(dur / time.Second)
		if p != -1 {
			exit(fmt.Errorf("p and t can't both be provided"))
		}

		vg := d / t
		va := calc.Va(vg, vw, dw, db)

		comp := calc.Pcomp(rho, cda, crr, va, vg, gr, mt, r, vg, vg, 0, t, calc.G, calc.Ec, calc.Fw, calc.I)
		ptot := comp.AT + comp.RR + comp.WB + comp.PE + comp.KE
		wkg := ptot / mr

		fmt.Printf("%s (%.2f km @ %.2f%%) = %.2f W (%.2f W/kg) = AT:%.2f W + RR:%.2f W + WB:%.2f W + PE:%.2f W\n",
			fmtDuration(dur), d/1000, gr*100, ptot, wkg, comp.AT, comp.RR, comp.WB, comp.PE)
	} else {
		exit(fmt.Errorf("p or t must be specified"))
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
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
