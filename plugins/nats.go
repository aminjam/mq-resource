package plugins

import (
	"encoding/json"
	"time"

	"github.com/aminjam/mq-resource"
	"github.com/nats-io/nats"
)

type natsPlugin struct {
	Plugin
}

func (p natsPlugin) Check() ([]resource.StringMap, error) {
	c, err := nats.Connect(p.Uri)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	response := make([]resource.StringMap, 0)
	c.Subscribe(p.Sub, func(m *nats.Msg) {
		var msg resource.StringMap
		err := json.Unmarshal(m.Data, &msg)
		if err != nil {
			return
		}
		response = append(response, msg)
	})
	time.Sleep(time.Second * 5)
	sub, err := c.Subscribe(p.Sub, nil)
	if err != nil {
		return nil, err
	}
	sub.Unsubscribe()
	return response, nil
}

func (p natsPlugin) Put(message []byte) error {
	nc, err := nats.Connect(p.Uri)
	if err != nil {
		return err
	}
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return err
	}
	err = c.Publish(p.Pub, message)
	if err != nil {
		return err
	}
	c.Close()
	return nil
}

func (base *Plugin) nats(params resource.StringMap) (Resourcer, error) {
	plugin := natsPlugin{
		Plugin: *base,
	}

	return plugin, nil
}
