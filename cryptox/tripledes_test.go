package cryptox

import (
	"encoding/base64"
	"testing"
)

func Test_TripleDesECBEncrypt(t *testing.T) {
	out, err := TripleDesECBEncrypt([]byte("Hello World!"), "www.huhutong.com")
	if err != nil {
		t.Error(err)
	}
	dst := base64.StdEncoding.EncodeToString(out)
	if dst != "ZXzboHCq6vON+5xjAgzb2Q==" {
		t.Error("加密失败")
	}
}

func Test_TripleDesECBDecrypt(t *testing.T) {
	encryped, _ := base64.StdEncoding.DecodeString("ZXzboHCq6vON+5xjAgzb2Q==")
	out, err := TripleDesECBDecrypt(encryped, "www.huhutong.com")
	if err != nil {
		t.Error(err)
	}
	if string(out) != "Hello World!" {
		t.Error("解密失败")
	}
}
