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
	Input     string
	Outputs   flagArray
	Pols      int
	Iter      int
	InputSize int
)

type flagArray []string

// TODO: dani

func (i *flagArray) String() string {
	return strings.Join(*i, ", ")
}

func (i *flagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	flag.StringVar(&Input, "i", "", "input image path")
	flag.Var(&Outputs, "o", "output image path")
	flag.IntVar(&Pols, "p", 50, "number of polygons")
	flag.IntVar(&Iter, "n", 1000, "number of iterations")
	flag.IntVar(&InputSize, "r", 256, "resize large input images to this size")
}

func main() {
	flag.Parse()
	ok := true
	if len(Input) == 0 {
		ok = poly.RaiseError("ERROR: input argument required")
	}
	if len(Outputs) == 0 {
		ok = poly.RaiseError("ERROR: output argument required")
	}
	if Pols <= 0 {
		ok = poly.RaiseError("ERROR: number of polygons should be > 0")
	}
	if Iter <= 0 {
		ok = poly.RaiseError("ERROR: number of iterations should be > 0")
	}
	if !ok {
		fmt.Println("Usage: polygonal [OPTIONS] -o output")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if Outputs[0] == Input {
		ok = poly.RaiseError("ERROR: input and output are the same file")
	}

	// Seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	input, err := poly.LoadImage(Input)
	poly.Check(err)

	// scale down input image if needed
	size := uint(InputSize)
	if size > 0 {
		input = resize.Thumbnail(size, size, input, resize.Bilinear)
	}

	// Main block
	model := poly.NewModel(input, Pols)
	fmt.Println(time.Now())
	start := time.Now()
	ratioMutations := model.Optimize(Iter)
	elapsed := time.Since(start)
	fmt.Printf("Mutations: %d\n", Iter)
	fmt.Printf("Took %d minutes\n", int(elapsed.Minutes()))
	speed := int(float64(Iter) * float64(Pols) / elapsed.Seconds())
	fmt.Println(speed, "polygons/s")
	fmt.Println(ratioMutations)
	fmt.Println("-------------------")

	for _, output := range Outputs {
		path := output
		extension := strings.ToLower(filepath.Ext(output))
		switch extension {
		default:
			poly.Check(fmt.Errorf("unrecognized file extension: %s", extension))
		case ".svg":
			poly.Check(poly.SaveFile(path, model.SVG()))
			app := "inkscape"
			arg0 := output
			arg1 := "--export-png=F.png"
			cmd := exec.Command(app, arg0, arg1)
			_, err := cmd.Output()
			poly.Check(err)
			// print(string(stdout))
		case ".png":
			model.PNG(output)
		}
	}
}
