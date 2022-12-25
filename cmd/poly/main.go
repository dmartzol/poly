package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"path/filepath"

	"github.com/dmartzol/poly/poly"
	"github.com/nfnt/resize"
)

var (
	inputPath    string
	Outputs      flagArray
	polygonCount int
	iterations   int
	maxImageSize int
	concurrency  int
	cpuprofile   string
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
	flag.IntVar(&concurrency, "c", 3, "number of workers to use")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")
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
	if polygonCount <= 0 {
		poly.PrintDefaultsWithError("number of polygons should be > 0")
	}
	if iterations <= 0 {
		poly.PrintDefaultsWithError("number of iterations should be > 0")
	}

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Printf("unable to create profile: %v", err)
			return
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

	}

	model := new(poly.Model)

	extension := strings.ToLower(filepath.Ext(inputPath))
	if extension == ".jpg" {
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
		model = poly.NewModel(inputImage, polygonCount, randomSeed, whiteColor)

	} else if extension == ".gob" {
		err := poly.ReadGob(inputPath, model)
		if err != nil {
			log.Printf("unable to read gob file: %v", err)
			return
		}
	}

	start := time.Now()
	score := model.Optimize(iterations, concurrency)
	elapsed := time.Since(start)

	// logging info
	fmt.Printf("Mutations: %d\n", iterations)
	fmt.Printf("took %v\n", elapsed)
	speed := int(float64(iterations) * float64(polygonCount) / elapsed.Seconds())
	fmt.Println(speed, "polygons/s")
	fmt.Printf("score: %v", score)

	// saving output
	for _, output := range Outputs {
		path := output
		extension := strings.ToLower(filepath.Ext(output))
		switch extension {
		default:
			log.Printf("unrecognized file extension: %s", extension)
			return
		case ".svg":
			err := poly.SaveFile(path, model.SVG())
			if err != nil {
				log.Printf("unable to save SVG file: %v", err)
				return
			}
		case ".gob":
			err := model.GOB(output)
			if err != nil {
				log.Printf("unable to save GOB file: %v", err)
				return
			}
		case ".png":
			err := model.PNG(output)
			if err != nil {
				log.Printf("unable to save PNG file: %v", err)
				return
			}
		}
	}
}
