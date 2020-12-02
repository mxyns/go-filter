package filter

import (
	im "image"
	"image/color"
	"strings"
)

type Filter struct {
	Name   string
	Usage  string
	Parser func(filter *Filter, args *map[string]interface{}) *map[string]interface{}
	Apply  func(read *im.Image, x int, y int, args *map[string]interface{}) *color.RGBA64
}

// nom -> filtre
var filters = map[string]*Filter{}

func GetFilter(filter_name string) *Filter {

	return filters[strings.ToLower(filter_name)]
}
func RegisterFilter(filter *Filter) {

	filters[strings.ToLower(filter.Name)] = filter
}
func GetFilterRegister() map[string]*Filter {

	return filters
}
