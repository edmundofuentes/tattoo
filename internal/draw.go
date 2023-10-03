package internal

import (
	"github.com/fogleman/gg"
)

func Draw(cfg Config, constellation *Constellation) *gg.Context {
	dc := gg.NewContext(cfg.CanvasSize, cfg.CanvasSize)

	// Draw the white background
	dc.DrawRectangle(0, 0, float64(cfg.CanvasSize), float64(cfg.CanvasSize))
	dc.SetRGB(1, 1, 1)
	dc.Fill()

	// DRAW THE LINKS
	// Links are drawn first because some nodes have no fill, so we want to make sure that the node's
	// circles are drawn AFTER the link's lines have been drawn
	for _, link := range constellation.links {
		drawLink(cfg, dc, link)
	}

	// Draw the Nodes
	for _, node := range constellation.nodes {
		drawNode(cfg, dc, node)
	}

	return dc
}

func drawLink(cfg Config, dc *gg.Context, link Link) {
	HC := cfg.HalfCanvas()
	dc.DrawLine(link.nodeA.polar.X()+HC, link.nodeA.polar.Y()+HC, link.nodeB.polar.X()+HC, link.nodeB.polar.Y()+HC)
	dc.SetRGB(0, 0, 0)

	if link.linkType == LinkType_PRIMARY {
		dc.SetLineWidth(cfg.PrimaryLinkStrokeWidth)
	} else {
		dc.SetLineWidth(cfg.SecondaryLinkStrokeWidth)
	}

	dc.Stroke()
}

func drawNode(cfg Config, dc *gg.Context, node Node) {
	if node.nodeType == NodeType_PRIMARY || node.nodeType == NodeType_EXTERNAL {
		// Primary Nodes are solid filled
		dc.DrawCircle(node.polar.X()+cfg.HalfCanvas(), node.polar.Y()+cfg.HalfCanvas(), cfg.PrimaryNodeRadius)
		dc.SetRGB(0, 0, 0)
		dc.Fill()
	} else {
		// Secondary Nodes have a stroke and are empty inside
		dc.DrawCircle(node.polar.X()+cfg.HalfCanvas(), node.polar.Y()+cfg.HalfCanvas(), cfg.SecondaryNodeRadius)
		dc.SetRGB(0, 0, 0)
		dc.Fill()

		dc.DrawCircle(node.polar.X()+cfg.HalfCanvas(), node.polar.Y()+cfg.HalfCanvas(), cfg.SecondaryNodeRadius-cfg.SecondaryNodeStrokeWidth)
		dc.SetRGB(1, 1, 1)
		dc.Fill()
	}
}
