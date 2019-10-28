package natx

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_defaultApp_ReqReply(t *testing.T) {
	app := &DefaultAPP{
		cli: NewNATClient([]string{"http://192.168.8.101:4222", "http://192.168.8.121:4222", "http://192.168.8.141:4222"}),
	}
	testTopic := "xyz"
	app.Subscribe(testTopic, app.defaultHandle)

	resp, err := app.Request(context.Background(), testTopic, &DefaultAPPReq{
		Arg1: 456,
		Arg2: "alimmm",
	})
	if err != nil {
		t.Fatal(err)
	}
	ret := resp.(*DefaultAPPResp)
	fmt.Println("客户端：", ret.Arg1, " ", ret.Arg2)
}

func Test_defaultApp_Publish(t *testing.T) {
	app := &DefaultAPP{
		cli: NewNATClient([]string{"http://192.168.8.101:4222", "http://192.168.8.121:4222", "http://192.168.8.141:4222"}),
	}
	testTopic := "xyz"
	app.Subscribe(testTopic, app.defaultHandle)

	err := app.Publish(testTopic, &DefaultAPPReq{
		Arg1: 456,
		Arg2: "alimmm",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 3)
}
