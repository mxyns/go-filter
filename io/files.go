package io

import (
	"fmt"
	"github.com/mxyns/go-tcp/filet/requests"
	dRequests "github.com/mxyns/go-tcp/filet/requests/defaultRequests"
	"os"
	"strings"
)

func FormatImageName(path *string, filter_name *string) (new_name, new_path string) {

	shards := strings.Split(*path, "/")
	filename := shards[len(shards)-1]

	new_path = strings.Join(shards[:len(shards)-1], "/")

	shards = strings.Split(filename, ".")
	new_name = ""
	if len(shards) > 1 {
		new_name = strings.Join(shards[:len(shards)-1], ".")
		new_name += "-" + *filter_name + "." + shards[len(shards)-1]
	} else {
		new_name = shards[0] + "-" + *filter_name + ".png"
	}
	new_path += "/" + new_name

	return
}
func FormatImageTempName(temp_path, input_file_path, filter_name *string) (new_name, new_path string) {

	// ./some_folder/some.input.file.png
	shards := strings.Split(*input_file_path, "/") // [. some_folder input.png]
	filename := shards[len(shards)-1]              // some.input.file.png

	// ./some_other_folder/output_folder/46513213676574.png
	shards = strings.Split(*temp_path, "/")              // [. some_other_folder output_folder 46513213676574.png]
	new_path = strings.Join(shards[:len(shards)-1], "/") // ./some_other_folder/output_folder

	shards = strings.Split(filename, ".") // [some input file png]
	new_name = ""
	if len(shards) > 1 {
		new_name = strings.Join(shards[:len(shards)-1], ".")         // some.input.file
		new_name += "-" + *filter_name + "." + shards[len(shards)-1] // some.input.file-myfilter.png
	} else {
		new_name = shards[0] + "-" + *filter_name + ".png" // or input-myfilter.png
	}
	new_path += "/" + new_name // ./some_other_folder/output_folder/some.input.file-myfilter.png

	return
}

func RenameFiles(pack *requests.Pack, input_file_path *string, filterList string) {

	filtersCmds := strings.Split(filterList, ";")
	for i, filter := range filtersCmds {
		filtersCmds[i] = strings.TrimSpace(filter)
	}

	for i, result := range pack.GetRequests() {
		switch (*result).(type) {

		case *dRequests.FileRequest:
			{
				temp_file_path := (*result).(*dRequests.FileRequest).GetPath()
				_, new_path := FormatImageTempName(&temp_file_path, input_file_path, &strings.Split(filtersCmds[i], " ")[0])

				err := os.Rename(temp_file_path, new_path)
				if err != nil {
					fmt.Printf("			/!\\ error while renaming : %v\n", err)
					new_path = temp_file_path
				}
				fmt.Printf("		-> file : path=%v size=%v\n",
					new_path,
					(*result).(*dRequests.FileRequest).GetFileSize())
			}
		case *dRequests.TextRequest:
			{
				fmt.Printf("		-> error : msg=%v\n",
					(*result).(*dRequests.TextRequest).GetText())
			}
		}
	}
}

func RemoveFile(filename string) error {

	return os.Remove(filename)
}
