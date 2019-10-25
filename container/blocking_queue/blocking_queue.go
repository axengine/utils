package blocking_queue

import (
	"container/list"
	"math"
	"sync"
)

// 阻塞队列
type BlockingQueue struct {
	cap  int
	list *list.List
	mu   *sync.RWMutex
	cond *sync.Cond
}

func New(cap int) *BlockingQueue {
	var al BlockingQueue
	al.mu = new(sync.RWMutex)
	al.cond = sync.NewCond(al.mu)
	al.list = list.New()
	al.cap = cap
	if al.cap <= 0 {
		al.cap = math.MaxInt32
	}
	return &al
}

func (p *BlockingQueue) Empty() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.list.Len() == 0
}

func (p *BlockingQueue) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.list.Len()
}

func (p *BlockingQueue) Front() interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.list.Front() != nil {
		return p.list.Front().Value
	}
	return nil
}

func (p *BlockingQueue) Back() interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.list.Back() != nil {
		return p.list.Back().Value
	}
	return nil
}

func (p *BlockingQueue) PushBack(e interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.list.Len() >= p.cap {
		p.cond.Wait()
	}
	p.list.PushBack(e)
	p.cond.Signal()
}

func (p *BlockingQueue) PopFront() interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.list.Len() == 0 {
		p.cond.Wait()
	}
	e := p.list.Front()
	p.cond.Signal()
	return p.list.Remove(e)
}
