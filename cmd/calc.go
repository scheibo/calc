// calc provides a CLI for calculating either the power required or the time
// achievable for a given performance.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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
	"WSW": 247.5,
	"W":   270,
	"WNW": 292.5,
	"NW":  315,
	"NNW": 337.5,
}

type DirectionFlag struct {
	Direction float64
}

func (df *DirectionFlag) String() string {
	return strconv.FormatFloat(df.Direction, 'f', -1, 64)
}

func (df *DirectionFlag) Set(v string) error {
	d, ok := COMPASS[strings.ToUpper(v)]
	if ok {
		df.Direction = d
		return nil
	}

	d, err := strconv.ParseFloat(v, 64)
	if err == nil {
		df.Direction = d
		return nil
	}

	return fmt.Errorf("invalid direction '%s'", v)
}

func main() {
	var rho, cda, crr, vw, e, gr, h, mt, mr, mb, r, t, d, p float64
	var dw, db DirectionFlag
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
	flag.Var(&dw, "dw", "the cardinal direction the wind originates from")
	flag.Var(&db, "db", "the cardinal direction the bicycle is travelling")

	flag.Float64Var(&e, "e", 0, "total elevation gained in m")
	flag.Float64Var(&gr, "gr", 0, "average grade")
	flag.Float64Var(&h, "h", 0, "median elevation")

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
	verify("h", h)
	if h != 0 {
		r := calc.Rho(h, calc.G)
		// if both are specified, make sure they agree
		if rho != calc.Rho0 && r != rho {
			exit(fmt.Errorf("specified both rho=%f and h=%f but they do not agree", rho, h))
		}
		rho = r
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

	fi, _ := os.Stdout.Stat()
	tty := (fi.Mode() & os.ModeCharDevice) == 0

	if p != -1 {
		verify("p", p)
		if dur != -1 {
			exit(fmt.Errorf("t and p can't both be provided"))
		}

		t = calc.T(p, d, rho, cda, crr, vw, dw.Direction, db.Direction, gr, mt, calc.G, calc.Ec, calc.Fw)
		dur = time.Duration(t) * time.Second
		wkg := p / mr

		if !tty {
			fmt.Println(dur)
		} else {
			fmt.Printf("%.2f km @ %.2f%% @ %.2f W (%.2f W/kg) = %s\n", d/1000, gr*100, p, wkg, fmtDuration(dur))
		}
	} else if dur != -1 {
		verify("t", float64(dur))
		t = float64(dur / time.Second)
		if p != -1 {
			exit(fmt.Errorf("p and t can't both be provided"))
		}

		vg := d / t
		va := calc.Va(vg, vw, dw.Direction, db.Direction)

		comp := calc.Pcomp(rho, cda, crr, va, vg, gr, mt, r, vg, vg, 0, t, calc.G, calc.Ec, calc.Fw, calc.I)
		ptot := comp.AT + comp.RR + comp.WB + comp.PE + comp.KE
		wkg := ptot / mr

		if !tty {
			fmt.Println(ptot)
		} else {
			fmt.Printf("%s (%.2f km @ %.2f%%) = %.2f W (%.2f W/kg) = AT:%.2f W + RR:%.2f W + WB:%.2f W + PE:%.2f W\n",
				fmtDuration(dur), d/1000, gr*100, ptot, wkg, comp.AT, comp.RR, comp.WB, comp.PE)
		}
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

func verify(s string, x float64) {
	if x < 0 {
		exit(fmt.Errorf("%s must be non negative but was %f", s, x))
	}
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	flag.PrintDefaults()
	os.Exit(1)
}
