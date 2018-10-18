# Tattoo

A simple Go program that generates a random design that might be tattooed on myself.

It takes an `int64` seed as the first command argument and generates a pseudo-random design from it.

This project uses the [Go Graphics](https://github.com/fogleman/gg) library by [Michael Fogleman](http://www.michaelfogleman.com/).

## Sample

```bash
$> go build
$> ./tattoo 1539903654746537380
Attempting design from seed 1539903654746537380 ..
Success! Image generated on output/1539903654746537380.png
```

![Sample design using seed 1539903654746537380](https://github.com/edmundofuentes/tattoo/blob/master/sample/1539903654746537380.png)