package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aminjam/mq-resource/in"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}
	destination := os.Args[1]

	var request in.Request
	err := json.NewDecoder(os.Stdin).Decode(&request)
	checkError("reading request from stdin", err)

	message, err := json.Marshal(request.Version)
	checkError("marshalling the result", err)

	err = os.MkdirAll(destination, 0755)
	checkError("creating destination", err)

	if request.Params.File == "" {
		request.Params.File = "message.json"
	}

	file, err := os.Create(filepath.Join(destination, request.Params.File))
	checkError("opening the file", err)
	defer file.Close()
	_, err = fmt.Fprintf(file, "%s", message)
	checkError("writing to the file", err)

	err = json.NewEncoder(os.Stdout).Encode(request.Version)
	checkError("encoding recieved data", err)
}
func checkError(message string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
