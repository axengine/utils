// timer 基于时间轮的定时器服务，精度为10ms
// refer:http://blog.csdn.net/yueguanghaidao/article/details/46290539
// 修改内容：为定时器增加类型和参数属性，修改回调函数类型，更通用
// 使用方法
//	tm = timer.New(time.Millisecond * 10)
//	go tm.Start()
//	tm.NewTimer(time.Second*60, ProcTimer, TIMER_TYPE_DEFAULT, "")
//  func  ProcTimer(tp int, arg interface{}){switch tp {	case TIMER_TYPE_DEFAULT:}}

package timer

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

//referer https://github.com/cloudwu/skynet/blob/master/skynet-src/skynet_timer.c

const (
	TIME_NEAR_SHIFT  = 8
	TIME_NEAR        = 1 << TIME_NEAR_SHIFT
	TIME_LEVEL_SHIFT = 6
	TIME_LEVEL       = 1 << TIME_LEVEL_SHIFT
	TIME_NEAR_MASK   = TIME_NEAR - 1
	TIME_LEVEL_MASK  = TIME_LEVEL - 1
)

type Timer struct {
	near [TIME_NEAR]*list.List     //256
	t    [4][TIME_LEVEL]*list.List //64
	sync.Mutex
	time uint32
	tick time.Duration
	quit chan struct{}
}

type Node struct {
	expire uint32
	t      int
	arg    interface{}
	f      func(int, interface{})
}

var T *Timer

func init() {
	timer := New(time.Millisecond * 10)
	T = timer
	//fmt.Println("start timer server...")
	go T.Start()
}

func (n *Node) String() string {
	return fmt.Sprintf("Node:expire,%d,t:%d,arg:%s", n.expire, n.t, n.arg)
}

func New(d time.Duration) *Timer {
	t := new(Timer)
	t.time = 0
	t.tick = d
	t.quit = make(chan struct{})

	var i, j int
	for i = 0; i < TIME_NEAR; i++ {
		t.near[i] = list.New()
	}

	for i = 0; i < 4; i++ {
		for j = 0; j < TIME_LEVEL; j++ {
			t.t[i][j] = list.New()
		}
	}

	return t
}

func (t *Timer) addNode(n *Node) {
	expire := n.expire
	current := t.time
	if (expire | TIME_NEAR_MASK) == (current | TIME_NEAR_MASK) {
		t.near[expire&TIME_NEAR_MASK].PushBack(n)
	} else {
		var i uint32
		var mask uint32 = TIME_NEAR << TIME_LEVEL_SHIFT
		for i = 0; i < 3; i++ {
			if (expire | (mask - 1)) == (current | (mask - 1)) {
				break
			}
			mask <<= TIME_LEVEL_SHIFT
		}

		t.t[i][(expire>>(TIME_NEAR_SHIFT+i*TIME_LEVEL_SHIFT))&TIME_LEVEL_MASK].PushBack(n)
	}

}

func (t *Timer) NewTimer(d time.Duration, f func(int, interface{}), tp int, arg interface{}) *Node {
	n := new(Node)
	n.f = f
	n.t = tp
	n.arg = arg
	t.Lock()
	n.expire = uint32(d/t.tick) + t.time
	t.addNode(n)
	t.Unlock()
	return n
}

func (t *Timer) String() string {
	return fmt.Sprintf("Timer:time:%d, tick:%s", t.time, t.tick)
}

func dispatchList(front *list.Element) {
	for e := front; e != nil; e = e.Next() {
		node := e.Value.(*Node)
		go node.f(node.t, node.arg)
	}
}

func (t *Timer) moveList(level, idx int) {
	vec := t.t[level][idx]
	front := vec.Front()
	vec.Init()
	for e := front; e != nil; e = e.Next() {
		node := e.Value.(*Node)
		t.addNode(node)
	}
}

func (t *Timer) shift() {
	t.Lock()
	var mask uint32 = TIME_NEAR
	t.time++
	ct := t.time
	if ct == 0 {
		t.moveList(3, 0)
	} else {
		time := ct >> TIME_NEAR_SHIFT
		var i int = 0
		for (ct & (mask - 1)) == 0 {
			idx := int(time & TIME_LEVEL_MASK)
			if idx != 0 {
				t.moveList(i, idx)
				break
			}
			mask <<= TIME_LEVEL_SHIFT
			time >>= TIME_LEVEL_SHIFT
			i++
		}
	}
	t.Unlock()
}

func (t *Timer) execute() {
	t.Lock()
	idx := t.time & TIME_NEAR_MASK
	vec := t.near[idx]
	if vec.Len() > 0 {
		front := vec.Front()
		vec.Init()
		t.Unlock()
		// dispatch_list don't need lock
		dispatchList(front)
		return
	}

	t.Unlock()
}

func (t *Timer) update() {
	// try to dispatch timeout 0 (rare condition)
	t.execute()

	// shift time first, and then dispatch timer message
	t.shift()

	t.execute()

}

func (t *Timer) Start() {
	tick := time.NewTicker(t.tick)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			t.update()
		case <-t.quit:
			return
		}
	}
}

func (t *Timer) Stop() {
	close(t.quit)
}
