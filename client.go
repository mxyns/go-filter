package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mxyns/go-filter/io"
	req "github.com/mxyns/go-filter/requests"
	"github.com/mxyns/go-tcp/filet"
	"github.com/mxyns/go-tcp/filet/requests"
	dRequests "github.com/mxyns/go-tcp/filet/requests/defaultRequests"
	"os"
	"strings"
)

var (
	filePathArg   *string
	filterListArg *string
)

func init() {

	filePathArg = flag.String("i", "", "client image file path to use")
	filterListArg = flag.String("fl", "copy", "client filter list to apply to image provided with -i. e.g: invert ; blur 5 ; etc")
}

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
		fmt.Println("Couldn't connect to server.")
		return
	}

	if len(*filePathArg) > 0 {
		firstCommunicationMethod(client)
	} else {
		secondCommunicationMethod(client, nil)
		terminalInput() // opération bloquant la goroutine principale
	}
}

func firstCommunicationMethod(client *filet.Client) {

	//send a pack with the args
	response := *client.Send(requests.MakeGenericPack(
		dRequests.MakeTextRequest(strings.TrimSpace(*filterListArg)),
		dRequests.MakeFileRequest(strings.TrimSpace(*filePathArg), true),
	))

	pack := response.(*requests.Pack)

	fmt.Printf("Step #1 got results : \n")
	io.RenameFiles(pack, filePathArg, *filterListArg)
}

func secondCommunicationMethod(client *filet.Client, previous *req.WorkRequest) {

	work_filter_list := previous
	scanner := bufio.NewScanner(os.Stdin)
	if previous == nil { // first run of the function (not after a "yes" when asked if wanna reuse img)

		*filePathArg = userInput(scanner, "File path > ") // use filePathArg to store value
		work_filter_list = (*client.Send(dRequests.MakeFileRequest(*filePathArg, true))).(*req.WorkRequest)
	}

	fmt.Printf("Step #%v -> \n", work_filter_list.GetStep())
	fmt.Printf("%v\n", work_filter_list.GetText())

	filters := userInput(scanner, "filters > ")

	// work_result_pack:Pack { results:Pack {img, img, img, ...}, work_continue:WorkRequest }
	work_result_pack := (*client.Send(work_filter_list.Answer(filters))).(*requests.Pack)
	results := (*work_result_pack.GetRequests()[0]).(*requests.Pack)
	work_continue := (*work_result_pack.GetRequests()[1]).(*req.WorkRequest)

	// treat result (rename & print)
	fmt.Printf("Step #1 got results : \n")
	io.RenameFiles(results, filePathArg, filters)

	fmt.Printf("Step #%v -> ", work_continue.GetStep())
	repeatInput := userInput(scanner, work_continue.GetText()+" > ")
	if repeatInput == "y" || repeatInput == "yes" {
		work_filter_list = (*client.Send(work_continue.Answer("yes"))).(*req.WorkRequest)
		secondCommunicationMethod(client, work_filter_list) // reuse this function with context
	} else {
		client.Send(work_continue.Answer("no"))
	}
}

func userInput(scanner *bufio.Scanner, text string) string {

	fmt.Printf(text)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
