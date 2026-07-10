package etcd

import (
	"context"
	"testing"
)

func TestLockUnlock(t *testing.T) {
	mutex := NewLock(_cli_, "/notary/txhash/")
	if err := mutex.Lock(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
	t.Log("lock success")
	if err := mutex.Unlock(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Log("unlock success")
}
