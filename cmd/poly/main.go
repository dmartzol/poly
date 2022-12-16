package main

import (
	"flag"
	"fmt"
	"log"
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
	maxImageSize int
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
	flag.IntVar(&maxImageSize, "r", 256, "resize large input images to this size")
}

func main() {
	flag.Parse()
	// flag validation
	if len(inputPath) == 0 {
		poly.PrintDefaultsWithError("input argument required")
	}
	if len(Outputs) == 0 {
		poly.PrintDefaultsWithError("output argument required")
	}
	if Outputs[0] == inputPath {
		poly.PrintDefaultsWithError("input and output are the same file")
	}
	if polygonCount <= 0 {
		poly.PrintDefaultsWithError("number of polygons should be > 0")
	}
	if iterations <= 0 {
		poly.PrintDefaultsWithError("number of iterations should be > 0")
	}

	inputImage, err := poly.LoadImage(inputPath)
	if err != nil {
		log.Printf("unable to load image: %v", err)
		return
	}

	// scale down input image if needed
	size := uint(maxImageSize)
	if size > 0 {
		inputImage = resize.Thumbnail(size, size, inputImage, resize.Bilinear)
	}

	// Main block
	whiteColor := poly.Color{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
	randomSeed := time.Now().UTC().UnixNano()
	model := poly.NewModel(inputImage, polygonCount, randomSeed, whiteColor)
	start := time.Now()
	ratioMutations := model.Optimize(iterations)
	elapsed := time.Since(start)

	// logging info
	fmt.Printf("Mutations: %d\n", iterations)
	fmt.Printf("took %v\n", elapsed)
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
