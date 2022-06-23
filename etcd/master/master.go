package master

import (
	"log"
	"time"

	"go.etcd.io/etcd/client/v2"
	"golang.org/x/net/context"
)

type Master struct {
	key              string         //as key
	id               string         //as value
	kapi             client.KeysAPI //etcd kapi
	defaultTTL       time.Duration  //
	defaultHeartbeat time.Duration  //
}

func NewMaster(key, id string, endpoints []string) *Master {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
		return nil
	}
	return &Master{
		kapi:             client.NewKeysAPI(etcdClient),
		key:              key,
		id:               id,
		defaultTTL:       time.Second * 5,
		defaultHeartbeat: time.Second * 3,
	}
}

func (m *Master) heartbeat() {
	tk := time.NewTicker(m.defaultHeartbeat)
	for {
		select {
		case <-tk.C:
			log.Println("heartbeat...k=", m.key, " v=", m.id)
			if _, err := m.kapi.Set(context.Background(), m.key, m.id, &client.SetOptions{
				TTL: m.defaultTTL,
			}); err != nil {
				log.Println("heartbeat err:", err)
			}
		}
	}
}

// Apply for master block until success or occur an error
func (m *Master) Apply() error {
	setOptions := &client.SetOptions{
		PrevExist: client.PrevNoExist,
		TTL:       m.defaultTTL,
	}
	for {
		resp, err := m.kapi.Set(context.Background(), m.key, m.id, setOptions)
		if err == nil {
			go m.heartbeat()
			return nil
		}

		e, ok := err.(client.Error)
		if !ok {
			return err
		}
		if e.Code != client.ErrorCodeNodeExist {
			return err
		}
		resp, err = m.kapi.Get(context.Background(), m.key, nil)
		if err != nil {
			return err
		}
		log.Println("Apply failed,watch ", m.key)
		watcherOptions := &client.WatcherOptions{
			AfterIndex: resp.Index,
			Recursive:  false,
		}
		watcher := m.kapi.Watcher(m.key, watcherOptions)
		for {
			resp, err = watcher.Next(context.Background())
			if err != nil {
				return err
			}
			if resp.Action == "set" {
				log.Println("watching event:", resp.Action, " key:", m.key, " value:", resp.Node.Value)
			}
			if resp.Action == "delete" || resp.Action == "expire" {
				log.Println("watching event:", resp.Action, " key:", m.key, " value:", resp.PrevNode.Value)
				break
			}
		}
	}
}
