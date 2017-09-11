// stack 线程安全的堆栈
package stack

import (
	"container/list"
	"sync"
)

type Stack struct {
	mutex sync.Mutex
	list  *list.List
}

func New() *Stack {
	return &Stack{
		list: list.New(),
	}
}

func (p *Stack) Init() {
	p.list = p.list.Init()
}

func (p *Stack) Push(value interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.list.PushBack(value)
}

func (p *Stack) PushFront(value interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.list.PushFront(value)
}

func (p *Stack) Pop() interface{} {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	e := p.list.Back()
	if e != nil {
		p.list.Remove(e)
		return e.Value
	}
	return nil
}

func (p *Stack) Peak() interface{} {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	e := p.list.Back()
	if e != nil {
		return e.Value
	}

	return nil
}

func (p *Stack) Len() int {
	return p.list.Len()
}

func (p *Stack) Empty() bool {
	return p.list.Len() == 0
}
