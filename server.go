package main

import (
	"fmt"
	filters "github.com/mxyns/go-filter/filter"
	"github.com/mxyns/go-filter/io"
	req "github.com/mxyns/go-filter/requests"
	disp "github.com/mxyns/go-filter/routines"
	"github.com/mxyns/go-tcp/filet"
	gotcp "github.com/mxyns/go-tcp/filet/requests"
	dRequests "github.com/mxyns/go-tcp/filet/requests/defaultRequests"
	"net"
	"strings"
	"sync"
)

var (
	clientStates = make(map[*net.Conn]*WorkState)
)

type WorkState struct {
	current_path *string
	step         byte
}

func startServer(address *string, proto *string, port *uint) {

	// start n worker goroutines
	for i := uint(0); i < *disp.WorkerCount; i++ {
		disp.StartWorker(i)
	}

	server := filet.Server{
		Address: &filet.Address{
			Proto: *proto,
			Addr:  *address,
			Port:  uint32(*port),
		},
		Clients:          make([]*net.Conn, 0),
		ConnectionWaiter: disp.JobWaiter,
		RequestHandler:   requestHandler,
	}
	defer server.Close()
	defer close(disp.WorkQueue)
	go server.Start()

	// TODO clean stop (server will process but will not send back last processed image when stopped during processing)
	terminalInput()       // bloque jusqu'à lire "stop"
	disp.JobWaiter.Wait() // bloque jusqu'à ce que tous les jobs soient terminés
}

func requestHandler(client *net.Conn, request *gotcp.Request) {

	if request == nil {
		_ = (*client).Close()
		return
	}

	switch (*request).(type) {
	case *gotcp.Pack:
		{
			handlePack(client, request)
		}
	case *dRequests.FileRequest:
		{
			path_file := (*request).(*dRequests.FileRequest).GetPath()
			clientStates[client] = &WorkState{
				step:         0,
				current_path: &path_file,
			}

			sendFilterList(client)
		}
	case *req.WorkRequest:
		{
			if clientStates[client] != nil {
				handleWorkRequest(client, request)
			} else {
				fmt.Printf("Client [%v] sent WorkRequest without File", client)
			}
		}
	}
}

// first way for clients to communicate with server : send a Pack with a Text and a File (in this order)
// Text content : "filter1 arg1 arg2 arg3 ... ; filter2 arg1 arg2 ..."
// File content : png/jpg image
func handlePack(client *net.Conn, request *gotcp.Request) {
	pack := (*request).(*gotcp.Pack)
	list := pack.GetRequests() // 0 : image, 1 : texte (liste filtres)
	text := (*list[0]).(*dRequests.TextRequest).GetText()

	fmt.Printf("Image : path=%v, size=%v\n", (*list[1]).(*dRequests.FileRequest).GetPath(), (*list[1]).(*dRequests.FileRequest).GetFileSize())
	fmt.Printf("Filtres : %v\n", strings.Split(text, " "))

	syncPoint := sync.WaitGroup{}
	jobs := filterJob(text, (*list[1]).(*dRequests.FileRequest).GetPath(), &syncPoint)
	syncPoint.Wait()

	if pack.Info().WantsResponse {

		response := packJobResults(gotcp.MAX_PACK_SUBID+1+pack.SubId, jobs)

		_, err, err_id := gotcp.SendRequestOn(client, &response)
		if err != nil {
			fmt.Printf("Error while responding to Client %v : %v", err_id, err)
		}
	}
}

/* second way for clients to communicate with server
step :

	client 		=> file
	0. server 	=> Work { etape:0, texte: "liste des filtres" }
	client 		=> Work { etape:0, texte: "filtre que jveux" }
	server 		=> Pack { Pack { requests: fichiers traités }, Work { etape:1, texte : "reuse img ?" } }
	client 		=> Work {etape:1, yes/no}
					-> yes -> server => goto 0.
						-> no -> close()
*/
func handleWorkRequest(client *net.Conn, request *gotcp.Request) {

	workReq := (*request).(*req.WorkRequest)
	step := workReq.GetStep()
	text := workReq.GetText()
	state := clientStates[client]
	state.step = step

	if step == 0 {

		syncPoint := sync.WaitGroup{}

		// process image
		jobs := filterJob(text, *state.current_path, &syncPoint)
		syncPoint.Wait() // wait for all jobs termination

		var response gotcp.Request = gotcp.MakePack(
			1,
			packJobResults(1, jobs), // make pack from jobs results
			req.MakeWorkRequest(1, "Do you want to reuse this image ? yes/no"),
		)

		// send result and get client response
		clientResponse, _, _ := gotcp.SendRequestOn(client, &response)
		// pass it through requestHandler (will come back to handleWorkRequest w/ step=1)
		requestHandler(client, clientResponse)

	} else if step == 1 {

		if strings.ToLower(text) == "yes" || strings.ToLower(text) == "y" {
			sendFilterList(client)
		} else {
			_ = (*client).Close()
			_ = io.RemoveFile(clientStates[client].current_path)
			clientStates[client] = nil
		}
	}
}

// send filter list to client
func sendFilterList(client *net.Conn) {

	filterList := ""
	filtersMap := filters.GetFilterRegister() // map : nom filtre -> Filtre
	for name := range filtersMap {
		filterList += name + " | usage : " + filtersMap[name].Usage + "\n"
	}
	filterList = strings.TrimSpace(filterList)

	var response gotcp.Request = req.MakeWorkRequest(0, "Filter list : \n"+filterList)

	clientResponse, _, _ := gotcp.SendRequestOn(client, &response)
	requestHandler(client, clientResponse)
}

func filterJob(filterList string, filename string, syncPoint *sync.WaitGroup) []*disp.Job {

	filter := strings.Split(filterList, ";")

	jobs := make([]*disp.Job, len(filter))
	for i := range filter {
		shards := strings.Split(strings.TrimSpace(filter[i]), " ")
		if filter := filters.GetFilter(shards[0]); filter != nil {
			jobs[i] = disp.QueueJob(&disp.Job{
				InName:          &filename,
				Filter:          filter,
				FilterArgs:      shards[1:],
				HorizSliceCount: disp.HorizSliceCount,
				VertSliceCount:  disp.VertSliceCount,
				SyncPoint:       syncPoint,
			})
		} else {
			jobs[i] = nil
		}
	}

	return jobs
}

func packJobResults(id uint32, jobs []*disp.Job) gotcp.Request {

	imagesRequests := make([]gotcp.Request, len(jobs))
	for i := range jobs {
		if job := jobs[i]; job != nil {
			imagesRequests[i] = dRequests.MakeFileRequest(*jobs[i].OutPath, false)
		} else {
			imagesRequests[i] = dRequests.MakeTextRequest("Unknown filter")
		}
	}

	return gotcp.MakePack(id, imagesRequests...)
}
