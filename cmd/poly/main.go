package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"path/filepath"

	"os/exec"

	"github.com/dmartzol/poly/poly"
	"github.com/nfnt/resize"
)

var (
	inputPath    string
	Outputs      flagArray
	polygonCount int
	iterations   int
	maxSize      int
)

type flagArray []string

func (i *flagArray) String() string {
	return strings.Join(*i, ", ")
}

func (i *flagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	flag.StringVar(&inputPath, "i", "", "input image path")
	flag.Var(&Outputs, "o", "output image path")
	flag.IntVar(&polygonCount, "p", 50, "number of polygons")
	flag.IntVar(&iterations, "n", 1000, "number of iterations")
	flag.IntVar(&maxSize, "r", 256, "resize large input images to this size")
}

func main() {
	flag.Parse()
	ok := true
	if len(inputPath) == 0 {
		ok = poly.RaiseError("ERROR: input argument required")
	}
	if len(Outputs) == 0 {
		ok = poly.RaiseError("ERROR: output argument required")
	}
	if polygonCount <= 0 {
		ok = poly.RaiseError("ERROR: number of polygons should be > 0")
	}
	if iterations <= 0 {
		ok = poly.RaiseError("ERROR: number of iterations should be > 0")
	}
	if !ok {
		fmt.Println("Usage: poly [OPTIONS] -o output")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if Outputs[0] == inputPath {
		ok = poly.RaiseError("ERROR: input and output are the same file")
	}

	// Seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	inputImage, err := poly.LoadImage(inputPath)
	poly.CheckError(err)

	// scale down input image if needed
	size := uint(maxSize)
	if size > 0 {
		inputImage = resize.Thumbnail(size, size, inputImage, resize.Bilinear)
	}

	// Main block
	model := poly.NewModel(inputImage, polygonCount)
	fmt.Println(time.Now())
	start := time.Now()
	ratioMutations := model.Optimize(iterations)
	elapsed := time.Since(start)
	fmt.Printf("Mutations: %d\n", iterations)
	fmt.Printf("Took %d minutes\n", int(elapsed.Minutes()))
	speed := int(float64(iterations) * float64(polygonCount) / elapsed.Seconds())
	fmt.Println(speed, "polygons/s")
	fmt.Println(ratioMutations)
	fmt.Println("-------------------")

	// saving output
	for _, output := range Outputs {
		path := output
		extension := strings.ToLower(filepath.Ext(output))
		switch extension {
		default:
			poly.CheckError(fmt.Errorf("unrecognized file extension: %s", extension))
		case ".svg":
			poly.CheckError(poly.SaveFile(path, model.SVG()))
			app := "inkscape"
			arg0 := output
			arg1 := "--export-png=F.png"
			cmd := exec.Command(app, arg0, arg1)
			_, err := cmd.Output()
			poly.CheckError(err)
			// print(string(stdout))
		case ".png":
			model.PNG(output)
		}
	}
}
