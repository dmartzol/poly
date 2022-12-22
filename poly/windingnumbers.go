package poly

import "image"

type Point struct {
	X, Y int
}

func (point *Point) clone() Point {
	return Point{point.X, point.Y}
}

// WindingNumber is the winding number test for a point in a polygon
//      Input:   P = a point,
//               V[] = vertex points of a polygon V[n+1] with V[n]=V[0]
//      Return:  wn = the winding number (=0 only when P is outside)
func windingNumber(P Point, V []Point) int {
	wn := 0
	lim := len(V) - 1
	for i := 0; i < lim; i++ {
		if V[i].Y <= P.Y {
			if V[i+1].Y > P.Y {
				if isLeft(V[i], V[i+1], P) > 0 {
					wn++
				}
			}
		} else {
			if V[i+1].Y <= P.Y {
				if isLeft(V[i], V[i+1], P) < 0 {
					wn--
				}
			}
		}
	}
	if V[lim].Y <= P.Y {
		if V[0].Y > P.Y {
			if isLeft(V[lim], V[0], P) > 0 {
				wn++
			}
		}
	} else {
		if V[0].Y <= P.Y {
			if isLeft(V[lim], V[0], P) < 0 {
				wn--
			}
		}
	}
	return wn
}

// isLeft(): tests if a point is Left|On|Right of an infinite line.
//    Input:  three points P0, P1, and P2
//    Return: >0 for P2 left of the line through P0 and P1
//            =0 for P2  on the line
//            <0 for P2  right of the line
//    See: Algorithm 1 "Area of Triangles and Polygons"
func isLeft(P0, P1, P2 Point) int {
	return (P1.X-P0.X)*(P2.Y-P0.Y) - (P2.X-P0.X)*(P1.Y-P0.Y)
}

func rasterizePolygonWWN(polygon Polygon, result *image.RGBA) {
	minX, maxX, minY, maxY := minMaxPoints(polygon.Vertices)
	// copy(polygon.subImage.Pix, result.Pix)
	if polygon.HasPoints {
		for _, point := range polygon.Points {
			drawPoint(point.X, point.Y, polygon.Color, result)
		}
	} else {
		polygon.HasPoints = true
		for x := minX; x <= maxX; x++ {
			for y := minY; y <= maxY; y++ {
				if windingNumber(Point{x, y}, polygon.Vertices) != 0 {
					drawPoint(x, y, polygon.Color, result)
					polygon.Points = append(polygon.Points, Point{x, y})
				}
			}
		}
	}
}

func subtractPolygonWWN(polygon *Polygon, result *image.RGBA) {
	minX, maxX, minY, maxY := minMaxPoints(polygon.Vertices)
	var point Point
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			point = Point{x, y}
			if windingNumber(point, polygon.Vertices) != 0 {
				subtractPoint(x, y, polygon.Color, result)
			}
		}
	}
}
