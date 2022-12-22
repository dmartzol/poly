package poly

import (
	"math/rand"
)

type Polygon struct {
	Color Color
	// Vertices represents the list of coordinates for the vertices of the polygon
	Vertices  []Point
	Points    []Point
	HasPoints bool
}

func (p Polygon) clone() Polygon {
	var polygon Polygon
	polygon.Vertices = make([]Point, len(p.Vertices))
	polygon.Points = make([]Point, len(p.Points))
	copy(polygon.Vertices, p.Vertices)
	copy(polygon.Points, p.Points)
	polygon.HasPoints = p.HasPoints
	polygon.Color = p.Color
	return polygon
}

type Polygons []Polygon

func (ps Polygons) clone() Polygons {
	var polygons Polygons
	for _, p := range ps {
		c := p.clone()
		polygons = append(polygons, c)
	}
	return polygons
}

func newRandomPolygon(order, maxX, maxY int) Polygon {
	points := newRandomVertices(order, maxX, maxY)
	polygon := Polygon{}
	polygon.Vertices = points
	polygon.Color = newRandomColor()
	polygon.HasPoints = false
	return polygon
}

func (polygon *Polygon) mutate(ratio float64, w, h int) {
	randomFloat := rand.Float64()
	if randomFloat > 0.5 {
		polygon.mutateVertex(ratio, w, h)
	} else {
		polygon.Color = NewRandomColor()
		//polygon.mutateColor()
	}

	return
}

func (polygon *Polygon) mutateColor() {
	amplitude := 50
	channel := rand.Intn(4)
	// displacement is an integer in the interval [-50, 100)
	displacement := rand.Intn(2*amplitude) - amplitude
	colorList := [4]uint8{
		polygon.Color.R,
		polygon.Color.G,
		polygon.Color.B,
		polygon.Color.A,
	}
	colorList[channel] = clamp(int(colorList[channel]), displacement, 0, 255)
	polygon.Color =
		Color{
			colorList[0],
			colorList[1],
			colorList[2],
			colorList[3],
		}
}

func (polygon *Polygon) mutateVertex(ratio float64, width, height int) {
	polygon.HasPoints = false
	randomVertexIndex := rand.Intn(len(polygon.Vertices))

	var p [2]int
	if ratio < 0.07 {
		p = [2]int{
			polygon.Vertices[randomVertexIndex].X,
			polygon.Vertices[randomVertexIndex].Y,
		}
		amplitude := 10
		displacement := rand.Intn(2*amplitude+1) - amplitude
		direction := rand.Intn(2)
		maxValue := [2]int{width, height}
		p[direction] = int(clamp(p[direction], displacement, 0, maxValue[direction]))
	} else {
		p = [2]int{rand.Intn(width), rand.Intn(height)}
	}
	polygon.Vertices[randomVertexIndex] = Point{p[0], p[1]}
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
