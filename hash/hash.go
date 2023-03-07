package hash

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
)

func Keccak256Hex(in []byte) string {
	sha := sha3.NewLegacyKeccak256()
	sha.Write(in)
	sum := sha.Sum(nil)
	return hex.EncodeToString(sum)
}

func Keccak256Base64(in []byte) string {
	sha := sha3.NewLegacyKeccak256()
	sha.Write(in)
	sum := sha.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}

func MD5Lower(b []byte) string {
	sha := md5.New()
	sha.Write(b)
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func MD5Upper(b []byte) string {
	sha := md5.New()
	sha.Write(b)
	return fmt.Sprintf("%X", sha.Sum(nil))
}
