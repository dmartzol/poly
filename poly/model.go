package poly

import (
	"fmt"
	"image"
	"math/rand"
	"strconv"
	"strings"
)

type Model struct {
	Width, Height int
	Target        *image.RGBA
	Current       *Individual
	NumPolygons   int
	Scale         float64
}

func NewModel(input image.Image, numPolygons int, seed int64, bgColor Color) *Model {
	rand.Seed(seed)
	bounds := input.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	model := Model{
		Width:       w,
		Height:      h,
		Scale:       1.0,
		Target:      imageToRGBA(input),
		NumPolygons: numPolygons,
	}

	model.Current = NewIndividual(model.Target, model.NumPolygons, bgColor)

	return &model
}

func (model *Model) Optimize(n int) []float64 {
	// Mutating the individual
	ratioMutations := []float64{100}
	var numMutations, succesfulMutations int

	for i := 0; i < n; i++ {
		numMutations++
		lastRatio := ratioMutations[len(ratioMutations)-1]
		if numMutations == 100 {
			ratioMutations =
				append(
					ratioMutations,
					float64(succesfulMutations)/float64(numMutations),
				)
			succesfulMutations = 0
			numMutations = 0
		}
		model.Current.mutate(lastRatio)
		if model.Current.compare() {
			succesfulMutations++
		}
	}
	return ratioMutations
}

func (model *Model) SVG() string {
	bg := model.Current.BackgroundColor
	var lines []string
	lines = append(lines, fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" version=\"1.1\" width=\"%d\" height=\"%d\">", 2*model.Width, 2*model.Height))
	lines = append(lines, fmt.Sprintf("<rect x=\"0\" y=\"0\" width=\"%d\" height=\"%d\" fill=\"#%02x%02x%02x\" />", 2*model.Width, 2*model.Height, bg.R, bg.G, bg.B))
	lines = append(lines, fmt.Sprintf("<g transform=\"scale(%f) translate(0.5 0.5)\">", 2*model.Scale))
	for _, polygon := range model.Current.Polygons {
		color := polygon.ColorRGBA
		attrs := "<polygon fill=\"#%02x%02x%02x\" fill-opacity=\"%f\""
		attrs = fmt.Sprintf(attrs, color.R, color.G, color.B, float64(color.A)/255)
		p := " points=\""
		for _, vertex := range polygon.Vertices {
			p = p + strconv.Itoa(vertex.X) + "," + strconv.Itoa(vertex.Y) + " "
		}
		p = p + "\"" + "/>"
		attrs = attrs + p
		lines = append(lines, attrs)
	}
	lines = append(lines, "</g>")
	lines = append(lines, "</svg>")
	return strings.Join(lines, "\n")
}

func (model *Model) PNG(fname string) {
	model.Current.decode(fname)
}
