package funcs

import (
	im "image"
	"image/color"
)

func GrayScaleAverage(read *im.Image, x int, y int, _ *map[string]interface{}) *color.RGBA64 {
	rr, gg, bb, aa := (*read).At(x, y).RGBA()
	grey := uint16((rr + gg + bb) / 3)

	return &color.RGBA64{R: grey, G: grey, B: grey, A: uint16(aa)}
}

func GrayScaleLuminosity(read *im.Image, x int, y int, _ *map[string]interface{}) *color.RGBA64 {
	rr, gg, bb, aa := (*read).At(x, y).RGBA()
	grey := uint16((0.3 * float32(rr)) + (0.59 * float32(gg)) + (0.11 * float32(bb)))

	return &color.RGBA64{R: grey, G: grey, B: grey, A: uint16(aa)}
}

func GrayScaleDesaturation(read *im.Image, x int, y int, _ *map[string]interface{}) *color.RGBA64 {
	rr, gg, bb, aa := (*read).At(x, y).RGBA()

	grey := uint16((triMaxU32(rr, gg, bb) + triMinU32(rr, gg, bb)) / 2)
	return &color.RGBA64{R: grey, G: grey, B: grey, A: uint16(aa)}
}
