package crypto

import (
	"encoding/base64"
	"testing"
)

func Benchmark_TripleDesECBEncrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TripleDesECBEncrypt([]byte("Hello World!"), "www.huhutong.com")
	}
}

func Benchmark_TripleDesECBDecrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encryped, _ := base64.StdEncoding.DecodeString("ZXzboHCq6vON+5xjAgzb2Q==")
		TripleDesECBDecrypt(encryped, "www.huhutong.com")
	}
}
