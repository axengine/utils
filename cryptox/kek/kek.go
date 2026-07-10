package kekcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// Argon2Params 保存密钥派生参数，需与密文一起持久化。
type Argon2Params struct {
	Memory      uint32 // KiB
	Iterations  uint32
	Parallelism uint8
	SaltLength  int // 字节数，建议 16
}

// EncryptedKeyData 是加密后的私钥信封，可直接序列化（例如 JSON）存入数据库。
type EncryptedKeyData struct {
	Salt         []byte       `json:"salt"`
	Argon2Params Argon2Params `json:"argon2_params"`
	// Ciphertext 格式：IV(12字节) || ciphertext || tag(16字节)
	Ciphertext []byte `json:"ciphertext"`
}

// 默认安全参数 (OWASP 推荐 2023)
var DefaultParams = Argon2Params{
	Memory:      256 * 1024, // 256 MiB
	Iterations:  4,
	Parallelism: 2,
	SaltLength:  16,
}

// DeriveKEK 从密码和盐派生 AES-256 密钥（KEK）。
// 注意：前端调用此函数派生 KEK，后端不应执行此步骤（保持零密码）。
func DeriveKEK(password []byte, salt []byte, params Argon2Params) []byte {
	return argon2.IDKey(
		password,
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		32, // 256 bits
	)
}

// EncryptPrivateKey 使用用户密码加密私钥，返回可持久化的信封结构。
// privKey: 私钥明文
// password: 用户输入的密码（应保证高强度）
func EncryptPrivateKey(privKey []byte, password []byte) (*EncryptedKeyData, error) {
	if len(password) == 0 {
		return nil, errors.New("password cannot be empty")
	}
	if len(privKey) == 0 {
		return nil, errors.New("private key cannot be empty")
	}

	params := DefaultParams

	// 1. 生成随机盐
	salt := make([]byte, params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}

	// 2. 派生 KEK
	kek := DeriveKEK(password, salt, params)
	defer ClearBytes(kek) // 使用后立即擦除

	// 3. AES-256-GCM 加密
	block, err := aes.NewCipher(kek)
	if err != nil {
		return nil, fmt.Errorf("aes cipher: %w", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("gcm: %w", err)
	}

	// 4. 生成随机 IV (Nonce)
	iv := make([]byte, aesGCM.NonceSize()) // 12 字节
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("generate iv: %w", err)
	}

	// 5. 加密，可选的附加认证数据 (AAD)
	aad := []byte("private-key-v1") // 绑定用途
	ciphertext := aesGCM.Seal(nil, iv, privKey, aad)

	// 6. 组装密文: IV || ciphertext || tag
	// GCM 的 Seal 已将 tag 附加在 ciphertext 末尾
	combined := make([]byte, 0, len(iv)+len(ciphertext))
	combined = append(combined, iv...)
	combined = append(combined, ciphertext...)

	return &EncryptedKeyData{
		Salt:         salt,
		Argon2Params: params,
		Ciphertext:   combined,
	}, nil
}

// DecryptPrivateKey 使用密钥加密密钥 (KEK) 解密私钥信封。
// KEK 应由前端从密码派生后，通过安全通道（如 TLS + 额外 payload 加密）传来。
// 解密后调用方应负责及时擦除 returned plaintext。
func DecryptPrivateKey(data *EncryptedKeyData, kek []byte) ([]byte, error) {
	if data == nil {
		return nil, errors.New("missing encrypted key data")
	}
	if len(kek) != 32 {
		return nil, errors.New("KEK must be 32 bytes")
	}
	if len(data.Ciphertext) < 12+16 { // 至少 IV + tag
		return nil, errors.New("ciphertext too short")
	}

	// 1. 创建 AES-GCM
	block, err := aes.NewCipher(kek)
	if err != nil {
		return nil, fmt.Errorf("aes cipher: %w", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("gcm: %w", err)
	}

	// 2. 提取 IV
	nonceSize := aesGCM.NonceSize() // 12
	iv := data.Ciphertext[:nonceSize]
	ciphertext := data.Ciphertext[nonceSize:]

	// 3. 解密（认证 tag 已包含在 ciphertext 中）
	aad := []byte("private-key-v1")
	plaintext, err := aesGCM.Open(nil, iv, ciphertext, aad)
	if err != nil {
		return nil, fmt.Errorf("decryption failed (wrong key or tampered data): %w", err)
	}

	return plaintext, nil
}

// ClearBytes 将字节切片全置零，用于安全擦除敏感数据。
func ClearBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
