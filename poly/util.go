package poly

import (
	"flag"
	"fmt"
	"image"
	"log"
	"math"

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

// minMaxPoints returns the boundaries of the smallest rectangle that contains all the given points
func minMaxPoints(points []Point) (int, int, int, int) {
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
func pow(a, b int) int {
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

func mse(target, candidate *image.RGBA) float64 {
	targetPixels := target.Pix
	w, h := candidate.Bounds().Max.X, candidate.Bounds().Max.Y
	size := w * h * 4
	sum := 0
	for i := 0; i < size; i++ {
		if i%3 != 0 { // avoiding calculating difference for transparency pixels
			d := absoluteDifferenceInt8(targetPixels[i], candidate.Pix[i])
			// TODO: Write a faster Pow func
			sum = sum + pow(d, 2)
		}
	}

	return math.Sqrt(float64(sum))
}

func absoluteDifferenceInt8(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}
