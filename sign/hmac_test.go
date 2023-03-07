package sign

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	hmacHash := hmac.New(sha1.New, []byte("xxx"))
	hmacHash.Write([]byte("hello world"))
	fmt.Println(base64.StdEncoding.EncodeToString(hmacHash.Sum(nil)))
}
