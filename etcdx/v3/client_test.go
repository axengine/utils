package etcd

import (
	"os"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var _cli_ *clientv3.Client

func TestMain(m *testing.M) {
	cli, err := NewClient("./certs/ca.pem", "./certs/client.pem", "./certs/client-key.pem",
		"root", "000000", []string{"https://192.168.1.51:2379", "https://192.168.1.53:2379", "https://192.168.1.105:2379"})
	if err != nil {
		panic(err)
	}
	_cli_ = cli
	os.Exit(m.Run())
}

func TestNewClient(t *testing.T) {
	_, err := NewClient("./certs/ca.pem", "./certs/client.pem", "./certs/client-key.pem",
		"root", "000000", []string{"https://192.168.1.51:2379", "https://192.168.1.53:2379", "https://192.168.1.105:2379"})
	if err != nil {
		t.Fatal(err)
	}
	// Can be different from TLS and passwords
	_, err = NewClient("", "", "", "", "", []string{"http://127.0.0.1:2379"})
	if err != nil {
		t.Fatal(err)
	}
}
