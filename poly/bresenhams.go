package poly

import "image"

// DrawLine is the Bresenham's algorithm for painting lines efficiently
func DrawLine(x0, y0, x1, y1 int, color Color, canvas *image.RGBA) {
	var cx = x0
	var cy = y0

	var dx = x1 - cx
	var dy = y1 - cy
	if dx < 0 {
		dx = 0 - dx
	}
	if dy < 0 {
		dy = 0 - dy
	}

	var sx int
	var sy int
	if cx < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if cy < y1 {
		sy = 1
	} else {
		sy = -1
	}
	var err = dx - dy

	var n int
	for n = 0; n < 1000; n++ {
		DrawPoint(cx, cy, color, canvas)
		if (cx == x1) && (cy == y1) {
			return
		}
		var e2 = 2 * err
		if e2 > (0 - dy) {
			err = err - dy
			cx = cx + sx
		}
		if e2 < dx {
			err = err + dx
			cy = cy + sy
		}
	}
}
