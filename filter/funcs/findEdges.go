package funcs

import (
	"github.com/mxyns/go-filter/filter"
	parse "github.com/mxyns/go-filter/io"
	im "image"
	"image/color"
	"math"
	"strconv"
	"strings"
)

var (
	BLACK             = &color.RGBA64{R: 0, G: 0, B: 0, A: math.MaxUint16}
	TRANSPARENT_WHITE = &color.RGBA64{R: math.MaxUint16, G: math.MaxUint16, B: math.MaxUint16, A: 0}
	distanceFunctions = make(map[string]func(a *color.Color, b *color.RGBA64) float64)
)

func init() {

	distanceFunctions["euclidean"] = ColorEuclideanDist
	distanceFunctions["norm1"] = ColorNorm1Dist
	distanceFunctions["inf"] = ColorNormInfDist
}

func ParseFindEdgesArgs(filter *filter.Filter, args *map[string]interface{}) *map[string]interface{} {

	argsMap := ParseNoiseReductionArgs(filter, args)

	if val := (*argsMap)["threshold"]; val != nil {
		(*argsMap)["threshold"], _ = strconv.ParseFloat((*argsMap)["threshold"].(string), 64)
	} else {
		(*argsMap)["threshold"] = 1000.0
	}

	if val := (*argsMap)["dist"]; val != nil {
		if distType := strings.ToLower((*argsMap)["dist"].(string)); distanceFunctions[distType] != nil {
			(*argsMap)["dist"] = distType
		}
	} else {
		(*argsMap)["dist"] = "euclidean"
	}

	if val := (*argsMap)["edge_color"]; val != nil {
		(*argsMap)["edge_color"] = parse.ParseColor((*argsMap)["edge_color"].(string))
	} else {
		(*argsMap)["edge_color"] = &parse.OptionalRGBA64{Color: BLACK, Empty: false}
	}

	if val := (*argsMap)["gap_color"]; val != nil {
		(*argsMap)["gap_color"] = parse.ParseColor((*argsMap)["gap_color"].(string))
	} else {
		(*argsMap)["gap_color"] = &parse.OptionalRGBA64{Color: TRANSPARENT_WHITE, Empty: false}
	}

	return argsMap
}

func FindEdges(read *im.Image, x int, y int, args *map[string]interface{}) *color.RGBA64 {

	moy_clr := NoiseReduction(read, x, y, args)
	ref := (*read).At(x, y)
	distFunc := distanceFunctions[((*args)["dist"]).(string)]
	threshold := (*args)["threshold"].(float64)

	r, g, b, a := ref.RGBA()
	original := &color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)} // or use original color

	if dist := distFunc(&ref, moy_clr); dist > threshold { // if very different, it's an edge

		if edge := (*args)["edge_color"]; edge != nil && !edge.(*parse.OptionalRGBA64).Empty { // if we have a custom color use it
			return (edge).(*parse.OptionalRGBA64).Color
		}
		return original

	} else if gap := (*args)["gap_color"]; gap != nil && !gap.(*parse.OptionalRGBA64).Empty { // if not very different, it's a gap. if we have a color for gap, use it
		return (gap).(*parse.OptionalRGBA64).Color
	}

	// else use original color (is gap & no custom color)
	return original
}

func ColorEuclideanDist(a *color.Color, b *color.RGBA64) float64 {

	aR, aG, aB, aA := (*a).RGBA()
	bR, bG, bB, bA := (*b).RGBA()

	dr := aR - bR
	dg := aG - bG
	db := aB - bB
	da := aA - bA

	return math.Sqrt(float64(dr*dr + dg*dg + db*db + da*da))
}

func ColorNorm1Dist(a *color.Color, b *color.RGBA64) float64 {

	aR, aG, aB, aA := (*a).RGBA()
	bR, bG, bB, bA := (*b).RGBA()

	dr := math.Abs(float64(aR) - float64(bR))
	dg := math.Abs(float64(aG) - float64(bG))
	db := math.Abs(float64(aB) - float64(bB))
	da := math.Abs(float64(aA) - float64(bA))

	return dr + dg + db + da
}

func ColorNormInfDist(a *color.Color, b *color.RGBA64) float64 {

	aR, aG, aB, aA := (*a).RGBA()
	bR, bG, bB, bA := (*b).RGBA()

	dr := math.Abs(float64(aR) - float64(bR))
	dg := math.Abs(float64(aG) - float64(bG))
	db := math.Abs(float64(aB) - float64(bB))
	da := math.Abs(float64(aA) - float64(bA))

	return quadMaxF64(dr, dg, db, da)
}
