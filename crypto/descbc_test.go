package crypto

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestDESEnDeCrypt(t *testing.T) {
	key := []byte("12345678")
	plaintext := []byte("hello world!")

	ciphertext, err := DESCBCPCSK5Encrypt(key, plaintext)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(ciphertext))

	plaintextnew, err := DESCBCPCSK5Decrypt(key, ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(plaintextnew, plaintext) {
		t.Fatal("encrypt or decrypt failed")
	}
	fmt.Println(string(plaintextnew))
}
