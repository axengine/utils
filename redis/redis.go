package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

type Rds struct {
	pool *redis.Pool
}

func New(dial, password string, db int) *Rds {
	pool := initPool(dial, password, db)
	return &Rds{
		pool: pool,
	}
}

func initPool(dial, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: 120 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", dial)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// Lock 锁住key
func (p *Rds) Lock(key string, ttl int) error {
	conn := p.pool.Get()
	defer conn.Close()
	mu := &Mutex{
		key: key,
	}

	reply, err := redis.String(conn.Do("SET", mu.Key(), mu.Key(), "EX", ttl, "NX"))
	if err != nil {
		return err
	}
	if reply != "OK" {
		return errors.New("xxx" + reply)
	}
	return err
}

// UnLock 解锁key
func (p *Rds) UnLock(key string) (err error) {
	conn := p.pool.Get()
	defer conn.Close()
	mu := &Mutex{
		key: key,
	}
	_, err = conn.Do("del", mu.Key())
	return
}

func (p *Rds) Get(key string) (v []byte, err error) {
	conn := p.pool.Get()
	defer conn.Close()
	v, err = redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		return v, nil
	}
	return
}

func (p *Rds) Set(key string, value []byte) (err error) {
	conn := p.pool.Get()
	defer conn.Close()
	_, err = redis.String(conn.Do("SET", key, value))
	return
}

func (p *Rds) SetNX(k string, v interface{}) (err error) {
	conn := p.pool.Get()
	defer conn.Close()
	_, err = redis.Int64(conn.Do("SETNX", k, v))
	return
}

func (p *Rds) SetEX(k string, expire int, v interface{}) (err error) {
	conn := p.pool.Get()
	defer conn.Close()
	_, err = redis.Bytes(conn.Do("SETEX", k, expire, v))
	return
}

func (p *Rds) Del(k string) (err error) {
	conn := p.pool.Get()
	defer conn.Close()
	_, err = conn.Do("DEL", k)
	return
}

func (p *Rds) Keys() (keys []string, err error) {
	conn := p.pool.Get()
	defer conn.Close()
	keys, err = redis.Strings(conn.Do("KEYS", "*"))
	return
}

func (p *Rds) FlushDB() (err error) {
	conn := p.pool.Get()
	defer conn.Close()
	_, err = conn.Do("FLUSHDB")
	return
}
