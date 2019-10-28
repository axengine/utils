package codec

import (
	"fmt"
	"testing"
)

func Test_ENCODE_BIT_NOT(t *testing.T) {
	str := "hello world"
	fmt.Println("length is ", len(str))
	buff := []byte(str)
	Encode(ENCODE_BIT_NOT, buff)
	fmt.Println("after encode is:", buff)
	Decode(ENCODE_BIT_NOT, buff)
	if string(buff) != str {
		t.Fatal("not equal")
	}
}

func Test_ENCODE_BYTE_RVS(t *testing.T) {
	str := "hello world"
	fmt.Println("length is ", len(str))
	buff := []byte(str)
	Encode(ENCODE_BYTE_RVS, buff)
	fmt.Println("after encode is:", buff)
	Decode(ENCODE_BYTE_RVS, buff)
	if string(buff) != str {
		t.Fatal("not equal")
	}
}

func Test_ENCODE_LOOP_XOR(t *testing.T) {
	str := "hello world"
	fmt.Println("length is ", len(str))
	buff := []byte(str)
	Encode(ENCODE_LOOP_XOR, buff)
	fmt.Println("after encode is:", buff)
	Decode(ENCODE_LOOP_XOR, buff)
	if string(buff) != str {
		t.Fatal("not equal")
	}
}

func Test_ENCODE_LOOP_XOR_N(t *testing.T) {
	str := "hello world"
	fmt.Println("length is ", len(str))
	buff := []byte(str)
	n := 10240
	for i := 0; i < n; i++ {
		Encode(ENCODE_LOOP_XOR, buff)
	}

	for i := 0; i < n; i++ {
		Decode(ENCODE_LOOP_XOR, buff)
	}
	fmt.Printf("after  %v\n", string(buff))
}

func Benchmark_ENCODE_LOOP_XOR(b *testing.B) {
	buff := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		Encode(ENCODE_LOOP_XOR, buff)
	}
}
