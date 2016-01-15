package plugins

import (
	"fmt"

	"github.com/aminjam/mq-resource"
)

type Resourcer interface {
	Check() ([]resource.StringMap, error)
	Put([]byte) error
}

type Plugin struct {
	Uri string
	Pub string
	Sub string
}

func NewPlugin(source resource.Source) (Resourcer, error) {
	plugin := &Plugin{
		Uri: source.Uri,
		Sub: source.Sub,
		Pub: source.Pub,
	}
	switch source.Queue {
	case "nsq":
		return plugin.nsq(source.Params)
	}
	return nil, fmt.Errorf("plugin not found.")
}
