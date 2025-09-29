# Tattoo

A simple Go program that generates a random design that might be tattooed on myself.

It takes an `int64` seed as the first command argument and generates a pseudo-random design from it.

This project uses the [Go Graphics](https://github.com/fogleman/gg) library by [Michael Fogleman](http://www.michaelfogleman.com/).

```bash
$> go build -o tattoo cmd/tattoo.go
```

## v2 (2023)
![Sample design using seed 1696298614598191000](https://github.com/edmundofuentes/tattoo/blob/master/sample/1696298614598191000.png)

![Sample design using seed 1696299702148896000](https://github.com/edmundofuentes/tattoo/blob/master/sample/1696299702148896000.png)

```toml
# General design SAMPLE
seed = 1696298614598191000

canvas_size = 1600

external_node_distance = 750 # this must be less than half of the canvas size

primary_nodes_alignment_min_offset_degrees = 15 # +- degrees to offset against the vertical axis (270º)
primary_nodes_alignment_max_offset_degrees = 25 # +- degrees to offset against the vertical axis (270º)

inner_ring_min_radius = 75 # the center primary Nodes will be at least DOUBLE this distance
inner_ring_max_radius = 100 # the center primary Nodes will be at most DOUBLE this distance

outer_ring_min_radius = 150
outer_ring_max_radius = 450

min_number_of_secondary_nodes = 6
max_number_of_secondary_nodes = 9
min_distance_between_nodes = 150

min_number_of_links = 12 # this includes the initial *2* Links created as base
max_number_of_links = 16
max_links_per_node = 3

# Drawing
primary_link_stroke_width = 9.0
secondary_link_stroke_width = 6.5 # in the original it was double the regular

primary_node_radius = 20.0 # primary node are filled in

secondary_node_radius = 20.0 # secondary Nodes have the same stroke_width, the outher radius is defined here
secondary_node_stroke_width = 6.5
```

Every parameter can now be set from the `config.toml`, including the seed.

```bash
$> go run cmd/tattoo.go
Attempting design from seed 1696298614598191000 ..
Success! Image generated on output/1696298614598191000.png
```


### General design notes
- The _constellation_ is initialized with 3 Nodes.
  - The first one is called "external" and it's always located at the very top of the canvas perfectly centered
  - The other two Nodes are called "primary" and they are closer to the center of the canvas, as defined by the `inner_ring` parameters. Those two Nodes are always 180º opposed from each other, and the angle against the vertical axis can be offset from the config.
- All other Nodes are called "secondary" and placed afterward within the bounds of the `outer_ring` with a `min_distance_between_nodes`
- Finally, Links (connections between Nodes) are placed. There are two algorithms for this:
  - Random.  Nodes are selected at random and are attempted to link. A link can fail if the `max_links_per_node` is exceeded.
  - Closest Neighbor. A node is selected at random, and then it tries to link itself to its closest neighbor. If it's already connected, it will attempt the second-closest and so forth.





## v1 (2018)
![Sample design using seed 1539903654746537380](https://github.com/edmundofuentes/tattoo/blob/master/sample/1539903654746537380.png)

```bash
$> go build
$> ./tattoo 1539903654746537380
Attempting design from seed 1539903654746537380 ..
Success! Image generated on output/1539903654746537380.png
```
