package routines

import (
	"flag"
	"fmt"
	filters "github.com/mxyns/go-filter/filter"
	"github.com/mxyns/go-filter/io"
	im "image"
	"strings"
	"sync"
	"time"
)

const (
	WorkQueueSize = 16
)

var (
	HorizSliceCount *uint // see init()
	VertSliceCount  *uint // see init()

	WorkerCount *uint     // number of worker routines wanted
	WorkQueue   chan *Job // Job input queue for workers

	JobWaiter = &sync.WaitGroup{} // global WaitGroup (+1 when a Job is queued, -1 when a Job is done)
)

type Job struct {
	InName                          *string         // path to input file
	OutPath                         *string         // path to output file
	VertSliceCount, HorizSliceCount *uint           // horizontal and vertical image fragmentation count
	Filter                          *filters.Filter // filter to apply
	FilterArgs                      []string        // arguments for the filter
	SyncPoint                       *sync.WaitGroup // WaitGroup used to synchronize wait for multiple Job termination
	Duration                        *time.Duration  // image processing duration
}
type ImageSlice struct {
	image                      *im.Image // image to read on
	writeImg                   *im.RGBA  // image to write on
	x_min, x_max, y_min, y_max uint      // slice bounds
}

func init() {
	WorkerCount = flag.Uint("r", 1, "number of image processing routines")
	VertSliceCount = flag.Uint("scvert", 5, "vertical slice count per image")
	HorizSliceCount = flag.Uint("schor", 5, "horizontal slice count per image")

	WorkQueue = make(chan *Job, WorkQueueSize)
}

func StartWorker(i uint) {
	fmt.Printf("Worker #%v => made\n", i)
	go func(i uint) {
		for {
			fmt.Printf("Worker #%v => waiting\n", i)
			job := <-WorkQueue

			if job == nil { // WorkQueue got closed
				break
			}

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

//Calculate image slice + call processSlice for each slice
func workerRoutine(job *Job, i uint) {

	fmt.Printf("\t%v -> Raw Arguments (%v) : %v\n", i, len(job.FilterArgs), job.FilterArgs)

	var args = make(map[string]interface{})
	var argsPtr *map[string]interface{}
	for _, val := range job.FilterArgs {
		if arr := strings.Split(val, "="); len(arr) == 2 {
			args[arr[0]] = arr[1]
		}
	}
	if parser := job.Filter.Parser; parser != nil {
		argsPtr = parser(job.Filter, &args)
	} else {
		argsPtr = nil
	}

	begin := time.Now()
	waitSlices := sync.WaitGroup{}

	pimage, _ := io.LoadImage(job.InName)

	fmt.Printf("Worker #%v => loaded\n\t%v -> Taille image (%v) : %v\n", i, i, *job.InName, (*pimage).Bounds())
	fmt.Printf("\t%v -> Arguments (%v) : %v\n", i, len(job.FilterArgs), args)
	imW := im.NewRGBA((*pimage).Bounds())

	sliceWidth := uint(imW.Bounds().Max.X) / *HorizSliceCount
	sliceHeight := uint(imW.Bounds().Max.Y) / *VertSliceCount
	for i := uint(0); i < *job.HorizSliceCount; i++ {
		for j := uint(0); j < *job.VertSliceCount; j++ {
			waitSlices.Add(1)
			go func(i, j uint) {

				slice := ImageSlice{
					image:    pimage,
					writeImg: imW,
					x_min:    sliceWidth * i,
					x_max:    sliceWidth * (i + 1),
					y_min:    sliceHeight * j,
					y_max:    sliceHeight * (j + 1),
				}

				// correct precision loss with integer division
				if i == *job.VertSliceCount-1 {
					slice.x_max = uint(imW.Bounds().Max.X)
				}
				if j == *job.VertSliceCount-1 {
					slice.y_max = uint(imW.Bounds().Max.Y)
				}

				processSlice(&slice, job.Filter, argsPtr)
				waitSlices.Done()
			}(i, j)
		}
	}

	waitSlices.Wait()

	duration := time.Since(begin)
	job.Duration = &duration
	job.OutPath = io.SaveImage(imW, job.InName, &job.Filter.Name, &job.FilterArgs)
	fmt.Printf("Worker #%v => finished\n   %v -> Taille image sortie (%v) : %v\n   %v -> Temps : %v\n", i, i, *job.OutPath, imW.Bounds(), i, job.Duration)

	if job.SyncPoint != nil {
		job.SyncPoint.Done()
	}
}

//Apply filter + write the image for each slice
func processSlice(slice *ImageSlice, filter *filters.Filter, args *map[string]interface{}) {

	for x := slice.x_min; x < slice.x_max; x++ {
		for y := slice.y_min; y < slice.y_max; y++ {
			(*slice.writeImg).Set(int(x), int(y), filter.Apply(slice.image, int(x), int(y), args))
		}
	}
}
