package funcs

import (
	"github.com/mxyns/go-filter/filter"
	im "image"
	"image/color"
	"strconv"
)

//Get the value of the radius
func ParseNoiseReductionArgs(_ *filter.Filter, args *map[string]interface{}) *map[string]interface{} {

	// ["radius"]
	if arg := (*args)["radius"]; arg != nil {
		(*args)["radius"], _ = strconv.Atoi(arg.(string))
	} else {
		(*args)["radius"] = 1
	}

	return args
}

//Calculate for each pixel the noise reduction applied according to radius
func NoiseReduction(read *im.Image, x int, y int, args *map[string]interface{}) *color.RGBA64 {

	var r = ((*args)["radius"]).(int)

	xMin := maxInt(x-r, (*read).Bounds().Min.X)
	yMin := maxInt(y-r, (*read).Bounds().Min.Y)
	xMax := minInt(x+r, (*read).Bounds().Max.X-1)
	yMax := minInt(y+r, (*read).Bounds().Max.Y-1)

	c := uint32((1 + xMax - xMin) * (1 + yMax - yMin)) // pixel counted (avoid border pbs)

	var mR, mG, mB, mA uint32 = 0, 0, 0, 0 // mean
	for i := xMin; i <= xMax; i++ {
		for j := yMin; j <= yMax; j++ {
			tR, tG, tB, tA := (*read).At(i, j).RGBA()
			mR += tR
			mG += tG
			mB += tB
			mA += tA
		}
	}
	_, _, _, _ = (*read).At(x, y).RGBA()
	return &color.RGBA64{R: uint16(mR / c), G: uint16(mG / c), B: uint16(mB / c), A: uint16(mA / c)}
}
