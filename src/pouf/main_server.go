package pouf

import (
	filters "filter"
	"fmt"
	"go-tcp/filet"
	gotcp "go-tcp/filet/requests"
	dRequests "go-tcp/filet/requests/default"
	"net"
	req "requests"
	disp "routines"
	"strings"
	"sync"
)

/*
	type Work struct {
		etape uint
		texte string
	}

	etapes :
		client => file
		0. server => Work {etape:0, texte: "liste des filtres"}
		client => Work {etape:0, texte: "filtre que jveux"}
		server => Pack { Pack {requests: fichiers traités}, Work {etape:1, texte : "reuse img ?"} }
		client => Work {etape:1, yes/no}
			-> yes -> server => 0.
			-> no -> close()

*/

var (
	clientStates = make(map[*net.Conn]*WorkState)
)

type WorkState struct {
	current_path *string
	step         byte
}

func MainServer(address *string, proto *string, port *uint) {

	fmt.Printf("\n")
	for i := uint(0); i < *disp.WorkerCount; i++ {
		disp.StartWorker(i)
	}

	server := filet.Server{
		Address: &filet.Address{
			Proto: *proto,
			Addr:  *address,
			Port:  uint32(*port),
		},
		Clients:          make([]*net.Conn, 5),
		ConnectionWaiter: &disp.JobWaiter,
		RequestHandler:   requestHandler,
	}
	defer server.Close()
	server.Start()

	terminalInput(&disp.JobWaiter)
	disp.JobWaiter.Wait()
	close(disp.WorkQueue)
}

func requestHandler(client *net.Conn, request *gotcp.Request) {
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

		response := packJobResults(gotcp.HalfUInt32+1+pack.SubId, jobs)

		_, err, err_id := gotcp.SendRequestOn(client, &response)
		if err != nil {
			fmt.Printf("Error while responding to Client %v : %v", err_id, err)
		}
	}
}

func handleWorkRequest(client *net.Conn, request *gotcp.Request) {

	workReq := (*request).(*req.WorkRequest)
	step := workReq.GetStep()
	text := workReq.GetText()
	state := clientStates[client]
	state.step = step

	if step == 0 {

		syncPoint := sync.WaitGroup{}
		jobs := filterJob(text, *state.current_path, &syncPoint)
		syncPoint.Wait()

		var response gotcp.Request = gotcp.MakePack(
			1,
			packJobResults(1, jobs),
			req.MakeWorkRequest(1, "Do you want to reuse this image ? yes/no"),
		)

		clientResponse, _, _ := gotcp.SendRequestOn(client, &response)
		requestHandler(client, clientResponse)

	} else if step == 1 {

		if strings.ToLower(text) == "yes" || strings.ToLower(text) == "y" {
			sendFilterList(client)
		} else {
			_ = (*client).Close()
			clientStates[client] = nil
		}
	}
}

func sendFilterList(client *net.Conn) {

	filterList := ""
	filtersMap := filters.GetFilterRegister() // nom filtre -> Filtre
	for name := range filtersMap {
		filterList += name + " "
	}
	filterList = strings.TrimSpace(filterList)

	var response gotcp.Request = req.MakeWorkRequest(0, "Filter list : "+filterList)

	clientResponse, _, _ := gotcp.SendRequestOn(client, &response)
	requestHandler(client, clientResponse)
}

func filterJob(filterList string, filename string, syncPoint *sync.WaitGroup) []*disp.Job {

	filter := strings.Split(filterList, " ")

	jobs := make([]*disp.Job, len(filter))
	for i := range filter {
		filename := filename
		jobs[i] = disp.QueueJob(&disp.Job{
			InName:      &filename,
			Filter:      filters.GetFilter(filter[i]),
			SliceWidth:  disp.SliceWidth,
			SliceHeight: disp.SliceHeight,
			SyncPoint:   syncPoint,
		})
	}

	return jobs
}

func packJobResults(id uint32, jobs []*disp.Job) gotcp.Request {

	imagesRequests := make([]gotcp.Request, len(jobs))
	for i := range jobs {
		imagesRequests[i] = dRequests.MakeFileRequest(*jobs[i].OutPath, false)
	}

	return gotcp.MakePack(id, imagesRequests...)
}
