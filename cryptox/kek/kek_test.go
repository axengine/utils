package kekcrypto

import (
	"encoding/json"
	"testing"
)

func TestEncryptAndDecryptPrivateKey(t *testing.T) {
	// --- 添加私钥时（前端执行） ---
	password := []byte("user-strong-password")
	privKey := []byte("-----BEGIN PRIVATE KEY-----...")

	envelope, err := EncryptPrivateKey(privKey, password)
	if err != nil {
		t.Fatalf("EncryptPrivateKey failed: %v", err)
	}
	b, _ := json.MarshalIndent(envelope, "", "  ")
	t.Logf("envelope: %s", b)
	// 将 envelope 序列化后发送到后端存储

	// --- 签名时 ---
	// 前端重新派生 KEK（需从后端获取盐和参数）
	kek := DeriveKEK(password, envelope.Salt, envelope.Argon2Params)
	// 将 kek 通过安全通道发送给后端（不传密码）
	// sendKekToBackend(kek, envelope)

	// --- 后端解密私钥并签名 ---
	privKeyDecrypted, err := DecryptPrivateKey(envelope, kek)
	if err != nil {
		t.Fatalf("DecryptPrivateKey failed: %v", err)
	}
	t.Logf("decrypted: %s", privKeyDecrypted)
	defer ClearBytes(privKeyDecrypted) // 使用后立即擦除
	// 执行签名操作...
}
