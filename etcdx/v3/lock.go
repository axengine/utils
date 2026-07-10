package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Lock struct {
	cli *clientv3.Client
	key string

	session *concurrency.Session
	mutex   *concurrency.Mutex
}

func NewLock(cli *clientv3.Client, key string) *Lock {
	return &Lock{cli: cli, key: key}
}

// Lock TTL is a lease, which will be automatically renewed, and the program must ensure that Unlock can be called after Lock,
// otherwise other programs cannot obtain the lock, unless the program exits, the network is disconnected, etc., and the lease is not actively renewed
func (l *Lock) Lock(ctx context.Context, ttl int) error {
	// Create a session for lock management
	session, err := concurrency.NewSession(l.cli, concurrency.WithTTL(ttl), concurrency.WithContext(ctx))
	if err != nil {
		return err
	}
	l.session = session

	// Create a lock
	mutex := concurrency.NewMutex(session, l.key)
	if err := mutex.Lock(ctx); err != nil {
		return err
	}
	l.mutex = mutex
	return nil
}

func (l *Lock) Unlock(ctx context.Context) error {
	defer l.session.Close()
	if err := l.mutex.Unlock(ctx); err != nil {
		return err
	}
	return nil
}
