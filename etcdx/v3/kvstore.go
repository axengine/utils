package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type KvStore struct {
	cli *clientv3.Client

	lease clientv3.Lease
}

func NewKvStore(cli *clientv3.Client) *KvStore {
	return &KvStore{cli: cli, lease: clientv3.NewLease(cli)}
}

func (kv *KvStore) Set(ctx context.Context, key, value string, ttl int) error {
	leaseResp, err := kv.lease.Grant(ctx, int64(ttl))
	if err != nil {
		return err
	}

	// Use the lease id to put the key-value pair into etcd
	putResp, err := kv.cli.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	_ = putResp
	return nil
}

func (kv *KvStore) Get(ctx context.Context, key string) (string, error) {
	// Check if the key-value pair still exists
	getResp, err := kv.cli.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if len(getResp.Kvs) == 0 {
		return "", nil
	}
	return string(getResp.Kvs[0].Value), nil
}
