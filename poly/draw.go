package poly

import "image"

func drawPoint(cx, cy int, color Color, canvas *image.RGBA) {
	p := canvas.PixOffset(cx, cy)
	var oldColor Color
	// Not using loops for better performance
	oldColor.R = canvas.Pix[p]
	oldColor.G = canvas.Pix[p+1]
	oldColor.B = canvas.Pix[p+2]
	oldColor.A = canvas.Pix[p+3]
	newColor := addColors(color, oldColor)
	// Not using loops for better performance
	canvas.Pix[p] = newColor.R
	canvas.Pix[p+1] = newColor.G
	canvas.Pix[p+2] = newColor.B
	canvas.Pix[p+3] = newColor.A
}

func subtractPoint(cx, cy int, color Color, canvas *image.RGBA) {
	p := canvas.PixOffset(cx, cy)
	var oldColor Color
	// Not using loops for better performance
	oldColor.R = canvas.Pix[p]
	oldColor.G = canvas.Pix[p+1]
	oldColor.B = canvas.Pix[p+2]
	oldColor.A = canvas.Pix[p+3]
	newColor := subtractColors(color, oldColor)
	// Not using loops for better performance
	canvas.Pix[p] = newColor.R
	canvas.Pix[p+1] = newColor.G
	canvas.Pix[p+2] = newColor.B
	canvas.Pix[p+3] = newColor.A
}

func fillRow(y, x0, x1 int, color Color, canvas *image.RGBA) {
	for x := x0; x < x1; x++ {
		drawPoint(x, y, color, canvas)
	}
}

func pixelIsNotClear(x, y int, canvas *image.RGBA) bool {
	if canvas.Pix[canvas.PixOffset(x, y)+3] != 0 {
		return true
	}
	return false
}

func clamp(a, b, min, max int) uint8 {
	if a+b < min {
		return uint8(min)
	}
	if a+b > max {
		return uint8(max)
	}
	return uint8(a + b)
}
