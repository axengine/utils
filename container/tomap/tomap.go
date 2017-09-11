// timeout map,设置key的过期时间，超时后key被删除
// 使用定时器为临时方案，后续会参考REDIS的TTL机制修改
package tomap

import (
	"log"
	"sync"
	"time"

	timer "github.com/axengine/utils/time/timer_minheap"
)

type TOMap struct {
	mu   sync.RWMutex
	data map[string]*node
	tm   *timer.Timer
}

type node struct {
	v  interface{}
	td *timer.TimerData
}

func New(cap int) *TOMap {
	return &TOMap{
		data: make(map[string]*node),
		tm:   timer.NewTimer(cap),
	}
}

// Add 向TOM加入KV并设置K的超时时间
// 超时后调用超时函数删除该K
func (p *TOMap) Add(key string, v interface{}, expire time.Duration) {
	n := new(node)
	n.v = v
	td := p.tm.Add(expire, p.toHandle, key)
	n.td = td

	p.mu.Lock()
	p.data[key] = n
	p.mu.Unlock()
}

// ReSet 在TOM中重新设置K的V和超时时间
func (p *TOMap) ReSet(key string, v interface{}, expire time.Duration) {
	n, ok := p.data[key]
	if ok {
		p.tm.Set(n.td, expire)
		n.v = v
	} else {
		n := new(node)
		n.v = v
		td := p.tm.Add(expire, p.toHandle, key)
		n.td = td
		p.mu.Lock()
		p.data[key] = n
		p.mu.Unlock()
	}
}

// Del 在TOM中删除K
func (p *TOMap) Del(key string) {
	n, ok := p.data[key]
	if ok {
		p.tm.Del(n.td)
		p.mu.Lock()
		delete(p.data, key)
		p.mu.Unlock()
	}
}

// Get 从TOM中获取K的V
func (p *TOMap) Get(key string) interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	n, ok := p.data[key]
	if ok {
		return n.v
	}
	return nil
}

func (p *TOMap) toHandle(i interface{}) {
	key, ok := i.(string)
	if !ok {
		return
	}
	p.mu.Lock()
	delete(p.data, key)
	p.mu.Unlock()
}

func (p *TOMap) Format() {
	for k, v := range p.data {
		log.Println("k=", k, " v=", v)
	}
}
