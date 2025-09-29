package internal

import (
	"fmt"
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

	// we will calculate 4 special "Nodes" that we will call "Primary"
	// the very first node will be outside the bounds of the connections circle
	externalNode := Node{
		Polar: Polar{
			t: DegreeToRadian(270),
			r: float64(cfg.PrimaryExternalNodeDistance), // this will always be perfectly at North
		},
		NodeType: NodeType_EXTERNAL,
	}
	constellation.appendNode(externalNode)

	// then we will place 3 within the INNER ring
	internalNodeA := Node{
		Polar: Polar{
			t: DegreeToRadian(270),
			r: float64(cfg.PrimaryInternalNodesDistance),
		},
		NodeType: NodeType_PRIMARY,
	}

	// The second and third Nodes are at equidistant angles
	internalNodeB := internalNodeA
	internalNodeB.Polar.t += DegreeToRadian(120)

	internalNodeC := internalNodeA
	internalNodeC.Polar.t += DegreeToRadian(-120)

	constellation.appendNode(internalNodeA)
	constellation.appendNode(internalNodeB)
	constellation.appendNode(internalNodeC)

	// Finally, link all the primary Nodes
	constellation.appendLink(Link{
		NodeA:    externalNode,
		NodeB:    internalNodeA,
		LinkType: LinkType_PRIMARY,
	})

	constellation.appendLink(Link{
		NodeA:    internalNodeA,
		NodeB:    internalNodeB,
		LinkType: LinkType_PRIMARY,
	})

	constellation.appendLink(Link{
		NodeA:    internalNodeA,
		NodeB:    internalNodeC,
		LinkType: LinkType_PRIMARY,
	})

	fmt.Println("> Generated primary Nodes!")

	// Place secondary nodes
	constellation = placeSecondaryNodes(cfg, constellation)

	// Place the links
	//constellation = placeLinksByNearest(cfg, constellation)
	constellation = placeLinksAll(cfg, constellation)

	return constellation
}

func placeSecondaryNodes(cfg Config, constellation *Constellation) *Constellation {
	// Generate all the other Secondary Nodes
	// determine how many Nodes we will actually place
	numberOfSecondaryNodes := rand.Intn(cfg.MaxNumberOfSecondaryNodes-cfg.MinNumberOfSecondaryNodes+1) + cfg.MinNumberOfSecondaryNodes
	fmt.Printf("> Attempting to place %d secondary Nodes\n", numberOfSecondaryNodes)
	sector := 0

	for {
		sector++

		// skip the vertical sector, ensure the external node is not blocked by a secondary node
		if (sector % cfg.NumberOfSectorsForSecondaryNodes) == 0 {
			continue
		}

		// maybe we also skip this sector? 10% chance
		if rand.Float64() < 0.10 {
			continue
		}

		node := Node{
			Polar: Polar{
				t: RandomRadianInSector(sector, cfg.NumberOfSectorsForSecondaryNodes),
				r: float64(rand.Intn(cfg.SecondaryNodesMaxRadius-cfg.SecondaryNodesMinRadius) + cfg.SecondaryNodesMinRadius),
			},
			NodeType: NodeType_SECONDARY,
		}

		tooClose := false
		for _, c := range constellation.Nodes {
			if DistanceBetweenPolars(node.Polar, c.Polar) < cfg.MinDistanceBetweenNodes {
				tooClose = true
				break
			}
		}

		if tooClose {
			continue
		}

		constellation.appendNode(node)

		//fmt.Printf("· Placed secondary node! Total: %d\n", len(constellation.Nodes)-4)

		if len(constellation.Nodes)-4 >= numberOfSecondaryNodes {
			break
		}
	}

	fmt.Printf("· Placed secondary Nodes! Total: %d\n", len(constellation.Nodes)-4)

	return constellation
}

func placeLinksAll(cfg Config, constellation *Constellation) *Constellation {

	for _, nodeA := range constellation.Nodes {
		if nodeA.NodeType == NodeType_EXTERNAL {
			// the original external node cannot have more than one link
			continue
		}

		for _, nodeB := range constellation.Nodes {
			if nodeA == nodeB {
				// same node, continue
				continue
			}

			if nodeB.NodeType == NodeType_EXTERNAL {
				// the original external node cannot have more than one link
				continue
			}

			if nodeA.NodeType == NodeType_PRIMARY && nodeB.NodeType == NodeType_PRIMARY {
				// prevent linking the inner Primary nodes
				continue
			}

			// So far so good, prepare the link
			link := Link{
				NodeA:    nodeA,
				NodeB:    nodeB,
				LinkType: LinkType_SECONDARY,
			}

			// check to see if the link is not repeated
			if constellation.similarLinkExists(link) {
				continue
			}

			// check that this link doesn't intersect any other link
			// we checked if the same link exists above, so that prevents a collinearity exception to the check
			if constellation.linkIntersectsExisting(link) {
				continue
			}

			constellation.appendLink(link)
		}
	}

	return constellation
}

func placeLinksByNearest(cfg Config, constellation *Constellation) *Constellation {
	numberOfLinks := rand.Intn(cfg.MaxNumberOfLinks-cfg.MinNumberOfLinks+1) + cfg.MinNumberOfLinks
	numberOfLinks = numberOfLinks - 3 // remove the initial three Links that have already been placed

	fmt.Printf("> Attempting to place %d secondary Links\n", numberOfLinks)

	for {
		// Select two random Nodes
		// node A will be prioritized from a list of Nodes that do not YET have a link
		nodesWithoutLinks := constellation.nodesWithLessThanXLinks(cfg.MinLinksPerNode) // TODO: make this configurable
		var nodeA Node
		if len(nodesWithoutLinks) != 0 {
			nodeA = nodesWithoutLinks[rand.Intn(len(nodesWithoutLinks))]
		} else {
			// all Nodes have at least X links, yes! we can use the default random method
			nodeA = constellation.Nodes[rand.Intn(len(constellation.Nodes))]
		}

		if nodeA.NodeType == NodeType_EXTERNAL {
			continue
		}

		fmt.Printf("· Selected node %f %f %d\n", nodeA.Polar.t, nodeA.Polar.r, constellation.countLinksForNode(nodeA))

		if constellation.countLinksForNode(nodeA) >= cfg.MaxLinksPerNode {
			continue
		}

		// then, we will order the Nodes by nearest neighbor from the source node
		nodesOrderedByDistanceToNode := constellation.nodesOrderedByDistanceFromPolar(nodeA.Polar)

		// and we will try to link it
		for _, nodeB := range nodesOrderedByDistanceToNode {
			if nodeA == nodeB {
				// same node, continue
				continue
			}

			if nodeB.NodeType == NodeType_EXTERNAL {
				// the original external node cannot have more than one link
				continue
			}

			if nodeA.NodeType == NodeType_PRIMARY && nodeB.NodeType == NodeType_PRIMARY {
				// prevent linking the inner Primary nodes
				continue
			}

			// Look up how many Links the other node has
			if constellation.countLinksForNode(nodeB) >= cfg.MaxLinksPerNode {
				continue
			}

			// So far so good, prepare the link
			link := Link{
				NodeA:    nodeA,
				NodeB:    nodeB,
				LinkType: LinkType_SECONDARY,
			}

			// check to see if the link is not repeated
			if constellation.similarLinkExists(link) {
				continue
			}

			// check that this link doesn't intersect any other link
			// we checked if the same link exists above, so that prevents a collinearity exception to the check
			if constellation.linkIntersectsExisting(link) {
				continue
			}

			constellation.appendLink(link)

			//fmt.Printf("· Placed link! Total: %d\n", len(constellation.Links)-3)

			break
		}

		if len(constellation.Links)-3 >= numberOfLinks {
			break
		}
	}

	fmt.Printf("· Placed secondary Links! Total: %d\n", len(constellation.Links)-3)

	return constellation
}
