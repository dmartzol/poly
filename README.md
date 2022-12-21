# Poly

Go CLI application to optimize the shape of a fixed number of polygons to minimize the difference to a target image using a genetic algorithm just for fun.

Project highly inspired in the work of [Michael Fogleman](https://github.com/fogleman/primitive) and [Roger Johansson](https://rogerjohansson.blog/2008/12/07/genetic-programming-evolution-of-mona-lisa/) but using a different approach and getting different results.

### Target image
![Target](https://github.com/dmartzol/dmartzol.github.io/raw/master/images/strawberry/strawberry.png)

### Resulting image
![Result](https://github.com/dmartzol/dmartzol.github.io/raw/master/images/strawberry/p200-n80000.png)
Image generated with 200 polygons and 80000 iterations.

## Features

* Select the number of polygons used to reproduce the target image.
* Select the number of generations for the algorithm to iterate.

## Installation

Install with:

```
go get -u https://github.com/dmartzol/poly.git
```



## Usage and examples

### Options
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

To generate an image with 200 polygons and 50k iterations input:
```
poly -i input.png -o output.svg -n 50000 -p 200
```

## TO DO
- [ ] Concurrency
- [ ] Implement save a frame every N iterations.
- [ ] Add more examples.
- [ ] Improve this README.

## License

This project is licensed under the MIT License

## Examples

![Target](https://github.com/dmartzol/dmartzol.github.io/raw/master/images/win/400poly64000n.png)
