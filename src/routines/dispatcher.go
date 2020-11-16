package routines

import (
	filters "filter"
	"flag"
	"fmt"
	im "image"
	"imageio"
	"sync"
	"time"
)

var WorkQueue chan *Job
var (
	SliceHeight *uint
	SliceWidth  *uint
	WorkerCount *uint
)
var JobWaiter = sync.WaitGroup{}

type Job struct {
	InName                                  *string
	OutPath                                 *string
	routines_count, SliceWidth, SliceHeight *uint
	Filter                                  *filters.Filter
	SyncPoint                               *sync.WaitGroup
}
type ImageSlice struct {
	image                      *im.Image
	writeImg                   *im.RGBA
	x_min, x_max, y_min, y_max uint32
}

func init() {

	WorkerCount = flag.Uint("r", 1, "number of image processing routines")
	SliceWidth = flag.Uint("swidth", 100, "image slice width")
	SliceHeight = flag.Uint("sheight", 100, "image slice height")

	WorkQueue = make(chan *Job, 16)
}

func StartWorker(i uint) {
	fmt.Printf("Worker #%v => made\n", i)
	go func(i uint) {
		for {
			fmt.Printf("Worker #%v => waiting\n", i)
			job := <-WorkQueue
			fmt.Printf("Worker #%v => new work\n", i)
			workerRoutine(job, i)
			JobWaiter.Done()
		}
	}(i)
}

func QueueJob(job *Job) *Job {

	JobWaiter.Add(1)
	if job.SyncPoint != nil {
		job.SyncPoint.Add(1)
	}

	WorkQueue <- job
	fmt.Printf("=> Put %v on waiting list\n", *job.InName)

	return job
}

func workerRoutine(job *Job, i uint) {

	begin := time.Now()
	waitSlices := sync.WaitGroup{}

	pimage, _ := imageio.LoadImage(job.InName)

	fmt.Printf("Worker #%v => loaded\n   -> Taille image (%v) : %v\n", i, *job.InName, (*pimage).Bounds())
	imW := im.NewRGBA((*pimage).Bounds())

	if uint32(imW.Bounds().Max.X) < uint32(*job.SliceWidth) {
		*job.SliceWidth = uint(imW.Bounds().Max.X)
	}
	if uint32(imW.Bounds().Max.Y) < uint32(*job.SliceHeight) {
		*job.SliceHeight = uint(imW.Bounds().Max.Y)
	}
	x_count := uint32(imW.Bounds().Max.X) / uint32(*job.SliceWidth)
	y_count := uint32(imW.Bounds().Max.Y) / uint32(*job.SliceHeight)

	//TODO si on a le temps faire un modulo pour faire une slice plus petite si jamais il y a un reste >0
	for i := uint32(0); i < x_count; i++ {
		for j := uint32(0); j < y_count; j++ {
			waitSlices.Add(1)
			go func(i, j uint32) {

				slice := ImageSlice{
					image:    pimage,
					writeImg: imW,
					x_min:    uint32(*job.SliceWidth) * i,
					x_max:    uint32(*job.SliceWidth) * (i + 1),
					y_min:    uint32(*job.SliceHeight) * j,
					y_max:    uint32(*job.SliceHeight) * (j + 1),
				}
				if i == x_count-1 {
					slice.x_max = uint32(imW.Bounds().Max.X)
				}
				if j == y_count-1 {
					slice.y_max = uint32(imW.Bounds().Max.Y)
				}

				processSlice(&slice, job.Filter)
				waitSlices.Done()
			}(i, j)
		}
	}

	waitSlices.Wait()

	job.OutPath = imageio.SaveImage(imW, job.InName, &job.Filter.Name)
	fmt.Printf("Worker #%v => finished\n   -> Taille image sortie (%v) : %v\n   -> Temps : %v\n", i, *job.OutPath, imW.Bounds(), time.Since(begin))

	if job.SyncPoint != nil {
		job.SyncPoint.Done()
	}
}

func processSlice(slice *ImageSlice, filter *filters.Filter) {

	for x := slice.x_min; x < slice.x_max; x++ {
		for y := slice.y_min; y < slice.y_max; y++ {
			(*slice.writeImg).Set(int(x), int(y), filter.Apply(slice.image, int(x), int(y)))
		}
	}
}
