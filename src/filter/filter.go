package filter

import (
	im "image"
	"image/color"
)

type Filter struct {
	Name  string
	Apply func(read im.Image, x int, y int) color.RGBA64
}

// nom -> filtre
var filters = map[string]*Filter{}

func GetSelectedFilter(args []string) *Filter {

	filter_name := ""
	for i := range args {
		if args[i] == "-f" || args[i] == "--filter" && i+1 < len(args) {
			filter_name = args[i+1]
		}
	}

	if len(filter_name) > 0 {
		return filters[filter_name]
	} else {
		panic("Pas de filtre donné")
	}
}
func GetFilter(filter_name string) *Filter {

	return filters[filter_name]
}
func RegisterFilter(filter Filter) {

	filters[filter.Name] = &filter
}
