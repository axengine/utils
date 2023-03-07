package referral

import (
	"errors"
	"fmt"
	"testing"
)

func init() {
	//SetBase("ETN6BGQF7P5IK3MUAR4HV8S2DJZX9WYL")
	//SetLength(8)
}

func TestEncode(t *testing.T) {
	for i := 1; i < 10000000; i++ {
		code := Encode(uint64(i))
		num := Decode(code)
		if num != uint64(i) {
			t.Fatal(errors.New("error"))
		}
		if i%100 == 0 {
			fmt.Println(code, num)
		}
	}
}

func TestDecode(t *testing.T) {
	code := Encode(1)
	fmt.Println(code)
	id := Decode(code)
	fmt.Println(id)

}

func TestDecode2(t *testing.T) {
	code := "NIBCN6"
	id := Decode(code)
	fmt.Println(id)
}

func TestEncode2(t *testing.T) {
	id := 1
	code := Encode(uint64(id))
	fmt.Println("code", code)
}

func TestDumpKey(t *testing.T) {
	keys := make(map[string]struct{}, 100000000)
	for i := 1; i < 1000000000; i++ {
		code := Encode(uint64(i))
		_, ok := keys[code]
		if ok {
			t.Fatal("dump key")
		}
		keys[code] = struct{}{}
		if Decode(code) != uint64(i) {
			t.Fatal("decode exception")
		}
	}
}
