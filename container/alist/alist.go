package alist

import (
	"container/list"
	"sync"
)

type AList struct {
	mu   *sync.Mutex
	cond *sync.Cond
	list *list.List
}

func New() *AList {
	var al AList
	al.mu = new(sync.Mutex)
	al.cond = sync.NewCond(al.mu)
	al.list = list.New()
	return &al
}

func (p *AList) Put(e interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.list.PushBack(e)
	p.cond.Signal()
}

func (p *AList) Get() interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.list.Len() == 0 {
		p.cond.Wait()
	}
	e := p.list.Front()
	return p.list.Remove(e)
}
