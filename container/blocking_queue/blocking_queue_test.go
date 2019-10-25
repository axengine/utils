package blocking_queue

import (
	"fmt"
	"testing"
	"time"
)

func Test_Base(t *testing.T) {
	bqueue := New(0)
	for i := 0; i < 10; i++ {
		bqueue.PushBack(i)
	}

	for i := 0; i < 10; i++ {
		v := bqueue.PopFront()
		if i != v.(int) {
			t.Error("not equal")
		}
	}
}

func Benchmark_Put(b *testing.B) {
	al := New(0)
	for i := 0; i < b.N; i++ {
		al.PushBack(i)
	}

	for i := 0; i < b.N; i++ {
		al.PopFront()
	}
}

func Test_GetPut(t *testing.T) {
	al := New(0)
	go func() {
		v := al.PopFront()
		fmt.Println("1 got ", v.(int), " at ", time.Now())
	}()
	go func() {
		v := al.PopFront()
		fmt.Println("2 got ", v.(int), " at ", time.Now())
	}()
	time.Sleep(time.Second)

	al.PushBack(1)
	al.PushBack(2)
	time.Sleep(time.Second * 3)
}

func Test_PutGet(t *testing.T) {
	al := New(5)

	time.AfterFunc(time.Second*3, func() {
		for {
			v := al.PopFront()
			fmt.Println("got ", v.(int), " at ", time.Now())
		}
	})
	for i := 0; i < 10; i++ {
		al.PushBack(i)
		fmt.Println("put ", i, " at ", time.Now())
	}

	time.Sleep(time.Second * 10)
}
