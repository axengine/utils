package random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func RandDigits(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	codes := make([]string, 0)
	for i := 0; i < length; i++ {
		codes = append(codes, fmt.Sprintf("%d", r.Intn(10)))
	}
	return strings.Join(codes, "")
}

func RandAlphaAndDigits(length int, lower ...bool) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result[i] = bytes[r.Intn(len(bytes))]
	}
	if len(lower) > 0 && lower[0] {
		return strings.ToLower(string(result))
	}
	return string(result)
}
