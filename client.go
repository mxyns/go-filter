package main

import (
	"fmt"
	req "github.com/mxyns/go-filter/requests"
	"github.com/mxyns/go-tcp/filet"
	"github.com/mxyns/go-tcp/filet/requests"
	dRequests "github.com/mxyns/go-tcp/filet/requests/defaultRequests"
)

func startClient(address *string, proto *string, port *uint, timeout *string) {

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

	firstCommunicationMethod(client)

	//TODO lire le terminal du client pour faire une request
	terminalInput() // opération bloquant la goroutine principale
}

func firstCommunicationMethod(client *filet.Client) {
	response := *client.Send(requests.MakeGenericPack(
		dRequests.MakeTextRequest("invert"),
		dRequests.MakeFileRequest("./in/7.png", true),
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
