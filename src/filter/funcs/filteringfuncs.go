package funcs

import (
	im "image"
	"image/color"
)

func InvertColor(read im.Image, x int, y int) color.RGBA64 {
	rr, gg, bb, aa := read.At(x, y).RGBA()
	const max = 65535
	r := (max - rr) * aa / max
	g := (max - gg) * aa / max
	b := (max - bb) * aa / max
	return color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(aa)}
}
func FindEdge(read im.Image, x int, y int, size int, threshold float64) color.RGBA64 {

	return color.RGBA64{}
}
func Nullify(read im.Image, x int, y int) color.RGBA64 {
	return color.RGBA64{}
}
