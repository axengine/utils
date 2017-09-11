// github.com/alex023/basekit/blob/master/counter.go
package count

import (
	"sync"
	"sync/atomic"
)

// Counter counter is a multi-thread safe counters
type Counter struct {
	mut     sync.Mutex
	currNum int64 //当前数量
	maxNum  int64 //最大数量
}

// AddOne 在原内部计数基础上，+1。
func (c *Counter) AddOne() int {
	new := atomic.AddInt64(&c.currNum, 1)

	c.mut.Lock()
	if c.maxNum < new {
		c.maxNum = new
	}
	c.mut.Unlock()

	return int(c.currNum)
}

// DecOne 在原内部计数基础上，-1。
func (c *Counter) DecOne() int {
	return int(atomic.AddInt64(&c.currNum, -1))

}

// Current 获取当前内部计数结果。
func (c *Counter) Current() int {
	return int(atomic.LoadInt64(&c.currNum))
}

// MaxNum 计数器生存周期内，最大的计数。
func (c *Counter) MaxNum() int {
	return int(atomic.LoadInt64(&c.maxNum))
}

// Reset 重置计数器
func (c *Counter) Reset() {
	c.mut.Lock()
	c.currNum = 0
	c.maxNum = 0
	c.mut.Unlock()
}

// NewCounter counter constructor
func NewCounter() *Counter {
	return &Counter{}
}
