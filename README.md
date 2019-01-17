# Poly

Go CLI application to optimize the shape of a fixed number of polygons to minimize the difference to a target image using a generative algorithm.

Project highly inspired in the work of [Michael Fogleman](https://github.com/fogleman/primitive) and [Roger Johansson](https://rogerjohansson.blog/2008/12/07/genetic-programming-evolution-of-mona-lisa/) but using a different approach and getting different results.

![Example](https://github.com/dmartzol/dmartzol.github.io/raw/master/images/strawberry/p200-n80000.svg)

## Features

* Select the number of polygons used to reproduce the target image.
* Select the number of generations for the algorithm to iterate.

## Installation

Install with:

```
go get -u https://github.com/dmartzol/poly.git
```



## Usage and examples

```
Usage: polygonal [OPTIONS] -o output
  -i string
    	input image path
  -n int
    	number of iterations (default 1000)
  -o value
    	output image path
  -p int
    	number of polygons (default 50)
  -r int
    	resize large input images to this size (default 256)
```

```
poly -i input.png -o output.svg -n 50000
```


## TO DO
- [ ] Implement save a frame every N iterations.
- [ ] Refactor


## License

This project is licensed under the MIT License
