package main

import (
	"fmt"
	"time"
	"os"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

const C int = 1600
const HC = C / 2

const INNER_RING_MIN int = 100
const INNER_RING_MAX int = 150

const OUTER_RING_MIN int = 200
const OUTER_RING_MAX = HC - 100

const MAIN_CIRCLE_RADIUS float64 = 20

const MINOR_CIRCLE_RADIUS float64 = 20

const STROKE_WIDTH float64 = 5


func main() {
	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)
	fmt.Println(seed)

	primaryCircleCenter := Polar{
		t: randomRadians(),
		r: float64(rand.Intn(INNER_RING_MAX - INNER_RING_MIN) + INNER_RING_MIN),
	}

	// The second circle should be opposite to the first one
	secondaryCircleCenter := primaryCircleCenter
	secondaryCircleCenter.t += math.Pi



	dc := gg.NewContext(C, C)

	// Draw the white background
	dc.DrawRectangle(0, 0, float64(C), float64(C))
	dc.SetRGB(1, 1, 1)
	dc.Fill()

	drawMainCircle(dc, primaryCircleCenter)
	drawMainCircle(dc, secondaryCircleCenter)


	for i := 0; i < 10; i++ {
		center := Polar{
			t: randomRadians(),
			r: float64(rand.Intn(OUTER_RING_MAX - OUTER_RING_MIN) + OUTER_RING_MIN),
		}

		drawMinorCircle(dc, center)
	}


	os.MkdirAll("./output", os.ModePerm)
	filename := fmt.Sprintf("./output/%d.png", seed)

	err := dc.SavePNG(filename)

	if err != nil {
		fmt.Println("Error writing output PNG file. Check permissions and try again.")
	}
}

func drawMainCircle(dc *gg.Context, p Polar) {
	dc.DrawCircle(p.X() +  float64(HC), p.Y() + float64(HC), MAIN_CIRCLE_RADIUS)
	dc.SetRGB(0,0,0)
	dc.Fill()
}

func drawMinorCircle(dc *gg.Context, p Polar) {
	dc.DrawCircle(p.X() +  float64(HC), p.Y() + float64(HC), MINOR_CIRCLE_RADIUS)
	dc.SetRGB(0,0,0)
	dc.Fill()

	dc.DrawCircle(p.X() +  float64(HC), p.Y() + float64(HC), MINOR_CIRCLE_RADIUS - STROKE_WIDTH)
	dc.SetRGB(1, 1, 1)
	dc.Fill()
}

func randomRadians() float64 {
	return rand.Float64() * 2 * math.Pi
}
