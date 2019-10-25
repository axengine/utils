package redis

import (
	"testing"
)

var rds *Rds

func init() {
	rds = New("192.168.8.145:6379", "12345678", 15)
}

func TestMutex_Lock(t *testing.T) {
	mu := rds.NewMutex("abc")
	ok, err := mu.Lock(3)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal(ok)
	}
	defer mu.Unlock()

	{
		_, err := mu.Lock(3)
		if err == nil {
			t.Fatal("error lock")
		}

	}
}

func TestMutex_Set(t *testing.T) {
	err := rds.Set("abc", []byte("abvc"))
	if err != nil {
		t.Fatal(err)
	}

}

func TestMutex_SetNX(t *testing.T) {
	err := rds.SetNX("abc", []byte("abvc"))
	if err != nil {
		t.Fatal(err)
	}

}

func TestMutex_Get(t *testing.T) {
	v, err := rds.Get("abc")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(v))

}
