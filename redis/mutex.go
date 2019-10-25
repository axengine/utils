package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

type Mutex struct {
	key  string
	conn redis.Conn
}

// NewMutex
// 返回互斥量Mutex
func (p *Rds) NewMutex(key string) *Mutex {
	return &Mutex{
		key:  key,
		conn: p.pool.Get(),
	}
}

func (mu *Mutex) Key() string {
	return "lock_" + mu.key
}

func (mu *Mutex) Lock(ttl int64) (bool, error) {
	conn := mu.conn
	reply, err := redis.String(conn.Do("SET", mu.Key(), mu.Key(), "EX", ttl, "NX"))
	if err != nil {
		conn.Close()
		return false, errors.Wrap(err, "lock failed")
	}

	if reply != "OK" {
		conn.Close()
		return false, errors.New("lock failed " + reply)
	}
	return true, nil
}

func (mu *Mutex) Unlock() (err error) {
	conn := mu.conn
	_, err = conn.Do("del", mu.Key())
	conn.Close()
	return
}
