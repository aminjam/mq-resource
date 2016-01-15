package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aminjam/mq-resource/check"
	"github.com/aminjam/mq-resource/plugins"
)

func main() {
	var request check.Request
	err := json.NewDecoder(os.Stdin).Decode(&request)

	//Capture Stdout while calling the plugin to aviod
	//polution and print out in stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	checkError("reading request from stdin", err)
	plugin, err := plugins.NewPlugin(request.Source)
	checkError("creating a plugin", err)

	versions, err := plugin.Check()
	checkError("running check on plugin", err)

	response := make(check.Response, 0)
	for _, v := range versions {
		if len(v) > 0 && !v.IsEqual(request.Version) {
			response = append(response, v)
		}
	}

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
