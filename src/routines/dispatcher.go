package routines

import (
	filters "filter"
	"fmt"
	im "image"
	"imageio"
	"os"
	"strconv"
	"sync"
	"time"
)

var WorkQueue chan *Job
var (
	SliceHeight uint32
	SliceWidth  uint32
	WorkerCount uint32
)
var JobWaiter = sync.WaitGroup{}

type Job struct {
	InName                                  *string
	routines_count, SliceWidth, SliceHeight uint32
	Filter                                  *filters.Filter
}
type ImageSlice struct {
	image                      *im.Image
	writeImg                   *im.RGBA
	x_min, x_max, y_min, y_max uint32
}

func init() {

	WorkerCount, SliceWidth, SliceHeight = GetRoutineInfoFromArgs(os.Args[1:])
	WorkQueue = make(chan *Job, 2)
}

func StartWorker(i uint32) {
	fmt.Printf("Worker #%v => made\n", i)
	go func(i uint32) {
		for {
			fmt.Printf("Worker #%v => waiting\n", i)
			job := <-WorkQueue
			fmt.Printf("Worker #%v => new work\n", i)
			workerRoutine(job, i)
			JobWaiter.Done()
		}
	}(i)
}

func QueueJob(job *Job) {

	WorkQueue <- job
	fmt.Printf("=> Put %v on waiting list\n", *job.InName)
}

func workerRoutine(job *Job, i uint32) {

	begin := time.Now()
	waitSlices := sync.WaitGroup{}

	pimage, _ := imageio.LoadImage(job.InName)

	fmt.Printf("Worker #%v => loaded\n   -> Taille image (%v) : %v\n", i, *job.InName, (*pimage).Bounds())
	imW := im.NewRGBA((*pimage).Bounds())

	x_count := uint32(imW.Bounds().Max.X) / job.SliceWidth
	y_count := uint32(imW.Bounds().Max.Y) / job.SliceHeight
	for i := uint32(0); i < x_count; i++ {
		for j := uint32(0); j < y_count; j++ {
			waitSlices.Add(1)
			go func(i, j uint32) {

				slice := ImageSlice{
					image:    pimage,
					writeImg: imW,
					x_min:    job.SliceWidth * i,
					x_max:    job.SliceWidth * (i + 1),
					y_min:    job.SliceHeight * j,
					y_max:    job.SliceHeight * (j + 1),
				}
				if i == x_count-1 {
					slice.x_max = uint32(imW.Bounds().Max.Y)
				}
				if j == y_count-1 {
					slice.y_max = uint32(imW.Bounds().Max.X)
				}

				processSlice(&slice, job.Filter)
				waitSlices.Done()
			}(i, j)
		}
	}

	waitSlices.Wait()

	imageio.SaveImage(imW, job.InName, &job.Filter.Name)
	fmt.Printf("Worker #%v => finished\n   -> Taille image sortie (%v) : %v\n   -> Temps : %v\n", i, *job.InName, imW.Bounds(), time.Since(begin))
}

func processSlice(slice *ImageSlice, filter *filters.Filter) {

	for x := slice.x_min; x < slice.x_max; x++ {
		for y := slice.y_min; y < slice.y_max; y++ {
			(*slice.writeImg).Set(int(x), int(y), filter.Apply(slice.image, int(x), int(y)))
		}
	}
}

func GetRoutineInfoFromArgs(args []string) (uint32, uint32, uint32) {

	coroutine_per_img_count, width, height := uint32(0), uint32(0), uint32(0)
	for i := range args {
		if args[i] == "--routines" || args[i] == "-r" && i+1 < len(args) {
			r, _ := strconv.ParseUint(args[i+1], 10, 32)
			coroutine_per_img_count = uint32(r)
		} else if args[i] == "--swidth" || args[i] == "-w" && i+1 < len(args) {
			r, _ := strconv.ParseUint(args[i+1], 10, 32)
			width = uint32(r)
		} else if args[i] == "--sheight" || args[i] == "-h" && i+1 < len(args) {
			r, _ := strconv.ParseUint(args[i+1], 10, 32)
			height = uint32(r)
		}
	}

	if width <= 0 || height <= 0 || coroutine_per_img_count <= 0 {
		panic("Wrong parameter values for swidth(" + strconv.Itoa(int(width)) + ") / sheight (" + strconv.Itoa(int(height)) + ") routines(" + strconv.Itoa(int(coroutine_per_img_count)) + ")")
	}

	return coroutine_per_img_count, width, height
}
