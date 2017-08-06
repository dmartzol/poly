package poly

import (
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
)

type Individual struct {
	Width, Height       int
	NumPolygons         int
	RGBA                *image.RGBA
	Target              *image.RGBA
	Polygons            []Polygon
	BackgroundColor     Color
	Score, OldScore     float64
	ChoosenPolygonIndex int
	PolygonBackup       Polygon
	RGBABackup          *image.RGBA
}

// TODO: bgColor function to be executed in NewIndividual
func NewIndividual(target *image.RGBA, numPolygons int, bgColor Color) *Individual {
	individual := &Individual{}
	individual.Width = target.Bounds().Max.X
	individual.Height = target.Bounds().Max.Y
	individual.NumPolygons = numPolygons
	individual.Target = CopyRGBA(target)
	individual.BackgroundColor = bgColor
	for i := 0; i < numPolygons; i++ {
		order := rand.Intn(3) + 3
		polygon := NewRandomPolygon(order, individual.Width, individual.Height)
		individual.Polygons = append(individual.Polygons, polygon)
	}
	individual.RGBA = image.NewRGBA(target.Bounds())
	individual.Score = individual.MSE()
	return individual
}

func (individual *Individual) drawPolygons() {
	var rgba *image.RGBA
	rect := image.Rect(0, 0, individual.Width, individual.Height)
	rgba = image.NewRGBA(rect)
	backgroundColor(individual.BackgroundColor, rgba)
	// rgba := individual.Polygons[individual.ChoosenPolygonIndex].subImage
	// TODO: Implement calculating average to decide white or black Bg
	// First subImage must be the Bg
	// First subImage should be printed in NewIndividual
	// setBackgroundColor(individual.BackgroundColor, rgba)
	for _, polygon := range individual.Polygons {
		RasterizePolygonWWN(polygon, rgba)
	}
	individual.RGBA = rgba
	individual.Score = 0
}

func (individual *Individual) decode(fname string) {
	file, _ := os.Create(fname)
	png.Encode(file, individual.RGBA)
}

func (individual *Individual) mutate(ratio float64) {
	// Select a random polygon
	individual.ChoosenPolygonIndex = rand.Intn(individual.NumPolygons)

	// Backing up elements
	individual.backup()

	//Modifying the individual
	individual.Score = 0
	individual.Polygons[individual.ChoosenPolygonIndex].mutate(ratio, individual.Width, individual.Height)
}

// compare returns true if the mutation has been succesful and
// restores the backup otherwise
func (individual *Individual) compare() bool {
	if individual.OldScore <= individual.MSE() {
		individual.restoreBackup()
		return false
	}
	return true
}

func (individual *Individual) backup() {
	chosen := individual.ChoosenPolygonIndex
	individual.PolygonBackup = individual.Polygons[chosen].clone()
	individual.RGBABackup = CopyRGBA(individual.RGBA)
	individual.OldScore = individual.MSE()
}

func (individual *Individual) restoreBackup() {
	chosen := individual.ChoosenPolygonIndex
	individual.Polygons[chosen] = individual.PolygonBackup
	individual.RGBA = CopyRGBA(individual.RGBABackup)
	individual.Score = individual.OldScore
}

func (individual *Individual) MSE() float64 {
	if individual.Score != 0 {
		return individual.Score
	}
	individual.drawPolygons()
	targetPixels := individual.Target.Pix
	w, h := individual.RGBA.Bounds().Max.X, individual.RGBA.Bounds().Max.Y
	size := w * h * 4
	sum := 0
	for i := 0; i < size; i++ {
		if i%3 != 0 {
			d := absoluteDifferenceInt8(targetPixels[i], individual.RGBA.Pix[i])
			// TODO: Write a faster Pow func
			sum = sum + Pow(d, 2)
		}
	}
	root := float64(sum)
	root = math.Sqrt(root)
	individual.Score = root
	return root
}

func backgroundColor(c Color, canvas *image.RGBA) {
	l := len(canvas.Pix)
	for i := 0; i < l; i += 4 {
		canvas.Pix[i] = c.R
		canvas.Pix[i+1] = c.G
		canvas.Pix[i+2] = c.B
		canvas.Pix[i+3] = c.A
	}
}
