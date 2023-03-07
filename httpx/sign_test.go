package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/axengine/utils/sign"
	"io"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"
)

var (
	exampleAccessKey    = "exampleAccessKey"
	exampleAccessSecret = "exampleAccessSecret"
	exampleTTL          = time.Second * 5
	exampleMethod       = HmacSha256Hex
)

func testHttpServer(sign *APISign) {
	http.HandleFunc("/sign", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "GET":
			msg := request.Method
			err := sign.Verify(request, "Authorization")
			if err != nil {
				msg += err.Error()
			}
			writer.WriteHeader(200)
			writer.Write([]byte(msg))
		case "POST", "PUT":
			msg := request.Method
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

func TestAPISign_Verify(t *testing.T) {
	go func() {
		testHttpServer(NewAPISign(exampleAccessKey, exampleAccessSecret, exampleTTL, exampleMethod))
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
	ttl := time.Now().Add(time.Second * 10).Unix()
	timestamp := fmt.Sprintf("%d", ttl)

	//signStr := getMethodSign("/sign", timestamp, map[string]interface{}{"2": 2, "1": 1})
	signStr := jsonSign(timestamp, bytes.NewBuffer(byts))
	req, err := http.NewRequest("POST", path, bytes.NewReader(byts))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Authorization", fmt.Sprintf("%s:%s:%s",
		exampleAccessKey, signStr, timestamp))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	respVal, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respVal))
}

func getMethodSign(path, timestamp string, kv map[string]interface{}) string {
	keys := make([]string, 0)
	for k, _ := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	vals := make([]string, 0)
	for _, k := range keys {
		vals = append(vals, fmt.Sprintf("%s=%v", k, kv[k]))
	}
	str := path + "?" + strings.Join(vals, "&") + timestamp
	fmt.Println(str)
	return sign.HMACSha256Hex([]byte(str), []byte(exampleAccessSecret))
}

func jsonSign(timestamp string, body io.Reader) string {
	b := make(RequestBodyMap)
	byts, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(byts, &b); err != nil {
		panic(err)
	}
	str, _ := b.SortToString("&")

	str = str + timestamp
	return sign.HMACSha256Hex([]byte(str), []byte(exampleAccessSecret))
}

type mockSignStruct struct {
	Txt               string         `json:"txt"`
	Integer           int            `json:"integer"`
	SliceBool         []bool         `json:"sliceBool"`
	SliceString       []string       `json:"sliceString"`
	Float             float64        `json:"float"`
	SliceFloat        []float64      `json:"float64s"`
	SliceCustomStruct []CustomStruct `json:"sliceCustomStruct"`
	CustomStruct      CustomStruct   `json:"customStruct"`
}
type CustomStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestSortToString(t *testing.T) {
	obj := mockSignStruct{
		Txt:         "Mock",
		Integer:     100,
		SliceBool:   []bool{true, false},
		SliceString: []string{"A", "B"},
		Float:       9.99,
		SliceFloat:  []float64{8.88, 9.99},
		SliceCustomStruct: []CustomStruct{
			{
				Name: "SliceCustomStruct",
				Age:  10,
			},
			{
				Name: "SliceCustomStruct",
				Age:  20,
			},
		},
		CustomStruct: CustomStruct{
			Name: "CustomStruct",
			Age:  10,
		},
	}
	str, _ := json.Marshal(obj)
	body := RequestBodyMap{}

	if err := json.Unmarshal(str, &body); err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(str))

	ts := fmt.Sprintf("%d", time.Now().Unix())
	sortStr, _ := body.SortToString("&")

	rawStr := sortStr + ts
	fmt.Println(rawStr)

	signed := SignHash(HmacSha256Hex, []byte(rawStr), []byte("lc3pptr2g2sumgvcbt5yw5g3e0tf8oni"))

	fmt.Println(signed)

	fmt.Println(SignHash(HmacSha256Hex, []byte("openId=8JCV80S92NDYLTCH3DPMPGT1XHKH1155H19CF0TH1632473877"), []byte("lc3pptr2g2sumgvcbt5yw5g3e0tf8oni")))
}
