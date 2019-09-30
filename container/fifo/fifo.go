package fifo

import (
	"sync"
)

type FIFO struct {
	front int
	rear  int
	size  int
	mu    sync.Mutex
	queue []interface{}
}

func New(size int) *FIFO {
	return &FIFO{
		front: 0,
		rear:  0,
		size:  size,
		queue: make([]interface{}, size),
	}
}

func (f *FIFO) Enqueue(i interface{}) bool {
	f.mu.Lock()
	if (f.rear+1)%(f.size+1) == f.front {
		f.mu.Unlock()
		return false
	}

	//if f.rear == f.front {
	//	f.queue[f.front] = i
	//}
	f.queue[f.rear] = i
	f.rear = (f.rear + 1) % (f.size + 1)
	f.mu.Unlock()
	return true
}

func (f *FIFO) Dequeue() interface{} {
	f.mu.Lock()
	if f.rear == f.front {
		f.mu.Unlock()
		return nil
	}
	x := f.queue[f.front]
	f.front = (f.front + 1) % (f.size + 1)
	f.mu.Unlock()
	return x
}
