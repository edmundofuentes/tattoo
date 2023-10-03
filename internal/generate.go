package internal

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

func Generate(cfg Config) *Constellation {

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Execution timeout. Could not generate a design that matches all the constraints using this seed.")
		os.Exit(1)
	}()

	constellation := NewConstellation()

	// we will calculate 3 special "nodes" that we will call "Primary"
	// the very first node will be outside the bounds of the connections circle
	externalNode := Node{
		polar: Polar{
			t: DegreeToRadian(270),
			r: float64(cfg.ExternalNodeDistance), // this will always be perfectly at North
		},
		nodeType: NodeType_EXTERNAL,
	}
	constellation.appendNode(externalNode)

	// then we will randomly place a "pair" within the INNER ring
	// calculate the alignment offset
	theta := (rand.Float64() * float64(cfg.PrimaryNodesAlignmentMaxOffsetDegrees-cfg.PrimaryNodesAlignmentMinOffsetDegrees)) + float64(cfg.PrimaryNodesAlignmentMinOffsetDegrees)
	theta = 270 - theta

	internalNodeA := Node{
		polar: Polar{
			//t: RandomRadianInOctant(7),
			t: DegreeToRadian(theta),
			r: float64(rand.Intn(cfg.InnerRingMaxRadius-cfg.InnerRingMinRadius) + cfg.InnerRingMinRadius),
		},
		nodeType: NodeType_PRIMARY,
	}

	// The second node should be opposite to the first one
	internalNodeB := internalNodeA
	internalNodeB.polar.t += math.Pi

	constellation.appendNode(internalNodeA)
	constellation.appendNode(internalNodeB)

	// Finally, link all the primary nodes
	constellation.appendLink(Link{
		nodeA:    externalNode,
		nodeB:    internalNodeA,
		linkType: LinkType_PRIMARY,
	})

	constellation.appendLink(Link{
		nodeA:    internalNodeA,
		nodeB:    internalNodeB,
		linkType: LinkType_PRIMARY,
	})

	fmt.Println("> Generated primary nodes!")

	// Generate all the other Secondary nodes
	// determine how many nodes we will actually place
	numberOfSecondaryNodes := rand.Intn(cfg.MaxNumberOfSecondaryNodes-cfg.MinNumberOfSecondaryNodes+1) + cfg.MinNumberOfSecondaryNodes
	fmt.Printf("> Attempting to place %d secondary nodes\n", numberOfSecondaryNodes)
	quadrant := 0

	for {
		quadrant++

		node := Node{
			polar: Polar{
				t: RandomRadianInOctant(quadrant),
				r: float64(rand.Intn(cfg.OuterRingMaxRadius-cfg.OuterRingMinRadius) + cfg.OuterRingMinRadius),
			},
			nodeType: NodeType_SECONDARY,
		}

		tooClose := false
		for _, c := range constellation.nodes {
			if DistanceBetweenPolars(node.polar, c.polar) < cfg.MinDistanceBetweenNodes {
				tooClose = true
				break
			}
		}

		if tooClose {
			continue
		}

		constellation.appendNode(node)

		fmt.Printf("· Placed secondary node! Total: %d\n", len(constellation.nodes)-3)

		if len(constellation.nodes)-3 >= numberOfSecondaryNodes {
			break
		}
	}

	// Finally, place the links
	//constellation = placeLinksRandomly(cfg, constellation)
	constellation = placeLinksByNearest(cfg, constellation)

	/*


		ConnectPrincipalLoop:
			for {
				// Select a minor circles at random
				m := rand.Intn(len(minors))

				// skip if the distance between them is bigger than 4 octants (4 * π/4)
				if math.Abs(math.Mod(principal.t, 2*math.Pi)-math.Mod(minors[m].t, 2*math.Pi)) > (math.Pi) {
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
				if math.Abs(math.Mod(secondary.t, 2*math.Pi)-math.Mod(minors[m].t, 2*math.Pi)) > (math.Pi) {
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
				if math.Abs(math.Mod(minors[a].t, 2*math.Pi)-math.Mod(minors[b].t, 2*math.Pi)) > (2 * math.Pi / 4) {
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
				if len(connections) >= MINIMUM_MINOR_CONNECTIONS && checkThatAllPointsHaveAtLeastOneConnection(N_MINOR_CIRCLES, connections) {
					break
				}
			}

	*/

	return constellation
}

func placeLinksRandomly(cfg Config, constellation *Constellation) *Constellation {
	numberOfLinks := rand.Intn(cfg.MaxNumberOfLinks-cfg.MinNumberOfLinks+1) + cfg.MinNumberOfLinks
	numberOfLinks = numberOfLinks - 2 // remove the initial two links that have already been placed

	for {
		// Select two random nodes
		// node A will be prioritized from a list of nodes that do not YET have a link
		nodesWithoutLinks := constellation.nodesWithoutLinks()
		var nodeA Node
		if len(nodesWithoutLinks) != 0 {
			nodeA = nodesWithoutLinks[rand.Intn(len(nodesWithoutLinks))]
		} else {
			// all nodes have at least one link, yes! we can use the default random method
			nodeA = constellation.nodes[rand.Intn(len(constellation.nodes))]
		}

		nodeB := constellation.nodes[rand.Intn(len(constellation.nodes))]

		if nodeA == nodeB {
			// same node, continue
			continue
		}

		if nodeA.nodeType == NodeType_EXTERNAL || nodeB.nodeType == NodeType_EXTERNAL {
			// the original external node cannot have more than one link
			continue
		}

		// Look up how many links a node has
		if constellation.countLinksForNode(nodeA) > cfg.MaxLinksPerNode {
			continue
		}

		if constellation.countLinksForNode(nodeB) > cfg.MaxLinksPerNode {
			continue
		}

		// So far so good, prepare the link
		link := Link{
			nodeA:    nodeA,
			nodeB:    nodeB,
			linkType: LinkType_SECONDARY,
		}

		// do one final check to see if the link is not repeated
		if constellation.similarLinkExists(link) {
			continue
		}

		constellation.appendLink(link)

		fmt.Printf("· Placed link! Total: %d\n", len(constellation.links)-2)

		if len(constellation.links)-2 >= numberOfLinks {
			break
		}
	}

	return constellation
}

func placeLinksByNearest(cfg Config, constellation *Constellation) *Constellation {
	numberOfLinks := rand.Intn(cfg.MaxNumberOfLinks-cfg.MinNumberOfLinks+1) + cfg.MinNumberOfLinks
	numberOfLinks = numberOfLinks - 2 // remove the initial two links that have already been placed

	for {
		// Select two random nodes
		// node A will be prioritized from a list of nodes that do not YET have a link
		nodesWithoutLinks := constellation.nodesWithoutLinks()
		var nodeA Node
		if len(nodesWithoutLinks) != 0 {
			nodeA = nodesWithoutLinks[rand.Intn(len(nodesWithoutLinks))]
		} else {
			// all nodes have at least one link, yes! we can use the default random method
			nodeA = constellation.nodes[rand.Intn(len(constellation.nodes))]
		}

		// Perform checks on NodeA first
		if nodeA.nodeType == NodeType_EXTERNAL {
			continue
		}

		if constellation.countLinksForNode(nodeA) > cfg.MaxLinksPerNode {
			continue
		}

		// then, we will order the nodes by nearest neighbor from the source node
		nodesOrderedByDistanceToNode := constellation.nodesOrderedByDistanceFromPolar(nodeA.polar)

		// and we will try to link it
		for _, nodeB := range nodesOrderedByDistanceToNode {
			if nodeA == nodeB {
				// same node, continue
				continue
			}

			if nodeB.nodeType == NodeType_EXTERNAL {
				// the original external node cannot have more than one link
				continue
			}

			// Look up how many links a node has
			if constellation.countLinksForNode(nodeB) > cfg.MaxLinksPerNode {
				continue
			}

			// So far so good, prepare the link
			link := Link{
				nodeA:    nodeA,
				nodeB:    nodeB,
				linkType: LinkType_SECONDARY,
			}

			// do one final check to see if the link is not repeated
			if constellation.similarLinkExists(link) {
				continue
			}

			constellation.appendLink(link)

			fmt.Printf("· Placed link! Total: %d\n", len(constellation.links)-2)

			break
		}

		if len(constellation.links)-2 >= numberOfLinks {
			break
		}
	}

	return constellation
}
