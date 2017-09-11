package bitset

import (
	"math/rand"
	"testing"
	"time"
)

func TestBitSet_Add(t *testing.T) {
	bitset := NewBitSet()
	var key uint32 = 123851235
	bitset.Add(key)
	if !bitset.IsExist(key) {
		t.Error("key can not found in the bitset")
	}
	bitset.Del(key)
	if bitset.IsExist(key) {
		t.Error("key should not found in the bitset")
	}

	key = 0
	bitset.Add(key)
	if !bitset.IsExist(0) {
		t.Error("key can not found in the bitset")
	}
}
func BenchmarkBitSet_Add(b *testing.B) {
	bitset := NewBitSet()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < b.N; i++ {
		key := uint32(r.Int31())
		bitset.Add(key)
		ok := bitset.IsExist(key)
		if !ok {
			b.Errorf("第%9d计算,出现数据不存在,key=%9d \n", i, key)
			break
		}

	}

}
func BenchmarkBitSet_IsExist(b *testing.B) {
	bitset := NewBitSet()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var max uint32 = 1024 * 1024 * 1024
	for i := uint32(0); i < max; i++ {
		key := uint32(i)
		bitset.Add(key)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := uint32(r.Int31())
		ok := bitset.IsExist(key)
		if (ok && key >= max) || (!ok && key < max) {
			b.Errorf("第%9d运算,发现存放的key:%9d不符合逻辑,max:%9d,ok=%v \n", i, key, max, ok)
			break
		}
	}
}
