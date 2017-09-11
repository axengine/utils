package nsq

import (
	"github.com/nsqio/go-nsq"
)

const (
	USER_AGENT = "go agent v0.1"
)

type Producer struct {
	producer *nsq.Producer
}

func NewProducer(nsqdTcpAddr string) *Producer {
	var p Producer
	cfg := nsq.NewConfig()
	cfg.UserAgent = USER_AGENT
	producer, err := nsq.NewProducer(nsqdTcpAddr, cfg)
	if err != nil {
		panic(err)
		return &p
	}
	p.producer = producer
	return &p
}

func (p *Producer) Publish(topic string, msg []byte) error {
	return p.producer.Publish(topic, msg)
}
