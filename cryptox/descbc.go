package cryptox

import (
	"crypto/cipher"
	"crypto/des"
)

func DESCBCPCSK5Encrypt(key []byte, plaintext []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext = PKCSPadding(plaintext, block.BlockSize())
	// 创建CBC模式的加密器
	iv := key // 使用密钥作为初始向量
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, err
}

func DESCBCPCSK5Decrypt(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 创建CBC模式的加密器
	iv := key // 使用密钥作为初始向量
	mode := cipher.NewCBCDecrypter(block, iv)
	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	// 去除填充数据
	plaintext = PKCSUnPadding(plaintext)
	return plaintext, err
}
