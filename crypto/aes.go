package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// AESDecrypt AES加密 初始向量16字节空 PKCS5 CBC
// 入参:src 待加密[]byte
// key:密钥[]byte 16/24/32
// 返回:加密后[]byte
func AESEncrypt(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = PKCSPadding(src, block.BlockSize())
	srcLen := len(src)
	if srcLen%block.BlockSize() != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	dst := make([]byte, srcLen)
	iv := make([]byte, 16)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(dst, src)
	return dst, nil
}

// AESEncrypt AES解密 初始向量16字节空 PKCS5 CBC
// 入参:src 已加密[]byte
// key:密钥[]byte 16/24/32
// 返回:解密后[]byte
func AESDecrypt(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	srcLen := len(src)
	if srcLen%block.BlockSize() != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	dst := make([]byte, srcLen)
	iv := make([]byte, 16)
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, src)
	return PKCSPadding(dst, block.BlockSize()), nil
}
