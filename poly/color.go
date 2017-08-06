package poly

type Color struct {
	R, G, B, A uint8
}

func addColors(foreground, background Color) Color {
	var r [4]uint16
	alpha := uint16(foreground.A) + 1
	inverseAlpha := 256 - alpha
	r[0] = ((alpha*uint16(foreground.R) + inverseAlpha*uint16(background.R)) >> 8)
	r[1] = ((alpha*uint16(foreground.G) + inverseAlpha*uint16(background.G)) >> 8)
	r[2] = ((alpha*uint16(foreground.B) + inverseAlpha*uint16(background.B)) >> 8)
	r[3] = 0xff
	return Color{uint8(r[0]), uint8(r[1]), uint8(r[2]), uint8(r[3])}
}

func subtractColors(c, background Color) Color {
	var old [4]uint16
	alpha := uint16(c.A) + 1
	inverseAlpha := 256 - alpha
	old[0] = (256*uint16(background.R) - alpha*uint16(c.R)) / inverseAlpha
	old[1] = (256*uint16(background.G) - alpha*uint16(c.G)) / inverseAlpha
	old[2] = (256*uint16(background.B) - alpha*uint16(c.B)) / inverseAlpha
	old[3] = 0xff
	return Color{uint8(old[0]), uint8(old[1]), uint8(old[2]), uint8(old[3])}
}
