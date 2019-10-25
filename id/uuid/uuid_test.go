package uuid

import (
	"fmt"
	"testing"
)

func Test_GetUUID(t *testing.T) {
	fmt.Println(GenUUID())
}

func Test_GetGUID(t *testing.T) {
	fmt.Println(GenGUID())
}

func BenchmarkGetUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenGUID()
	}
}
