package discovery

import (
	"encoding/json"
	"go.etcd.io/etcd/client/v2"
	"log"
	"runtime"
	"time"

	"golang.org/x/net/context"
)

type Worker struct {
	Name    string
	IP      string
	KeysAPI client.KeysAPI
}

// workerInfo is the service register information to etcdx
type WorkerInfo struct {
	Name string
	IP   string
	CPU  int
}

func NewWorker(name, IP string, endpoints []string) *Worker {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcdx:", err)
	}

	w := &Worker{
		Name:    name,
		IP:      IP,
		KeysAPI: client.NewKeysAPI(etcdClient),
	}
	return w
}

func (w *Worker) HeartBeat() {
	api := w.KeysAPI

	for {
		info := &WorkerInfo{
			Name: w.Name,
			IP:   w.IP,
			CPU:  runtime.NumCPU(),
		}

		key := "workers/" + w.Name
		value, _ := json.Marshal(info)

		_, err := api.Set(context.Background(), key, string(value), &client.SetOptions{
			TTL: time.Second * 10,
		})
		if err != nil {
			log.Println("Error update workerInfo:", err)
		}
		time.Sleep(time.Second * 3)
	}
}
