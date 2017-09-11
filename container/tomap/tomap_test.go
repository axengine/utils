package tomap

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	size := 100000
	tm := New(size)

	for i := 0; i < size; i++ {
		tm.Add(fmt.Sprintf("%d", i), i, time.Second*time.Duration(1))
	}

	time.Sleep(time.Second * 2)
	if len(tm.data) != 0 {
		t.FailNow()
	}

	for i := 0; i < size; i++ {
		tm.Add(fmt.Sprintf("%d", i), i, time.Second*time.Duration(10))
	}

	for i := 0; i < size; i++ {
		tm.ReSet(fmt.Sprintf("%d", i), i+100, time.Second*time.Duration(1))
	}

	//tm.Format()

	time.Sleep(time.Second * 2)
	if len(tm.data) != 0 {
		t.FailNow()
	}

	for i := 0; i < size; i++ {
		tm.Add(fmt.Sprintf("%d", i), i, time.Second*time.Duration(10))
	}
	for i := 0; i < size; i++ {
		tm.Del(fmt.Sprintf("%d", i))
	}
	if len(tm.data) != 0 {
		t.FailNow()
	}
}

func BenchmarkAll(b *testing.B) {
	size := 1000000
	tm := New(size)
	b_add(tm, size, b)
	b_reset(tm, size, b)
	b_del(tm, size, b)
	if len(tm.data) != 0 {
		b.FailNow()
	}
}

func b_add(tm *TOMap, size int, b *testing.B) {
	for i := 0; i < size; i++ {
		tm.Add(fmt.Sprintf("%d", i), i, time.Second*time.Duration(30))
	}

	if len(tm.data) != size {
		b.FailNow()
	}
}

func b_reset(tm *TOMap, size int, b *testing.B) {
	for i := 0; i < size; i++ {
		tm.ReSet(fmt.Sprintf("%d", i), i, time.Second*time.Duration(60))
	}

	if len(tm.data) != size {
		b.FailNow()
	}
}

func b_del(tm *TOMap, size int, b *testing.B) {
	for i := 0; i < size; i++ {
		tm.Del(fmt.Sprintf("%d", i))
	}
}
