package internal

type Config struct {
	Seed int64 `toml:"seed"`

	CanvasSize int `toml:"canvas_size"`

	PrimaryExternalNodeDistance  int `toml:"primary_external_node_distance"`
	PrimaryInternalNodesDistance int `toml:"primary_internal_nodes_distance"`

	SecondaryNodesMinRadius int `toml:"secondary_nodes_min_radius"`
	SecondaryNodesMaxRadius int `toml:"secondary_nodes_max_radius"`

	MinNumberOfSecondaryNodes        int     `toml:"min_number_of_secondary_nodes"`
	MaxNumberOfSecondaryNodes        int     `toml:"max_number_of_secondary_nodes"`
	MinDistanceBetweenNodes          float64 `toml:"min_distance_between_nodes"`
	NumberOfSectorsForSecondaryNodes int     `toml:"number_of_sectors_for_secondary_nodes"`

	MinNumberOfLinks int `toml:"min_number_of_links"`
	MaxNumberOfLinks int `toml:"max_number_of_links"`
	MinLinksPerNode  int `toml:"min_links_per_node"`
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
