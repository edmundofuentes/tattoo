package internal

import (
	"sort"
)

type Constellation struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func NewConstellation() *Constellation {
	return &Constellation{
		Nodes: make([]Node, 0),
		Links: make([]Link, 0),
	}
}

func (constellation *Constellation) appendNode(node Node) {
	constellation.Nodes = append(constellation.Nodes, node)
}

func (constellation *Constellation) appendLink(link Link) {
	constellation.Links = append(constellation.Links, link)
}

type NodeType string

const (
	NodeType_EXTERNAL  NodeType = "external"
	NodeType_PRIMARY   NodeType = "primary"
	NodeType_SECONDARY NodeType = "secondary"
)

type LinkType string

const (
	LinkType_PRIMARY   LinkType = "primary"
	LinkType_SECONDARY LinkType = "secondary"
)

type Node struct {
	Polar    Polar    `json:"polar"`
	NodeType NodeType `json:"node_type"`
}

type Link struct {
	NodeA    Node     `json:"node_a"`
	NodeB    Node     `json:"node_b"`
	LinkType LinkType `json:"link_type"`
}

func (constellation *Constellation) countLinksForNode(node Node) int {
	c := 0
	for _, link := range constellation.Links {
		if link.NodeA == node || link.NodeB == node {
			c++
		}
	}

	return c
}

func (constellation *Constellation) similarLinkExists(link Link) bool {
	for _, l := range constellation.Links {
		if (l.NodeA == link.NodeA && l.NodeB == link.NodeB) ||
			(l.NodeA == link.NodeB && l.NodeB == link.NodeA) {
			return true
		}
	}

	return false
}

func (constellation *Constellation) nodesWithLessThanXLinks(x int) []Node {
	nodes := make([]Node, 0)

	for _, node := range constellation.Nodes {
		// Skip the External (original) node
		if node.NodeType == NodeType_EXTERNAL {
			continue
		}

		if constellation.countLinksForNode(node) < x {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (constellation *Constellation) nodesOrderedByDistanceFromPolar(polar Polar) []Node {
	nodes := make([]Node, len(constellation.Nodes))
	copy(nodes, constellation.Nodes)

	sort.Slice(nodes, func(a, b int) bool {
		return DistanceBetweenPolars(polar, nodes[a].Polar) < DistanceBetweenPolars(polar, nodes[b].Polar)
	})

	return nodes
}

// this is HIGHLY inefficient, oh well.
/*
func (constellation *Constellation) checkThatAllPointsHaveAtLeastOneConnection(n int, connections []Connection) bool {
	// TODO: we should use our own Nodes
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
*/

/*
func LinksIntersect(linkA Link, linkB Link) bool {
	var ccw = func(A Polar, B Polar, C Polar) bool {
		return (C.Y()-A.Y())*(B.X()-A.X()) > (B.Y()-A.Y())*(C.X()-A.X())
	}

	var A = linkA.NodeA.Polar
	var B = linkA.NodeB.Polar
	var C = linkB.NodeA.Polar
	var D = linkB.NodeB.Polar

	return (ccw(A, C, D) != ccw(B, C, D)) && (ccw(A, B, C) != ccw(A, B, D))
}
*/

func LinksIntersect(linkA Link, linkB Link) bool {
	var A = linkA.NodeA.Polar
	var B = linkA.NodeB.Polar
	var C = linkB.NodeA.Polar
	var D = linkB.NodeB.Polar

	var ccw = func(A Polar, B Polar, C Polar) bool {
		return (C.Y()-A.Y())*(B.X()-A.X()) > (B.Y()-A.Y())*(C.X()-A.X())
	}

	var A = linkA.NodeA.Polar
	var B = linkA.NodeB.Polar
	var C = linkB.NodeA.Polar
	var D = linkB.NodeB.Polar

	return (ccw(A, C, D) != ccw(B, C, D)) && (ccw(A, B, C) != ccw(A, B, D))
}

func (constellation *Constellation) linkIntersectsExisting(link Link) bool {
	for _, existingLink := range constellation.Links {
		if LinksIntersect(existingLink, link) {
			return true
		}
	}

	return false
}
