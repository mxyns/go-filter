package filter

import (
	im "image"
	"image/color"
)

type Filter struct {
	Name   string
	Usage  string
	Parser func(filter *Filter, args *[]string) *map[string]interface{}
	Apply  func(read *im.Image, x int, y int, args *map[string]interface{}) *color.RGBA64
}

// nom -> filtre
var filters = map[string]*Filter{}

func GetFilter(filter_name string) *Filter {

	return filters[filter_name]
}
func RegisterFilter(filter *Filter) {

	filters[filter.Name] = filter
}
func GetFilterRegister() map[string]*Filter {

	return filters
}
