package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Scrypy/poly/poly"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

var (
	Input     string
	Output    string
	Nth       int
	InputSize int
)

func init() {
	flag.StringVar(&Input, "i", "", "input image path")
	flag.StringVar(&Output, "o", "", "output image path")
	flag.IntVar(&Nth, "n", 10, "number of polygons")
	flag.IntVar(&InputSize, "r", 256, "resize large input images to this size")
}

func raiseError(message string) bool {
	fmt.Fprintln(os.Stderr, message)
	return false
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	ok := true
	if len(Input) == 0 {
		ok = raiseError("ERROR: input argument required")
	}
	if len(Output) == 0 {
		ok = raiseError("ERROR: output argument required")
	}
	if Nth <= 0 {
		ok = raiseError("ERROR: number of polygons should be > 0")
	}
	if !ok {
		fmt.Println("Usage: polygonal [OPTIONS] -o output")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	input, err := poly.LoadImage(Input)
	check(err)

	// scale down input image if needed
	size := uint(InputSize)
	if size > 0 {
		input = resize.Thumbnail(size, size, input, resize.Bilinear)
	}
	bounds := input.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	var n = 50
	var verticesCount int
	var polygons [][][2]int
	var polygon [][2]int
	var vertex [2]int
	for i := 0; i < n; i++ {
		verticesCount = rand.Intn(2) + 3
		for j := 0; j < verticesCount; j++ {
			vertex[0] = rand.Intn(w)
			vertex[1] = rand.Intn(h)
			polygon = append(polygon, vertex)
		}
		polygons = append(polygons, polygon)
	}
	dc := gg.NewContext(w, h)
	for i := range polygons {
		printPolygon(dc, polygons[i])
	}
	dc.SavePNG(Output)
}

func printPolygon(dc *gg.Context, polygon [][2]int) image.Image {
	dc.NewSubPath()
	for i := range polygon {
		dc.LineTo(float64(polygon[i][0]), float64(polygon[i][1]))
	}
	dc.ClosePath()
	dc.SetRGBA(0, 0, 0, 0.05)
	dc.Fill()
	return dc.Image()
}

// Get the bi-dimensional pixel array
func getPixels(img image.Image) [][]Pixel {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}
	return pixels
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}
