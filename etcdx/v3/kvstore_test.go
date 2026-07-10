package etcd

import (
	"context"
	"testing"
	"time"
)

func TestGetSet(t *testing.T) {
	kvstore := NewKvStore(_cli_)
	if err := kvstore.Set(context.Background(), "key", "value", 3); err != nil {
		t.Fatal(err)
	}

	if value, err := kvstore.Get(context.Background(), "key"); err != nil {
		t.Fatal(err)
	} else if value != "value" {
		t.Fatal("not equal")
	}
}

func TestSetExpire(t *testing.T) {
	kvstore := NewKvStore(_cli_)
	if err := kvstore.Set(context.Background(), "key", "value", 3); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 4)

	if value, err := kvstore.Get(context.Background(), "key"); err != nil {
		t.Fatal(err)
	} else if value == "value" {
		t.Fatal("not expire:", value)
	}
}
