package internal

type Config struct {
	Seed int64 `toml:"seed"`

	CanvasSize int `toml:"canvas_size"`

	ExternalNodeDistance int `toml:"external_node_distance"`

	PrimaryNodesAlignmentMinOffsetDegrees int `toml:"primary_nodes_alignment_min_offset_degrees"`
	PrimaryNodesAlignmentMaxOffsetDegrees int `toml:"primary_nodes_alignment_max_offset_degrees"`

	InnerRingMinRadius int `toml:"inner_ring_min_radius"`
	InnerRingMaxRadius int `toml:"inner_ring_max_radius"`

	OuterRingMinRadius int `toml:"outer_ring_min_radius"`
	OuterRingMaxRadius int `toml:"outer_ring_max_radius"`

	MinNumberOfSecondaryNodes int     `toml:"min_number_of_secondary_nodes"`
	MaxNumberOfSecondaryNodes int     `toml:"max_number_of_secondary_nodes"`
	MinDistanceBetweenNodes   float64 `toml:"min_distance_between_nodes"`

	MinNumberOfLinks int `toml:"min_number_of_links"`
	MaxNumberOfLinks int `toml:"max_number_of_links"`
	MaxLinksPerNode  int `toml:"max_links_per_node"`

	PrimaryLinkStrokeWidth   float64 `toml:"primary_link_stroke_width"`
	SecondaryLinkStrokeWidth float64 `toml:"secondary_link_stroke_width"`

	PrimaryNodeRadius float64 `toml:"primary_node_radius"`

	SecondaryNodeRadius      float64 `toml:"secondary_node_radius"`
	SecondaryNodeStrokeWidth float64 `toml:"secondary_node_stroke_width"`
}

func (cfg Config) HalfCanvas() float64 {
	return float64(cfg.CanvasSize / 2)
}
