package natx

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type DefaultAPP struct {
	cli *NATClient
}

type DefaultAPPReq struct {
	Arg1 int
	Arg2 string
}

type DefaultAPPResp struct {
	Arg1 int
	Arg2 string
}

func (app *DefaultAPP) Request(ctx context.Context, subj string, req interface{}) (interface{}, error) {
	var resp DefaultAPPResp
	err := app.cli.conn.RequestWithContext(context.Background(), SUBJ_PRE+subj, req, &resp)
	return &resp, err
}

func (app *DefaultAPP) Subscribe(subj string, handle Handle) error {
	f := func(subj, reply string, req *DefaultAPPReq) {
		go func() {
			subj = strings.TrimPrefix(subj, SUBJ_PRE)
			resp := handle(subj, req)
			if reply != "" {
				if err := app.cli.conn.Publish(reply, resp); err != nil {
					log.Println(err, " reply is ", reply)
				}
			}
		}()
	}

	_, err := app.cli.conn.Subscribe(SUBJ_PRE+subj, f)
	return err
}

func (app *DefaultAPP) QueueSubscribe(subj string, handle Handle) error {
	f := func(subj, reply string, req *DefaultAPPReq) {
		go func() {
			subj = strings.TrimPrefix(subj, SUBJ_PRE)
			resp := handle(subj, req)
			if reply != "" {
				if err := app.cli.conn.Publish(reply, resp); err != nil {
					log.Println(err, " reply is ", reply)
				}
			}
		}()
	}

	_, err := app.cli.conn.QueueSubscribe(SUBJ_PRE+subj, QUEUE_PRE+subj, f)
	return err
}

func (app *DefaultAPP) Publish(subj string, data interface{}) error {
	return app.cli.conn.Publish(SUBJ_PRE+subj, data)
}

func (app *DefaultAPP) defaultHandle(subj string, req interface{}) interface{} {
	fmt.Println("服务端：", subj)
	r := req.(*DefaultAPPReq)
	fmt.Println("服务端：", r.Arg1, " ", r.Arg2)
	return &DefaultAPPResp{
		Arg1: 123,
		Arg2: "hello world",
	}
}
