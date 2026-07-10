package cryptox

import (
	"encoding/base64"
	"fmt"
	"testing"
)

var (
	srcStr = `964432afc9c34739a06d6b877b7033b0217e96fa3e2db71896e708fd2bbc5a35`
	dstStr = `ns6Mi2GyaFJJ8Z3HiUwLs3kjShAqQepSyKRFzNv1FViXgOpccwFH6Gab1MSkyZ25OSoCWuQV1a5c+pHG/ZnJnCcoGqfPbkWeMcLpaLYhh0M=`
	ivStr  = `ri34iHc5LzgiWAhw`
	keyStr = `dZgfIU0XsfRzUFbOVRI39LSytTXs4Mvs`
)

func TestAES256_CBC_PKCS0Encrpt(t *testing.T) {
	src := []byte(srcStr)
	//iv, _ := base64.StdEncoding.DecodeString(ivStr)
	iv := []byte(ivStr)
	key := []byte(keyStr)
	dst, err := AES256CBCPKCS0Encrypt(src, iv, key)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(dst))
	if base64.StdEncoding.EncodeToString(dst) != dstStr {
		//t.Error("error encrpt ", base64.StdEncoding.EncodeToString(dst), "\n", dst)
	}
}

func TestAES256_CBC_PKCS0Decrpt(t *testing.T) {
	src, _ := base64.StdEncoding.DecodeString(dstStr)
	iv, _ := base64.StdEncoding.DecodeString(ivStr)
	key := []byte(keyStr)
	dst, err := AES256CBCPKCS0Decrypt(src, iv, key)
	if err != nil {
		t.Error(err)
	}
	if srcStr != string(dst) {
		t.Error("error decrpt ", string(dst), "\n", dst)
	}
}
