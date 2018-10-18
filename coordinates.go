package main

import (
	"math"
)

type Polar struct {
	t float64 // theta in radians
	r float64
}

func (p Polar) X() float64 {
	return p.r * math.Cos(p.t)
}

func (p Polar) Y() float64 {
	return p.r * math.Sin(p.t)
}

func (p Polar) ToCartesian() Cartesian {
	return Cartesian{
		x: p.X(),
		y: p.Y(),
	}
}



type Cartesian struct {
	x float64
	y float64
}

func (c Cartesian) T() float64 {
	return math.Tanh(c.y / c.x)
}

func (c Cartesian) R() float64 {
	return math.Sqrt(math.Pow(c.x, 2) + math.Pow(c.y, 2))
}

func (c Cartesian) ToPolar() Polar {
	return Polar{
		t: c.T(),
		r: c.R(),
	}
}



func RoundFloatToInt(f float64) int {
	return int(math.Round(f))
}

func RadianToDegree(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

func DegreeToRadian(deg float64) float64 {
	return deg * math.Pi / 180.0
}
