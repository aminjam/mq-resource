package out

import "github.com/aminjam/mq-resource"

type Request struct {
	Source  resource.Source    `json:"source"`
	Version resource.StringMap `json:"version"`
	Params  Params             `json:"params"`
}

type Response resource.StringMap

type Params struct {
	File string `json:"file"`
}
