package internal

import (
	"sort"
)

type Constellation struct {
	nodes []Node
	links []Link
}

func NewConstellation() *Constellation {
	return &Constellation{
		nodes: make([]Node, 0),
		links: make([]Link, 0),
	}
}

func (constellation *Constellation) appendNode(node Node) {
	constellation.nodes = append(constellation.nodes, node)
}

func (constellation *Constellation) appendLink(link Link) {
	constellation.links = append(constellation.links, link)
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
	polar    Polar
	nodeType NodeType
}

type Link struct {
	nodeA    Node
	nodeB    Node
	linkType LinkType
}

func (constellation *Constellation) countLinksForNode(node Node) int {
	c := 0
	for _, link := range constellation.links {
		if link.nodeA == node || link.nodeB == node {
			c++
		}
	}

	return c
}

func (constellation *Constellation) similarLinkExists(link Link) bool {
	for _, l := range constellation.links {
		if (l.nodeA == link.nodeA && l.nodeB == link.nodeB) ||
			(l.nodeA == link.nodeB && l.nodeB == link.nodeA) {
			return true
		}
	}

	return false
}

func (constellation *Constellation) nodesWithoutLinks() []Node {
	nodes := make([]Node, 0)

	for _, node := range constellation.nodes {
		if constellation.countLinksForNode(node) == 0 {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (constellation *Constellation) nodesOrderedByDistanceFromPolar(polar Polar) []Node {
	nodes := make([]Node, len(constellation.nodes))
	copy(nodes, constellation.nodes)

	sort.Slice(nodes, func(a, b int) bool {
		return DistanceBetweenPolars(polar, nodes[a].polar) < DistanceBetweenPolars(polar, nodes[b].polar)
	})

	return nodes
}

// this is HIGHLY inefficient, oh well.
/*
func (constellation *Constellation) checkThatAllPointsHaveAtLeastOneConnection(n int, connections []Connection) bool {
	// TODO: we should use our own nodes
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
