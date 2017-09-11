// rdslock 提供一种rds悲观锁
package lock

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type RdsLock struct {
	conn redis.Conn
}

func New(dial, pswd string, db int) *RdsLock {
	o1 := redis.DialPassword(pswd)
	o2 := redis.DialDatabase(db)
	c, err := redis.Dial("tcp", dial, o1, o2)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &RdsLock{
		conn: c,
	}
}

type Lock struct {
	resource string
	token    string
	conn     redis.Conn
	timeout  int
}

func (p *RdsLock) TryLock(resource string, token string, timeout int) (lock *Lock, err error) {
	lock = &Lock{resource, token, p.conn, timeout}
	ok, err := lock.tryLock()
	if !ok || err != nil {
		lock = nil
	}
	return
}

func (lock *Lock) tryLock() (ok bool, err error) {
	_, err = redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int(lock.timeout), "NX"))
	if err == redis.ErrNil {
		// The lock was not successful, it already exists.
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (lock *Lock) Unlock() (err error) {
	_, err = lock.conn.Do("del", lock.key())
	return
}

func (lock *Lock) key() string {
	return fmt.Sprintf("redislock:%s", lock.resource)
}

func (lock *Lock) AddTimeout(ex_time int64) (ok bool, err error) {
	ttl_time, err := redis.Int64(lock.conn.Do("TTL", lock.key()))
	if err != nil {
		return
	}
	if ttl_time > 0 {
		_, err := redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int(ttl_time+ex_time)))
		if err == redis.ErrNil {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

//func main() {
//	rdsLock := New("192.168.1.63:6379", "12345678", 14)
//	if rdsLock == nil {
//		log.Fatal("rdsLock is nil")
//	}

//	lock, ok, err := rdsLock.TryLock("iamkey", "iamvalue", 10)
//	if err != nil {
//		log.Fatal("Error while attempting lock")
//	}
//	if !ok {
//		log.Fatal("Lock")
//	}
//	lock.AddTimeout(10)

//	time.Sleep(time.Duration(10) * time.Second)
//	fmt.Println("end")
//	defer lock.Unlock()
//}
