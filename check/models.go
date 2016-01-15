package check

import "github.com/aminjam/mq-resource"

type Request struct {
	Source  resource.Source    `json:"source"`
	Version resource.StringMap `json:"version"`
}

type Response []resource.StringMap
