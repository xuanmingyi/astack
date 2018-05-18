package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

var base_url = "http://localhost:9090"

type command_func func([]string)

func panic_exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func do_get(url string) {
}


func args_parse(args []string, name string, require_value bool) string {
	for index, value := range args {
		if name == value {
			if require_value {
				if index+1 <= len(args) {
					return args[index+1]
				} else {
					panic_exit(fmt.Sprintf("%s must need a value", name))
				}
			} else {
				return "true"
			}
		}
	}
	return "false"
}

func image_create(args []string) {
	//astack image-create --name cirros --file cirros.img --format qcow2
	url := fmt.Sprintf("%s/image", base_url)

	//name := args_parse(args, "--name", true)
	file := args_parse(args, "--file", true)
	//format := args_parse(args, "--format", true)

	body_buf := &bytes.Buffer{}
	body_writer := multipart.NewWriter(body_buf)

	file_writer, err := body_writer.CreateFormFile("uploadfile", file)
	if err != nil {
		panic_exit("Error writing to buffer")
	}

	file_hander, err := os.Open(file)
	if err != nil {
		panic_exit("Error opening file")
	}
	defer file_hander.Close()

	_, err = io.Copy(file_writer, file_hander)
	if err != nil {
		panic_exit("Error copy file")
	}
	content_type := body_writer.FormDataContentType()
	body_writer.Close()

	resp, err := http.Post(url, content_type, body_buf)
	if err != nil {
		panic_exit("Error in post")
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic_exit("Error in read body")
	}
	fmt.Println(string(resp_body))
}

func image_list(args []string) {
	fmt.Println("image_list", args)
}

var command_mappings = map[string]command_func{
	"image-create": image_create,
	"image-list":   image_list,
}

func main() {
	if len(os.Args) == 1 {
		panic_exit("Error")
	}
	if _, ok := command_mappings[os.Args[1]]; ok {
		command := os.Args[1]
		command_mappings[command](os.Args[2:])
	} else {
		panic_exit("Error")
	}
}
