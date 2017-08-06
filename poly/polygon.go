package poly

import (
	"image"
	"math/rand"
)

type Polygon struct {
	// Order represents the number of vertices in the polygon
	Order     int
	ColorRGBA Color
	// Vertices represents the list of coordinates for the vertices of the polygon
	Vertices []Point
	// Rectangle represents the minimum rectangle that contains the polygon
	Rectangle  image.Rectangle
	subImage   *image.RGBA
	IsModified bool
	Points     []Point
	HasPoints  bool
}

func (polygon *Polygon) clone() Polygon {
	polygonCopied := &Polygon{}
	polygonCopied.Order = polygon.Order
	polygonCopied.ColorRGBA = polygon.ColorRGBA
	polygonCopied.Vertices = make([]Point, polygon.Order)
	for i := range polygon.Vertices {
		polygonCopied.Vertices[i] = polygon.Vertices[i].clone()
	}
	polygonCopied.Rectangle = polygon.Rectangle
	polygonCopied.subImage = image.NewRGBA(polygon.subImage.Bounds())
	copy(polygonCopied.subImage.Pix, polygon.subImage.Pix)
	polygonCopied.HasPoints = polygon.HasPoints
	return *polygonCopied
}

func NewRandomPolygon(order, maxX, maxY int) Polygon {
	points := newRandomVertices(order, maxX, maxY)
	x0, x1, y0, y1 := MinMaxPoints(points)
	polygon := Polygon{}
	polygon.Order = order
	polygon.Vertices = points
	polygon.Rectangle = image.Rect(x0, y0, x1, y1)
	polygon.ColorRGBA = newRandomColor()
	polygon.subImage = image.NewRGBA(image.Rect(0, 0, maxX, maxY))
	polygon.HasPoints = false
	return polygon
}

func (polygon *Polygon) mutate(ratio float64, w, h int) {
	randomFloat := rand.Float64()
	if randomFloat > 0.5 {
		polygon.HasPoints = false
		polygon.mutateVertex(ratio, w, h)
		return
	}
	polygon.mutateColor()
	return
}

func (polygon *Polygon) mutateColor() {
	amplitude := 50
	channel := rand.Intn(4)
	displacement := rand.Intn(2*amplitude) - amplitude
	colorList := [4]uint8{
		polygon.ColorRGBA.R,
		polygon.ColorRGBA.G,
		polygon.ColorRGBA.B,
		polygon.ColorRGBA.A,
	}
	colorList[channel] = clamp(int(colorList[channel]), displacement, 0, 255)
	polygon.ColorRGBA =
		Color{
			colorList[0],
			colorList[1],
			colorList[2],
			colorList[3],
		}
}

func (polygon *Polygon) mutateVertex(ratio float64, width, height int) {
	randomVertex := rand.Intn(len(polygon.Vertices))
	var p [2]int
	if ratio < 0.07 {
		p = [2]int{
			polygon.Vertices[randomVertex].X,
			polygon.Vertices[randomVertex].Y,
		}
		amplitude := 10
		displacement := rand.Intn(2*amplitude+1) - amplitude
		direction := rand.Intn(2)
		maxValue := [2]int{width, height}
		p[direction] = int(clamp(p[direction], displacement, 0, maxValue[direction]))
	} else {
		p = [2]int{rand.Intn(width), rand.Intn(height)}
	}
	polygon.Vertices[randomVertex] = Point{p[0], p[1]}
}
func newRandomColor() Color {
	var color [4]uint8
	for i := 0; i < 3; i++ {
		color[i] = uint8(rand.Intn(256))
	}
	color[3] = uint8(75)
	return Color{color[0], color[1], color[2], color[3]}
}

func newRandomVertices(order, maxX, maxY int) []Point {
	points := make([]Point, order)
	for i := 0; i < order; i++ {
		points[i] = Point{rand.Intn(maxX), rand.Intn(maxY)}
	}
	return points
}
