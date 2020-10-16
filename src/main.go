package main

import (
	filters "filter"
	filfuncs "filter/funcs"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"imageio"
	"os"
	disp "routines"
)

// TODO : image mutable
// TODO filtres : negatif, gris, reduction bruit (moyenne pix alentours), bords (diff), code barre
func main() {

	registerFilters()

	fmt.Printf("\n")
	images_names := imageio.GetImageNames(os.Args[1:])
	selected_filter := filters.GetSelectedFilter(os.Args[1:])

	fmt.Printf("Images : %v\n", images_names)

	for i := uint32(0); i < disp.WorkerCount; i++ {
		disp.StartWorker(i)
	}

	disp.JobWaiter.Add(len(images_names))
	for i := range images_names {

		disp.QueueJob(&disp.Job{
			InName:      &images_names[i],
			Filter:      selected_filter,
			SliceWidth:  disp.SliceWidth,
			SliceHeight: disp.SliceHeight,
		})
	}

	disp.JobWaiter.Wait()
	close(disp.WorkQueue)
}

func registerFilters() {

	filters.RegisterFilter(&filters.Filter{Name: "invert", Apply: filfuncs.InvertColor})
	filters.RegisterFilter(&filters.Filter{Name: "nullify", Apply: filfuncs.Nullify})
	filters.RegisterFilter(&filters.Filter{Name: "copy", Apply: filfuncs.Identity})
	filters.RegisterFilter(&filters.Filter{Name: "identity", Apply: filfuncs.Identity})
	filters.RegisterFilter(&filters.Filter{Name: "print", Apply: filfuncs.Print})
}
