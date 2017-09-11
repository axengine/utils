// bimap 是一个简单的线程安全的双向Map，要求k和v都唯一
package bimap

import (
	"sync"
)

type BiMap struct {
	mu    sync.Mutex
	left  map[interface{}]interface{}
	right map[interface{}]interface{}
}

func NewBiMap() *BiMap {
	return &BiMap{
		left:  make(map[interface{}]interface{}),
		right: make(map[interface{}]interface{}),
	}
}

func (p *BiMap) Insert(k, v interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.left[k], p.left[v] = v, k
}

func (p *BiMap) RightInsert(k, v interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.left[v], p.left[k] = k, v
}

func (p *BiMap) Value(k interface{}) interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	v, _ := p.left[k]
	return v
}

func (p *BiMap) RightValue(k interface{}) interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	v, _ := p.right[k]
	return v
}

func (p *BiMap) Exist(k interface{}) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.left[k]
	return ok
}

func (p *BiMap) RightExist(k interface{}) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.right[k]
	return ok
}

func (p *BiMap) Delete(k interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	v, ok := p.left[k]
	delete(p.left, k)
	if ok {
		delete(p.right, v)
	}
}

func (p *BiMap) RightDelete(k interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	v, ok := p.right[k]
	delete(p.right, k)
	if ok {
		delete(p.left, v)
	}
}
