package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

const (
	PublicRsa  = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0KVb7RfEzRVY+4x0kHwhPfCbIWKDQ8jUEIcO0E5noQdg4CUaeyLfJLNIQaAFENx2D0OO+VXertlo7X4tMXtXALBqYu8OAv2+0rNmKfuIejNf4iCYvi2ScBdqbom8mWViDQ/Zpe9cw22t1Y5NdzLHFJ0vf31nKHXRWZI69uLqYkQIDAQAB"
	PrivateRsa = "MIICWgIBAAKBgQC0KVb7RfEzRVY+4x0kHwhPfCbIWKDQ8jUEIcO0E5noQdg4CUaeyLfJLNIQaAFENx2D0OO+VXertlo7X4tMXtXALBqYu8OAv2+0rNmKfuIejNf4iCYvi2ScBdqbom8mWViDQ/Zpe9cw22t1Y5NdzLHFJ0vf31nKHXRWZI69uLqYkQIDAQABAoGAIzb5Y4AWV1k0NHKcHZVbQH4Z7os0U+Mj7bzTzf0O1lEYfEuD3EGLeu0h2kcaCEVlpYBNI8T88Tlhhb11MuZOjT+P5scABh10Ia3IbjbvIiwtJSRhpZHvIvnBx4iyfgTsg8FO11f9ASR+ASPnWw+OCTdwyCxHXQHJSIPAWF7E9gkCQQDeZF3ETy7md/tLUzgP+pYEts+JUsPA1LOXzKwzvdAwNPPo+VBHGlZhfEJpgr+rdhhrwEoU/+9kiyHh/CgsPEprAkEAz2M1Wlxvonw8636l5pV+yw0HjQ8YbL3GHmvqgbejJvhM2dZvhaCnIGGjF0/osX+bMl9Y78+d3eqO29h8tLcf8wJAE21XF5gHM9DVXe4mHpc4Va8WkBtvyD+MdL1HabmyHxPxHq/wyFVPqHJvZsIqNjM5zOfeUNlOs0zIJ/KcG8kkgwI/LI4j6EXztfT7IZ0UB3YWx4kFFkkn9jTPW7nTqArMApNV73cifpMFVO+lGl0QoRHJRgk2Ek+ImyTJjHH2WNz/AkB6Kvz09v7WL/dKlFqLBXRevcuu0WFMLZ0+Qdsl6D2o/dTay+HryMV7upxIDedkgU4EXWjzz6Kl0AZuy9AIYuIH"
)

// RSAEncrypt RSA加密
// plainText 要加密的数据
// path 公钥匙文件地址
func RSAEncrypt(plainText []byte) ([]byte, error) {
	var text []byte
	block, err := base64.StdEncoding.DecodeString(PublicRsa)
	if err != nil {
		return nil, err
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block)
	if err != nil {
		return nil, err
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	base64.StdEncoding.Encode(text, cipherText)
	//返回密文
	return text, err
}

// RSADecrypt RSA解密
// cipherText 需要解密的byte数据
// path 私钥文件路径
func RSADecrypt(cipherText []byte) ([]byte, error) {
	block, err := base64.StdEncoding.DecodeString(PrivateRsa)
	if err != nil {
		return nil, err
	}
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block)
	if err != nil {
		return nil, err
	}
	text, err := base64.StdEncoding.DecodeString(string(cipherText))
	if err != nil {
		return nil, err
	}
	//对密文进行解密
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(text))
	//返回明文
	return plainText, err
}

// RsaSignWithSha1Hex 签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
func RsaSignWithSha1Hex(data string, prvKey string) (string, error) {
	//keyByts, err := hex.DecodeString(prvKey)
	//if err != nil {
	//	fmt.Println(err)
	//	return "", err
	//}
	//打开文件
	keyBytes, err := ioutil.ReadFile(prvKey)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(keyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		return "", err
	}
	h := sha1.New()
	h.Write([]byte(data))
	hash := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hash[:])
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return "", err
	}
	out := base64.StdEncoding.EncodeToString(signature)
	return out, nil
}

// RsaVerySignWithSha1Base64 验签：对采用sha1算法进行签名后转base64格式的数据进行验签
func RsaVerySignWithSha1Base64(originalData, signData string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(BasePemPublicData)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	hash := sha1.New()
	hash.Write([]byte(originalData))
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
}
