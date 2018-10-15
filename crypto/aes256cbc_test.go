package crypto

import (
	"encoding/base64"
	"testing"
)

var (
	srcStr = `a8971729fbc199fb3459529cebcd8704791fc699d88ac89284f23ff8e7fca7d6`
	dstStr = `ns6Mi2GyaFJJ8Z3HiUwLs3kjShAqQepSyKRFzNv1FViXgOpccwFH6Gab1MSkyZ25OSoCWuQV1a5c+pHG/ZnJnCcoGqfPbkWeMcLpaLYhh0M=`
	ivStr  = `JqCvx/OxR/MN4REmBGDJxQ==`
	keyStr = `69u92Jg5SBWOgH41oB0tKY5rzTIrsjhu`
)

func TestAES256_CBC_PKCS0Encrpt(t *testing.T) {
	src := []byte(srcStr)
	iv, _ := base64.StdEncoding.DecodeString(ivStr)
	key := []byte(keyStr)
	dst, err := AES256_CBC_PKCS0Encrpt(src, iv, key)
	if err != nil {
		t.Error(err)
	}

	if base64.StdEncoding.EncodeToString(dst) != dstStr {
		t.Error("error encrpt ", base64.StdEncoding.EncodeToString(dst), "\n", dst)
	}
}

func TestAES256_CBC_PKCS0Decrpt(t *testing.T) {
	src, _ := base64.StdEncoding.DecodeString(dstStr)
	iv, _ := base64.StdEncoding.DecodeString(ivStr)
	key := []byte(keyStr)
	dst, err := AES256_CBC_PKCS0Decrpt(src, iv, key)
	if err != nil {
		t.Error(err)
	}
	if srcStr != string(dst) {
		t.Error("error decrpt ", string(dst), "\n", dst)
	}
}
