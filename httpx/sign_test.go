package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

var (
	exampleAccessKey    = "exampleAccessKey"
	exampleAccessSecret = "exampleAccessSecret"
	exampleTTL          = int64(60)
	exampleMethod       = HmacSha256Hex
)

func testHttpServer(sign *APISign) {
	http.HandleFunc("/sign", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "GET":
			msg := request.Method + "\n"
			err := sign.Verify(request, "Authorization")
			if err != nil {
				msg += err.Error()
			}
			writer.WriteHeader(200)
			writer.Write([]byte(msg))
		case "POST", "PUT":
			msg := request.Method + "\n"
			err := sign.Verify(request, "Authorization")
			if err != nil {
				msg += err.Error()
			}

			// 读出Body
			bodyData, _ := io.ReadAll(request.Body)
			defer request.Body.Close()

			writer.WriteHeader(200)
			writer.Write([]byte(msg + "\n"))
			writer.Write(bodyData)
		}
	})
	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func TestAPISignAndVerify(t *testing.T) {
	sn := NewAPISign(exampleAccessKey, exampleAccessSecret, exampleTTL, exampleMethod)
	go func() {
		testHttpServer(sn)
	}()
	time.Sleep(time.Second)
	path := "http://localhost:8080/sign?2=2&1=1"
	pointVal := float32(40.5)
	type Add struct {
		Address string `json:"address"`
		No      int
	}
	b := struct {
		Name    string
		Age     int
		Balance float64
		Point   *float32
		Adds    []Add
	}{Name: "nice", Age: 5, Balance: 102.22222, Point: &pointVal, Adds: []Add{{Address: "XX", No: 88}, {Address: "AA", No: 77}}}
	byts, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", path, bytes.NewReader(byts))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	deadline := time.Now().Unix() + sn.TTL

	signature, err := sn.Sign(req, deadline)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s:%s:%d",
		exampleAccessKey, signature, deadline))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	respVal, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respVal))
}
