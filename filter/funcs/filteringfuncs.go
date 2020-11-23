package funcs

import (
	"fmt"
	im "image"
	"image/color"
)

func InvertColor(read *im.Image, x int, y int) *color.RGBA64 {
	rr, gg, bb, aa := (*read).At(x, y).RGBA()

	r := aa - rr
	g := aa - gg
	b := aa - bb
	return &color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(aa)}
}
func FindEdge(read *im.Image, x int, y int, size int, threshold float64) *color.RGBA64 {

	return &color.RGBA64{}
}
func Nullify(read *im.Image, x int, y int) *color.RGBA64 {
	return &color.RGBA64{}
}
func Identity(read *im.Image, x, y int) *color.RGBA64 {
	rr, gg, bb, aa := (*read).At(x, y).RGBA()
	return &color.RGBA64{R: uint16(rr), G: uint16(gg), B: uint16(bb), A: uint16(aa)}
}
func Print(read *im.Image, x, y int) *color.RGBA64 {
	rr, gg, bb, aa := (*read).At(x, y).RGBA()

	if aa != 0 {
		rr = rr * 65535 / aa
		gg = gg * 65535 / aa
		bb = bb * 65535 / aa
	}

	fmt.Printf("%v %v %v %v\n", rr/255, gg/255, bb/255, aa/255)

	if aa != 0 {
		rr = rr * aa / 65535
		gg = gg * aa / 65535
		bb = bb * aa / 65535
	}

	return &color.RGBA64{R: uint16(rr), G: uint16(gg), B: uint16(bb), A: uint16(aa)}
}
