package poly

import (
	"flag"
	"fmt"
	"image"
	"log"

	// Decoding jpg images
	"image/draw"
	_ "image/jpeg"
	"os"
)

func PrintDefaultsWithError(errorMessage string) {
	log.Printf("invalid input parameters: %v", errorMessage)
	fmt.Println("Usage: poly [OPTIONS] -o output")
	flag.PrintDefaults()
	os.Exit(1)
}

func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	im, _, err := image.Decode(file)
	return im, err
}

func imageToRGBA(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
	return dst
}

func CopyRGBA(source *image.RGBA) *image.RGBA {
	copy := image.NewRGBA(source.Bounds())
	for i := range source.Pix {
		copy.Pix[i] = source.Pix[i]
	}
	copy.Stride = source.Stride
	return copy
}

func MinMax(values []int) (int, int) {
	min, max := values[0], values[0]
	for i := 0; i < len(values); i++ {
		if values[i] > max {
			max = values[i]
		}
		if values[i] < min {
			min = values[i]
		}
	}
	return min, max
}

// MinMaxPoints panics if length of the input array is 0
func MinMaxPoints(points []Point) (int, int, int, int) {
	xmin, xmax, ymin, ymax := points[0].X, points[0].X, points[0].Y, points[0].Y
	for _, point := range points {
		if point.X > xmax {
			xmax = point.X
		}
		if point.X < xmin {
			xmin = point.X
		}
		if point.Y > ymax {
			ymax = point.Y
		}
		if point.Y < ymin {
			ymin = point.Y
		}
	}
	return xmin, xmax, ymin, ymax
}

func SaveFile(path, contents string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(contents)
	return err
}

// TODO: Could do it faster assuming a, b > 0
func Pow(a, b int) int {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}

func absoluteDifferenceInt8(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
