# calc

![release-candidate](http://img.shields.io/badge/status-release--candidate-green.svg)&nbsp;
[![Build Status](http://img.shields.io/travis/scheibo/calc.svg)](https://travis-ci.org/scheibo/calc)

calc implements the mathematical formulas for modelling road cycling power as
outlined by the paper ['Validation of a Mathematical Model for Road Cycling
Power'][1] by Martin et al and other miscellaenous cycling power related
formulas.

> $ go install github.com/scheibo/calc
> $ ./calc -t=16m05s -d=4800 -gr=8.125 -mr=70
> 16:05 (4.80 km @ 8.12%) = 357.37 W (5.11 W/kg) = AT:25.44 W + RR:15.54 W + WB:0.68 W + PE:315.70 W

The generated GoDoc can be viewed at [godoc.org/github.com/scheibo/calc][2].

[1]: https://www.ncbi.nlm.nih.gov/pubmed/28121252
[2]: https://godoc.org/github.com/scheibo/calc
