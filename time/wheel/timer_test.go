package timer

import (
	"fmt"
	//	"sync/atomic"
	"testing"
	"time"
)

//var tt *Timer

func fn(t int, arg interface{}) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("type:", t, " arg:", arg)
}

func TestTimer(t *testing.T) {
	timer := New(time.Millisecond * 10)
	//	tt = timer
	//fmt.Println(timer)
	timer.NewTimer(time.Millisecond*time.Duration(1000*1), fn, 14, 15)
	go timer.Start()
	T.NewTimer(time.Millisecond*time.Duration(3000*1), fn, 12, 13)
	time.Sleep(time.Second * 100)

}
