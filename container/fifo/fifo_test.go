package fifo

import (
	"fmt"
	"testing"
)

func TestFIFO_Enqueue(t *testing.T) {
	f := New(3)
	for i := 0; i < 3; i++ {
		b := f.Enqueue(i)
		if !b {
			t.Fatal("enqueue false i:", i)
		}
		fmt.Println("h:", f.front)
		fmt.Println("t:", f.rear)
		fmt.Println("d:", f.queue)
	}

	b := f.Enqueue(1)
	if b {
		t.Fatal("not full check:")
	}
	fmt.Println(b)
}

func TestFIFO_Dequeue(t *testing.T) {
	f := New(3)
	for i := 0; i < 3; i++ {
		b := f.Enqueue(i)
		if !b {
			t.Fatal("enqueue false i:", i)
		}
	}
	fmt.Println("h:", f.front)
	fmt.Println("t:", f.rear)
	fmt.Println("d:", f.queue)
	for i := 0; i < 3; i++ {
		x := f.Dequeue()
		if x == nil {
			t.Fatal("dequeue failed")
		}
		fmt.Println(x)
	}
	x := f.Dequeue()
	if x != nil {
		t.Fatal("dequeue failed")
	}
	fmt.Println(x)
}

func BenchmarkFIFO_Put(b *testing.B) {
	var max = int(10000000)
	f := New(max)
	for i := 0; i < b.N; i++ {
		ok := f.Enqueue(i)
		if i < max && !ok {
			b.Fatal(ok)
		}
	}
}

func BenchmarkFIFO_Get(b *testing.B) {
	var max = int(10000000)
	f := New(max)
	for i := 0; i < max; i++ {
		ok := f.Enqueue(i)
		if !ok {
			b.Fatal(ok)
		}
	}

	for i := 0; i < b.N; i++ {
		x := f.Dequeue()
		if i < max && x == nil {
			b.Fatal("dequeue failed")
		}
	}
}
