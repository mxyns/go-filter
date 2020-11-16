package pouf

import (
	"bufio"
	"fmt"
	"go-tcp/filet"
	"go-tcp/filet/requests"
	dRequests "go-tcp/filet/requests/default"
	"os"
	req "requests"
	"sync"
)

func MainClient(address *string, proto *string, port *uint, timeout *string) {

	client := &filet.Client{
		Address: &filet.Address{
			Proto: *proto,
			Addr:  *address,
			Port:  uint32(*port),
		},
	}
	defer client.Close()
	_, err := client.Start(*timeout)
	if err != nil {
		return
	}

	secondCommunicationMethod(client, true)

	//TODO lire le terminal du client pour faire une request
	//TODO quand l'user dit "no" on supprime aussi son fichier pr libérer de l'espace
	terminalInput(nil)
}

func firstCommunicationMethod(client *filet.Client) {
	response := *client.Send(requests.MakeGenericPack(
		dRequests.MakeTextRequest("invert nullify copy identity copy invert identity"),
		dRequests.MakeFileRequest("./in/sample-3.png", true),
	))

	pack := response.(*requests.Pack)
	fmt.Println("Received files :")
	for i := range pack.GetRequests() {
		fmt.Printf("	-> %v\n", (*pack.GetRequests()[i]).(*dRequests.FileRequest).GetPath())
	}
}

func secondCommunicationMethod(client *filet.Client, doTwice bool) {
	work_filter_list := (*client.Send(dRequests.MakeFileRequest("./in/sample-3.png", true))).(*req.WorkRequest)
	fmt.Printf("Step #%v -> server told me : %v\n", work_filter_list.GetStep(), work_filter_list.GetText())

	work_result_pack := (*client.Send(work_filter_list.Answer("invert"))).(*requests.Pack)
	results := (*work_result_pack.GetRequests()[0]).(*requests.Pack)
	work_continue := (*work_result_pack.GetRequests()[1]).(*req.WorkRequest)
	for i := range results.GetRequests() {
		fmt.Printf(
			"Step #%v -> got result file : path=%v size=%v\n",
			work_continue.GetStep(),
			(*results.GetRequests()[i]).(*dRequests.FileRequest).GetPath(),
			(*results.GetRequests()[i]).(*dRequests.FileRequest).GetFileSize(),
		)
	}

	fmt.Printf("Step #%v -> server told me : %v\n", work_continue.GetStep(), work_continue.GetText())

	if doTwice {
		client.Send(work_continue.Answer("yes"))
		secondCommunicationMethod(client, false)
	} else {
		client.Send(work_continue.Answer("no"))
	}
}

func terminalInput(group *sync.WaitGroup) {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		line := scanner.Text()
		fmt.Printf("got > %v\n", line)
		if line == "stop" {
			if group != nil {
				group.Done()
			}
			break
		}
	}
}
