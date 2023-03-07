package sign

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

// HMACSha1B64 HMAC&SHA1->base64
func HMACSha1B64(raw []byte, key []byte) string {
	hmacHash := hmac.New(sha1.New, key)
	hmacHash.Write(raw)
	return base64.StdEncoding.EncodeToString(hmacHash.Sum(nil))
}

// HMACSha256B64 HMAC&SHA256->base64
func HMACSha256B64(raw []byte, key []byte) string {
	hmacHash := hmac.New(sha256.New, key)
	hmacHash.Write(raw)
	return base64.StdEncoding.EncodeToString(hmacHash.Sum(nil))
}

// HMACSha512B64 HMAC&SHA512->base64
func HMACSha512B64(raw []byte, key []byte) string {
	hmacHash := hmac.New(sha512.New, key)
	hmacHash.Write(raw)
	return base64.StdEncoding.EncodeToString(hmacHash.Sum(nil))
}

// HMACSha1Hex HMAC&SHA1->hex
func HMACSha1Hex(raw []byte, key []byte) string {
	hmacHash := hmac.New(sha1.New, key)
	hmacHash.Write(raw)
	return hex.EncodeToString(hmacHash.Sum(nil))
}

// HMACSha256Hex HMAC&SHA256->hex
func HMACSha256Hex(raw []byte, key []byte) string {
	hmacHash := hmac.New(sha256.New, key)
	hmacHash.Write(raw)
	return hex.EncodeToString(hmacHash.Sum(nil))
}

// HMACSha512Hex HMAC&SHA512->hex
func HMACSha512Hex(raw []byte, key []byte) string {
	hmacHash := hmac.New(sha512.New, key)
	hmacHash.Write(raw)
	return hex.EncodeToString(hmacHash.Sum(nil))
}
