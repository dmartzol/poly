package poly

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type AdvancedModel struct {
	Width, Height           int
	TargetImage             *image.RGBA
	NumPolygons             int
	Polygons                Polygons
	Scale                   float64
	Score                   float64
	BackgroundColor         Color
	MutateVertexProbability float64
}

func NewAdvancedModel(input image.Image, numPolygons int, seed int64, bgColor Color) *AdvancedModel {
	rand.Seed(seed)
	bounds := input.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	m := AdvancedModel{
		Width:           w,
		Height:          h,
		Scale:           1.0,
		TargetImage:     imageToRGBA(input),
		NumPolygons:     numPolygons,
		BackgroundColor: bgColor,
	}

	for i := 0; i < numPolygons; i++ {
		order := rand.Intn(3) + 3
		polygon := NewRandomPolygon(order, m.Width, m.Height)
		m.Polygons = append(m.Polygons, polygon)
	}

	rgbaCandidate := polygonsToRGBA(m.Polygons, m.BackgroundColor, m.Width, m.Height)
	m.Score = MSE(m.TargetImage, rgbaCandidate)

	return &m
}

type Model struct {
	Width, Height int
	TargetImage   *image.RGBA
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
		TargetImage: imageToRGBA(input),
		NumPolygons: numPolygons,
	}

	model.Current = NewIndividual(model.TargetImage, model.NumPolygons, bgColor)

	return &model
}

func (m *AdvancedModel) Optimize(iterations int) []float64 {
	var scores []float64

	for i := 0; i < iterations; i++ {
		polygons := make(Polygons, m.NumPolygons)
		copy(polygons, m.Polygons)
		randomIndex := rand.Intn(m.NumPolygons)
		polygons[randomIndex].mutate(0.15, m.Width, m.Height)
		rgbaCandidate := polygonsToRGBA(polygons, m.BackgroundColor, m.Width, m.Height)

		newScore := MSE(m.TargetImage, rgbaCandidate)
		if newScore < m.Score {
			m.Polygons = polygons
			m.Score = newScore
			scores = append(scores, newScore)
		}
	}

	return scores
}

func polygonsToRGBA(polygons Polygons, bgColor Color, w, h int) *image.RGBA {
	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)

	// first paints branckground with the selected color
	l := len(rgba.Pix)
	for i := 0; i < l; i += 4 {
		rgba.Pix[i] = bgColor.R
		rgba.Pix[i+1] = bgColor.G
		rgba.Pix[i+2] = bgColor.B
		rgba.Pix[i+3] = bgColor.A
	}

	// rgba := individual.Polygons[individual.ChoosenPolygonIndex].subImage
	// TODO: Implement calculating average to decide white or black Bg
	// First subImage must be the Bg
	// First subImage should be printed in NewIndividual
	// setBackgroundColor(individual.BackgroundColor, rgba)

	for _, polygon := range polygons {
		RasterizePolygonWWN(polygon, rgba)
	}

	return rgba
}

func MSE(target, candidate *image.RGBA) float64 {
	targetPixels := target.Pix
	w, h := candidate.Bounds().Max.X, candidate.Bounds().Max.Y
	size := w * h * 4
	sum := 0
	for i := 0; i < size; i++ {
		if i%3 != 0 { // avoiding calculating difference for transparency pixels
			d := absoluteDifferenceInt8(targetPixels[i], candidate.Pix[i])
			// TODO: Write a faster Pow func
			sum = sum + Pow(d, 2)
		}
	}

	return math.Sqrt(float64(sum))
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

func (m *AdvancedModel) SVG() string {
	bg := m.BackgroundColor
	var lines []string
	lines = append(lines, fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" version=\"1.1\" width=\"%d\" height=\"%d\">", 2*m.Width, 2*m.Height))
	lines = append(lines, fmt.Sprintf("<rect x=\"0\" y=\"0\" width=\"%d\" height=\"%d\" fill=\"#%02x%02x%02x\" />", 2*m.Width, 2*m.Height, bg.R, bg.G, bg.B))
	lines = append(lines, fmt.Sprintf("<g transform=\"scale(%f) translate(0.5 0.5)\">", 2*m.Scale))
	for _, polygon := range m.Polygons {
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

func (m *AdvancedModel) PNG(fname string) error {
	file, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}

	rgbaImage := polygonsToRGBA(m.Polygons, m.BackgroundColor, m.Width, m.Height)

	err = png.Encode(file, rgbaImage)
	if err != nil {
		return fmt.Errorf("unable to encode png: %w", err)
	}

	return nil
}
