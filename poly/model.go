package poly

import (
	"encoding/gob"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Model struct {
	Width, Height           int
	TargetImage             *image.RGBA
	NumPolygons             int
	Polygons                Polygons
	Scale                   float64
	Score                   float64
	Iteration               int
	BackgroundColor         Color
	MutateVertexProbability float64
}

func NewModel(input image.Image, numPolygons int, seed int64, bgColor Color) *Model {
	rand.Seed(seed)
	bounds := input.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	m := Model{
		Width:           w,
		Height:          h,
		Scale:           1.0,
		TargetImage:     imageToRGBA(input),
		NumPolygons:     numPolygons,
		BackgroundColor: bgColor,
	}

	for i := 0; i < numPolygons; i++ {
		order := rand.Intn(3) + 3
		polygon := newRandomPolygon(order, m.Width, m.Height)
		m.Polygons = append(m.Polygons, polygon)
	}

	rgbaCandidate := polygonsToRGBA(m.Polygons, m.BackgroundColor, m.Width, m.Height)
	m.Score = mse(m.TargetImage, rgbaCandidate)

	return &m
}

func (m *Model) mutate() Polygons {
	polygons := m.Polygons.clone()
	randomIndex := rand.Intn(m.NumPolygons)
	r := rand.Float64()
	polygons[randomIndex].mutate(r, m.Width, m.Height)
	return polygons
}

func (m *Model) step(id int, jobs <-chan int, results chan<- float64) {
	for j := range jobs {
		polygons := m.mutate()
		rgbaCandidate := polygonsToRGBA(polygons, m.BackgroundColor, m.Width, m.Height)

		newScore := mse(m.TargetImage, rgbaCandidate)
		if newScore < m.Score {
			m.Polygons = polygons
			m.Score = newScore
			m.Iteration = j
			results <- newScore
		}
		results <- math.Inf(1)
	}
}

func (m *Model) Optimize(iterations, concurrency, logFrequency int) float64 {
	var score float64
	var successful int

	jobs := make(chan int, iterations)
	results := make(chan float64, iterations)

	for w := 1; w <= concurrency; w++ {
		go m.step(w, jobs, results)
	}

	for i := 1; i <= iterations; i++ {
		jobs <- i
	}
	close(jobs)

	for a := 1; a <= iterations; a++ {
		current := <-results
		if current < score || a == 1 {
			score = current
		}
		if current != 0.0 {
			successful++
		}
		if a%logFrequency == 0 {
			fmt.Printf("%v,%v,%v\n", time.Now().Format(time.RFC3339), a, score)
		}
	}

	log.Printf("successful iterations: %v", successful)

	return score
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

	for _, polygon := range polygons {
		rasterizePolygonWWN(polygon, rgba)
	}

	return rgba
}

func (m *Model) GOB(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(m)
	if err != nil {
		return fmt.Errorf("unable to encode file: %w", err)
	}

	return nil
}

func ReadGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(object)
	if err != nil {
		return fmt.Errorf("unable to decode file: %w", err)
	}

	return nil
}

func (m *Model) SVG() string {
	bg := m.BackgroundColor
	var lines []string
	lines = append(lines, fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" version=\"1.1\" width=\"%d\" height=\"%d\">", 2*m.Width, 2*m.Height))
	lines = append(lines, fmt.Sprintf("<rect x=\"0\" y=\"0\" width=\"%d\" height=\"%d\" fill=\"#%02x%02x%02x\" />", 2*m.Width, 2*m.Height, bg.R, bg.G, bg.B))
	lines = append(lines, fmt.Sprintf("<g transform=\"scale(%f) translate(0.5 0.5)\">", 2*m.Scale))
	for _, polygon := range m.Polygons {
		color := polygon.Color
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

func (m *Model) PNG(fname string) error {
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
