package internal

import (
	"math"
	"math/rand"
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

func DistanceBetweenPolars(a, b Polar) float64 {
	// https://math.tutorvista.com/geometry/distance-formula-for-polar-coordinates.html
	return math.Sqrt(math.Pow(a.r, 2) + math.Pow(b.r, 2) - 2*a.r*b.r*math.Cos(a.t-b.t))
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

func RandomRadianInOctant(octant int) float64 {
	// one octant (1/8 of a circle) = π/4

	// calculate current octant
	o := octant % 8

	// randomize position inside an octant
	r := rand.Float64() * (math.Pi / 4.0)

	// add both of them together
	return r + (float64(o) * (math.Pi / 4.0))
}
