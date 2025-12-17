package cryptox

import (
	"crypto/des"
	"errors"
)

// TripleDesECBEncrypt 3des加密 使用加密模式为ECB 填充方式为pkcs5
// 入参:待加密原始数据 字符串key
// 返回值:加密后的base64字符串
func TripleDesECBEncrypt(origData []byte, key string) ([]byte, error) {
	bKey := transKey(key)
	block, err := des.NewTripleDESCipher(bKey)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	origData = PKCSPadding(origData, bs)
	if len(origData)%bs != 0 {
		return nil, errors.New("need a multiple of the blocksize")
	}
	out := make([]byte, len(origData))
	dst := out
	for len(origData) > 0 {
		block.Encrypt(dst, origData[:bs])
		origData = origData[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

// TripleDesECBDecrypt 3des解密
func TripleDesECBDecrypt(crypted []byte, key string) ([]byte, error) {
	bKey := transKey(key)
	block, err := des.NewTripleDESCipher(bKey)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(crypted))
	dst := out
	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}
	out = PKCSPadding(out, block.BlockSize())
	return out, nil
}

// 0填充
func transKey(key string) []byte {
	//key只取24位
	bKey := make([]byte, 24)
	keys := []byte(key)
	length := len(keys)
	if length == 24 {
		copy(bKey, keys)
	} else if length < 24 {
		copy(bKey, keys)
		for i := length; i < length; i++ {
			bKey[i] = 0
		}
	} else {
		copy(bKey, keys[:24])
	}
	return bKey
}
