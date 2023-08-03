package aes

import (
	"bytes"
	aesCry "crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Aes struct {
	key   []byte
	iv    []byte
	block cipher.Block
}

func CreateCipher(key string) Aes {
	var aes Aes
	aes.key = []byte(key)
	aes.iv = []byte("launcher20170601")
	var err error
	aes.block, err = aesCry.NewCipher(aes.key)
	if err != nil {
		panic(err)
	}
	return aes
}

// AesBase64Encrypt 加密
func (a *Aes) AesBase64Encrypt(in string) (string, error) {
	origData := []byte(in)
	origData = PKCS5Padding(origData, a.block.BlockSize())
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	bm := cipher.NewCBCEncrypter(a.block, a.iv)
	bm.CryptBlocks(crypted, origData)
	var b = base64.StdEncoding.EncodeToString(crypted)
	return b, nil
}

// AesBase64Decrypt 解密
func (a *Aes) AesBase64Decrypt(b string) (string, error) {
	crypted, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		return "", err
	}
	origData := make([]byte, len(crypted))
	bm := cipher.NewCBCDecrypter(a.block, a.iv)
	bm.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	var out = string(origData)
	return out, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
