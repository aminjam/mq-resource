package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/aminjam/mq-resource/out"
	"github.com/aminjam/mq-resource/plugins"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}
	destination := os.Args[1]

	var request out.Request
	err := json.NewDecoder(os.Stdin).Decode(&request)
	checkError("reading request from stdin", err)

	//Capture Stdout while calling the plugin to aviod
	//polution and print out in stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	plugin, err := plugins.NewPlugin(request.Source)
	checkError("creating a plugin", err)

	file := request.Params.File
	if file == "" {
		file = "message.json"
	}

	message, err := ioutil.ReadFile(path.Join(destination, file))
	checkError("reading the output file", err)

	var response out.Response
	err = json.Unmarshal(message, &response)
	checkError("unmarshaling response", err)

	err = plugin.Put(message)
	checkError("running check on plugin", err)

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old
	//uncomment the line below if needed to see the stdout from calling the plugin
	//out := <-outC

	err = json.NewEncoder(os.Stdout).Encode(response)
	checkError("writing response", err)
}

func checkError(message string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s:%s", message, err.Error()))
		os.Exit(1)
	}
}
