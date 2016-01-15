package plugins

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/aminjam/mq-resource"
	"github.com/nsqio/go-nsq"
)

type nsqPlugin struct {
	Plugin
	Channel string
}

func (p nsqPlugin) Check() ([]resource.StringMap, error) {
	cfg := nsq.NewConfig()
	response := make([]resource.StringMap, 0)

	consumer, err := nsq.NewConsumer(p.Sub, p.Channel, cfg)
	if err != nil {
		return nil, err
	}
	consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		var msg resource.StringMap
		err := json.Unmarshal(m.Body, &msg)
		if err != nil {
			return err
		}
		response = append(response, msg)
		return nil
	}))
	err = consumer.ConnectToNSQD(p.Uri)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second * 5)
	consumer.Stop()
	<-consumer.StopChan
	return response, err
}

func (p nsqPlugin) Put(message []byte) error {
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(p.Uri, cfg)
	if err != nil {
		return err
	}
	err = producer.Publish(p.Pub, message)
	if err != nil {
		return err
	}
	producer.Stop()
	return nil
}

func (base *Plugin) nsq(params resource.StringMap) (Resourcer, error) {
	plugin := nsqPlugin{
		Plugin:  *base,
		Channel: params["channel"],
	}

	if plugin.Channel == "" {
		rand.Seed(time.Now().UnixNano())
		plugin.Channel = fmt.Sprintf("tail%06d#ephemeral", rand.Int()%999999)
	}
	return plugin, nil
}
