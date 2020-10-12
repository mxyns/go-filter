package main

import (
	filters "filter"
	filfuncs "filter/funcs"
	"fmt"
	im "image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"imageio"
	"os"
)

// TODO : image mutable
// TODO filtres : negatif, gris, reduction bruit (moyenne pix alentours), bords (diff), code barre
func main() {

	registerFilters()

	fmt.Printf("\n")
	images_names := imageio.GetImageNames(os.Args[1:])
	selected_filter := filters.GetSelectedFilter(os.Args[1:])

	fmt.Printf("Images : %v\n", images_names)
	for i := range images_names {
		pimage, _ := imageio.LoadImage(images_names[i])

		fmt.Printf("Taille image (%v) : %v\n", images_names[i], (*pimage).Bounds())
		imW := im.NewRGBA((*pimage).Bounds())
		for x := imW.Bounds().Min.X; x < imW.Bounds().Max.X; x++ {
			for y := imW.Bounds().Min.Y; y < imW.Bounds().Max.Y; y++ {
				imW.Set(x, y, selected_filter.Apply(pimage, x, y))
			}
		}

		imageio.SaveImage(imW, images_names[i], selected_filter.Name)

		fmt.Printf("Taille image sortie (%v) : %v\n", images_names[i], imW.Bounds())
	}
}

func registerFilters() {

	filters.RegisterFilter(&filters.Filter{Name: "invert", Apply: filfuncs.InvertColor})
	filters.RegisterFilter(&filters.Filter{Name: "nullify", Apply: filfuncs.Nullify})
}
