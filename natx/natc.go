package natx

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

// 默认为subj和queue名称加前缀
const (
	SUBJ_PRE  = "NCSUBJ-"
	QUEUE_PRE = "NCQUE-"
)

type NATClient struct {
	conn *nats.EncodedConn
}

func NewNATClient(endpoints []string) *NATClient {
	opts := nats.Options{
		Servers:        endpoints,
		AllowReconnect: true,
		MaxReconnect:   10,
		ReconnectWait:  5 * time.Second,
		Timeout:        3 * time.Second,
		PingInterval:   30 * time.Second,
	}
	nc, err := opts.Connect()
	if err != nil {
		panic(err)
	}
	nats.RegisterEncoder(MSGP_ENCODER, &MsgpEncoder{})
	c, err := nats.NewEncodedConn(nc, MSGP_ENCODER)
	if err != nil {
		panic(err)
	}
	return &NATClient{
		conn: c,
	}
}

func (p *NATClient) Destroy() {
	p.conn.Close()
}

func (p *NATClient) Request(ctx context.Context, subj string, req interface{}, resp interface{}) error {
	return p.conn.RequestWithContext(ctx, SUBJ_PRE+subj, req, &resp)
}

func (p *NATClient) Publish(subj string, data interface{}) error {
	return p.conn.Publish(SUBJ_PRE+subj, data)
}

func (p *NATClient) Subscribe(subj string, handle Handle) error {
	f := func(subj, reply string, req interface{}) {
		go func() {
			subj = strings.TrimPrefix(subj, SUBJ_PRE)
			resp := handle(subj, req)
			if reply != "" {
				if err := p.conn.Publish(reply, resp); err != nil {
					log.Println(err, " reply is ", reply)
				}
			}
		}()
	}

	_, err := p.conn.Subscribe(SUBJ_PRE+subj, f)
	return err
}

func (p *NATClient) QueueSubscribe(subj string, handle Handle) error {
	f := func(subj, reply string, req interface{}) {
		go func() {
			subj = strings.TrimPrefix(subj, SUBJ_PRE)
			resp := handle(subj, req)
			if reply != "" {
				if err := p.conn.Publish(reply, resp); err != nil {
					log.Println(err, " reply is ", reply)
				}
			}
		}()
	}

	_, err := p.conn.QueueSubscribe(SUBJ_PRE+subj, QUEUE_PRE+subj, f)
	return err
}
