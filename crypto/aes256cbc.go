package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func AES256CBCPKCS0Encrypt(src []byte, iv []byte, key []byte) ([]byte, error) {
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
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(dst, src)
	return dst, nil
}

func AES256CBCPKCS0Decrypt(src []byte, iv []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	srcLen := len(src)
	if srcLen%block.BlockSize() != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	dst := make([]byte, srcLen)

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(dst, src)
	dst = PKCSUnPadding(dst)
	return dst, nil
}
