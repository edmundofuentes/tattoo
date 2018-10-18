package main

import (
	"fmt"
	"time"
	"os"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
	"strconv"
)

const C int = 1600
const HC = float64(C / 2)

const INNER_RING_MIN int = 100
const INNER_RING_MAX int = 150

const OUTER_RING_MIN int = 200
const OUTER_RING_MAX int = 700

const MAIN_CIRCLE_RADIUS float64 = 20

const MINOR_CIRCLE_RADIUS float64 = 20

const N_PRINCIPAL_CONNECTIONS int = 3
const N_SECONDARY_CONNECTIONS int = 2
const MAX_DISTANCE_FOR_MAJOR_CONNECTION float64 = 400

const N_MINOR_CIRCLES int = 16
const MINIMUM_MINOR_CONNECTIONS int = 8
const N_MAX_CONNECTIONS_PER_MINOR int = 3
const MIN_DISTANCE_BETWEEN_MINOR_CIRCLES float64 = 200

const STROKE_WIDTH float64 = 5

type Connection struct {
	a int
	b int
}


func main() {
	var seed int64
	if len(os.Args) > 1 {
		seedRaw, err := strconv.Atoi(os.Args[1])

		if err != nil {
			fmt.Println("Invalid seed provided, must be a valid int64.")
			os.Exit(1)
		}

		seed = int64(seedRaw)
	} else {
		seed = time.Now().UTC().UnixNano()
	}

	rand.Seed(seed)
	fmt.Printf("Attempting design from seed %d ..\n", seed)

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Execution timeout. Could not generate a design that matches all the constraints using this seed.")
		os.Exit(1)
	}()

	principal := Polar{
		t: randomRadians(rand.Intn(8)),
		r: float64(rand.Intn(INNER_RING_MAX - INNER_RING_MIN) + INNER_RING_MIN),
	}

	// The second circle should be opposite to the first one
	secondary := principal
	secondary.t += math.Pi


	dc := gg.NewContext(C, C)

	// Draw the white background
	dc.DrawRectangle(0, 0, float64(C), float64(C))
	dc.SetRGB(1, 1, 1)
	dc.Fill()

	// Draw the main 2 circles
	drawMainCircle(dc, principal)
	drawMainCircle(dc, secondary)

	// Connect the main 2 circles
	dc.DrawLine(principal.X()+HC, principal.Y()+HC, secondary.X()+HC, secondary.Y()+HC)
	dc.SetRGB(0,0,0)
	dc.SetLineWidth(STROKE_WIDTH * 2)
	dc.Stroke()


	// GENERATE THE MINOR CIRCLES
	minors := make([]Polar, 0)

	i := 0
	for {
		center := Polar{
			t: randomRadians(i),
			r: float64(rand.Intn(OUTER_RING_MAX - OUTER_RING_MIN) + OUTER_RING_MIN),
		}

		// check the distance to other points
		tooClose := false
		for _, c := range minors {
			if DistanceBetweenPolars(center, c) < MIN_DISTANCE_BETWEEN_MINOR_CIRCLES {
				tooClose = true
				break
			}
		}

		if tooClose {
			continue
		}

		minors = append(minors, center)

		if len(minors) == N_MINOR_CIRCLES {
			break
		}

		i++
	}



	// Randomize the connection between the principals
	principalConnections := make([]int, 0)

	ConnectPrincipalLoop:
	for {
		// Select a minor circles at random
		m := rand.Intn(len(minors))

		// skip if the distance between them is bigger than 4 octants (4 * π/4)
		if math.Abs( math.Mod(principal.t, 2*math.Pi) - math.Mod(minors[m].t, 2*math.Pi)) > (math.Pi) {
			continue
		}

		// check if the connection is repeated
		for _, c := range principalConnections {
			if c == m {
				continue ConnectPrincipalLoop
			}
		}

		// check that the distance between the principal and the minor is not too big
		if DistanceBetweenPolars(principal, minors[m]) > MAX_DISTANCE_FOR_MAJOR_CONNECTION {
			continue
		}

		// create a connection
		principalConnections = append(principalConnections, m)

		// end the loop once we have reached the minimum number of connections
		if len(principalConnections) == N_PRINCIPAL_CONNECTIONS {
			break
		}
	}

	secondaryConnections := make([]int, 0)


	ConnectSecondaryLoop:
	for {
		// Select a minor circles at random
		m := rand.Intn(len(minors))

		// skip if the distance between them is bigger than 4 octants (4 * π/4)
		if math.Abs( math.Mod(secondary.t, 2*math.Pi) - math.Mod(minors[m].t, 2*math.Pi)) > (math.Pi) {
			continue
		}

		// check if the connection is repeated
		for _, c := range secondaryConnections {
			if c == m {
				continue ConnectSecondaryLoop
			}
		}

		// check that the distance between the secondary and the minor is not too big
		if DistanceBetweenPolars(secondary, minors[m]) > MAX_DISTANCE_FOR_MAJOR_CONNECTION {
			continue
		}

		// create a connection
		secondaryConnections = append(secondaryConnections, m)

		// end the loop once we have reached the minimum number of connections
		if len(secondaryConnections) == N_SECONDARY_CONNECTIONS {
			break
		}
	}

	// Randomize the connection between the minors
	connections := make([]Connection, 0)

	ConnectMinorCirclesLoop:
	for {

		// Select to minor circles at random
		a := rand.Intn(len(minors))
		b := rand.Intn(len(minors))

		// if they are the same, skip
		if a == b {
			continue
		}

		// skip if the distance between them is bigger than 2 octants (2 * π/4)
		if math.Abs( math.Mod(minors[a].t, 2*math.Pi) - math.Mod(minors[b].t, 2*math.Pi)) > (2 * math.Pi / 4) {
			continue
		}

		// check if the connection is repeated
		for _, c := range connections {
			if (c.a == a && c.b == b) || (c.a == b && c.b == a) {
				continue ConnectMinorCirclesLoop
			}
		}

		// check how many times each node has a connection
		totalA := 0
		totalB := 0
		for _, c := range connections {
			if c.a == a || c.b == a {
				totalA++
			}
			if c.a == b || c.b == b {
				totalB++
			}
		}

		if totalA >= N_MAX_CONNECTIONS_PER_MINOR || totalB >= N_MAX_CONNECTIONS_PER_MINOR {
			continue
		}

		// create a connection
		connections = append(connections, Connection{a, b})

		// end the loop once we have reached the minimum number of connections
		if len(connections) >= MINIMUM_MINOR_CONNECTIONS && checkThatAllPointsHaveAtLeastOneConnection(N_MINOR_CIRCLES, connections){
			break
		}
	}

	// Draw all the principal connections
	for _, conn := range principalConnections {
		dc.DrawLine(principal.X()+HC, principal.Y()+HC, minors[conn].X()+HC, minors[conn].Y()+HC)
		dc.SetRGB(0,0,0)
		dc.SetLineWidth(STROKE_WIDTH)
		dc.Stroke()
	}

	// Draw all the secondary connections
	for _, conn := range secondaryConnections {
		dc.DrawLine(secondary.X()+HC, secondary.Y()+HC, minors[conn].X()+HC, minors[conn].Y()+HC)
		dc.SetRGB(0,0,0)
		dc.SetLineWidth(STROKE_WIDTH)
		dc.Stroke()
	}

	// Draw all the minor connections
	for _, conn := range connections {
		dc.DrawLine(minors[conn.a].X()+HC, minors[conn.a].Y()+HC, minors[conn.b].X()+HC, minors[conn.b].Y()+HC)
		dc.SetRGB(0,0,0)
		dc.SetLineWidth(STROKE_WIDTH)
		dc.Stroke()
	}

	// Draw all the minor circles
	for _, circle := range minors {
		drawMinorCircle(dc, circle)
	}


	os.MkdirAll("./output", os.ModePerm)
	filename := fmt.Sprintf("./output/%d.png", seed)

	err := dc.SavePNG(filename)

	if err != nil {
		fmt.Println("Error writing output PNG file. Check permissions and try again.")
	}

	fmt.Printf("Success! Image generated on output/%d.png\n", seed)
	os.Exit(0)
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

// this is HIGHLY inefficient, oh well.
func checkThatAllPointsHaveAtLeastOneConnection(n int, connections []Connection) bool {
	for i := 0; i < n; i++ {
		total := 0

		for _, c := range connections {
			if c.a == i || c.b == i {
				total++
				break
			}
		}

		if total == 0 {
			return false
		}
	}

	return true
}

func randomRadians(i int) float64 {
	// one octant (1/8 of a circle) = π/4

	// calculate current octant
	o := i % 8

	// randomize position inside an octant
	r := rand.Float64() * (math.Pi / 4.0)

	// add both of them together
	return  r + (float64(o) * (math.Pi / 4.0))
}
