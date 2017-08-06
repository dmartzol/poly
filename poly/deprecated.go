package poly

func addColorsWithTransparentBg(foreground, background [4]uint8) [4]uint8 {
	alphaF := float64(foreground[3]) / 256
	alphaB := float64(background[3]) / 256
	alphaR := alphaF + alphaB*(1-alphaF)
	R := (float64(foreground[0])*alphaF + float64(background[0])*alphaB*(1-alphaF)) / alphaR
	G := (float64(foreground[1])*alphaF + float64(background[1])*alphaB*(1-alphaF)) / alphaR
	B := (float64(foreground[2])*alphaF + float64(background[2])*alphaB*(1-alphaF)) / alphaR
	alphaR = alphaR * 256
	return [4]uint8{uint8(R), uint8(G), uint8(B), uint8(alphaR)}
}
