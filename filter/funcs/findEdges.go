package funcs

import (
	"fmt"
	"github.com/mxyns/go-filter/filter"
	im "image"
	"image/color"
	"strconv"
)

func ParseFindEdgesArgs(filter *filter.Filter, args *[]string) *map[string]interface{} {

	argsMap := make(map[string]interface{})
	num, _ := strconv.Atoi((*args)[0])
	argsMap["size"] = num
	argsMap["test"] = (*args)[1]

	return &argsMap
}

func FindEdges(read *im.Image, x int, y int, args *map[string]interface{}) *color.RGBA64 {

	fmt.Printf("[FindEdges] args => %v\n", args)
	return &color.RGBA64{}
}
