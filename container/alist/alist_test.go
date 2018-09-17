package alist

import (
	"testing"
)

func Test_Base(t *testing.T) {
	al := New()
	for i := 0; i < 10; i++ {
		al.Put(i)
	}
	var count int
	for {
		v := al.Get().(int)
		t.Log("got ", v)
		count = count + 1
		if count == 10 {
			t.Log("success")
			return
		}
	}
	t.Error("failed")
}

func Benchmark_Put(b *testing.B) {
	b.StopTimer()
	al := New()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		al.Put(i)
	}
}

func Benchmark_PutGet(b *testing.B) {
	b.StopTimer()
	al := New()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		al.Put(i)
	}
	for i := 0; i < b.N; i++ {
		if i != al.Get().(int) {
			b.Error("got err value")
		}
	}
}
