package cryptox

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestEncryptDecryptAES256GCM(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	plaintext := []byte("Hello, AES-256-GCM!")

	ciphertext, err := EncryptAES256GCM(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := DecryptAES256GCM(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Expected %s, got %s", plaintext, decrypted)
	}
}

func TestEncryptAES256GCM_InvalidKey(t *testing.T) {
	key := make([]byte, 16)
	plaintext := []byte("test")

	_, err := EncryptAES256GCM(plaintext, key)
	if err == nil {
		t.Error("Expected error for invalid key length")
	}
}

func TestDecryptAES256GCM_InvalidKey(t *testing.T) {
	key := make([]byte, 16)
	ciphertext := []byte("test")

	_, err := DecryptAES256GCM(ciphertext, key)
	if err == nil {
		t.Error("Expected error for invalid key length")
	}
}

func TestDecryptAES256GCM_ShortCiphertext(t *testing.T) {
	key := make([]byte, 32)
	ciphertext := []byte("short")

	_, err := DecryptAES256GCM(ciphertext, key)
	if err == nil {
		t.Error("Expected error for short ciphertext")
	}
}

func TestDecryptAES256GCM_WrongKey(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	rand.Read(key1)
	rand.Read(key2)

	plaintext := []byte("test data")
	ciphertext, _ := EncryptAES256GCM(plaintext, key1)

	_, err := DecryptAES256GCM(ciphertext, key2)
	if err == nil {
		t.Error("Expected error for wrong key")
	}
}
