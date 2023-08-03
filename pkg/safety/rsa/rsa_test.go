package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/lstack-org/utils/pkg/encipherment"
	"golang.org/x/crypto/scrypt"
)

func TestRSAEncrypt(t *testing.T) {
	encrypt := "cKMu47ZXBnpxZcwpVWORS7ETslBY8oa1ReR9NX3ORO4cX/0kVzjL9wFII0cU8uXgKY+Fg6jqNOc8qkHsZdgOPr/jQs2GbyQp521g8N7vFoUQGHtBHQIupEvGxBGvEubV4x+aJgm9n5hXOlLAObZsHXGSnmUumQlesjKOLG0BxaI="
	privateRsa := "MIICWgIBAAKBgQC0KVb7RfEzRVY+4x0kHwhPfCbIWKDQ8jUEIcO0E5noQdg4CUaeyLfJLNIQaAFENx2D0OO+VXertlo7X4tMXtXALBqYu8OAv2+0rNmKfuIejNf4iCYvi2ScBdqbom8mWViDQ/Zpe9cw22t1Y5NdzLHFJ0vf31nKHXRWZI69uLqYkQIDAQABAoGAIzb5Y4AWV1k0NHKcHZVbQH4Z7os0U+Mj7bzTzf0O1lEYfEuD3EGLeu0h2kcaCEVlpYBNI8T88Tlhhb11MuZOjT+P5scABh10Ia3IbjbvIiwtJSRhpZHvIvnBx4iyfgTsg8FO11f9ASR+ASPnWw+OCTdwyCxHXQHJSIPAWF7E9gkCQQDeZF3ETy7md/tLUzgP+pYEts+JUsPA1LOXzKwzvdAwNPPo+VBHGlZhfEJpgr+rdhhrwEoU/+9kiyHh/CgsPEprAkEAz2M1Wlxvonw8636l5pV+yw0HjQ8YbL3GHmvqgbejJvhM2dZvhaCnIGGjF0/osX+bMl9Y78+d3eqO29h8tLcf8wJAE21XF5gHM9DVXe4mHpc4Va8WkBtvyD+MdL1HabmyHxPxHq/wyFVPqHJvZsIqNjM5zOfeUNlOs0zIJ/KcG8kkgwI/LI4j6EXztfT7IZ0UB3YWx4kFFkkn9jTPW7nTqArMApNV73cifpMFVO+lGl0QoRHJRgk2Ek+ImyTJjHH2WNz/AkB6Kvz09v7WL/dKlFqLBXRevcuu0WFMLZ0+Qdsl6D2o/dTay+HryMV7upxIDedkgU4EXWjzz6Kl0AZuy9AIYuIH"
	des, err := base64.StdEncoding.DecodeString(privateRsa)
	if err != nil {
		panic(err)
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(des)
	if err != nil {
		panic(err)
	}
	//base64.StdEncoding.Decode(text, cipherText)
	//解密 对称加密的密钥
	ciphertext, err := base64.StdEncoding.DecodeString(encrypt)
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	fmt.Println("plainText:", string(plainText))
}

func Encrypt() (string, string) {
	//加密
	plainText := "Password@123"
	//对称加密
	aesEncrypt, _ := encipherment.AesEncrypt(plainText, "launcherSoNBPlus")
	// 内容base64
	accesskeySecret := base64.StdEncoding.EncodeToString([]byte(aesEncrypt))
	//base64.StdEncoding.Encode(accesskeySecret, []byte(aesEncrypt))
	publicRsa := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0KVb7RfEzRVY+4x0kHwhPfCbIWKDQ8jUEIcO0E5noQdg4CUaeyLfJLNIQaAFENx2D0OO+VXertlo7X4tMXtXALBqYu8OAv2+0rNmKfuIejNf4iCYvi2ScBdqbom8mWViDQ/Zpe9cw22t1Y5NdzLHFJ0vf31nKHXRWZI69uLqYkQIDAQAB"
	decodeString, _ := base64.StdEncoding.DecodeString(publicRsa)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(decodeString)
	if err != nil {
		panic(err)
	}
	//类型断言
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if ok {
		fmt.Println(publicKey)
	} else {
		panic("不是publicKey")
	}
	encrypt, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte("launcherSoNBPlus"))

	if err != nil {
		fmt.Println("File Encrypt error", err)
	}
	return base64.StdEncoding.EncodeToString(encrypt), accesskeySecret
}

func TestDecrypt(t *testing.T) {
	encrypt, accesskeySecret := Encrypt()
	privateRsa := "MIICWgIBAAKBgQC0KVb7RfEzRVY+4x0kHwhPfCbIWKDQ8jUEIcO0E5noQdg4CUaeyLfJLNIQaAFENx2D0OO+VXertlo7X4tMXtXALBqYu8OAv2+0rNmKfuIejNf4iCYvi2ScBdqbom8mWViDQ/Zpe9cw22t1Y5NdzLHFJ0vf31nKHXRWZI69uLqYkQIDAQABAoGAIzb5Y4AWV1k0NHKcHZVbQH4Z7os0U+Mj7bzTzf0O1lEYfEuD3EGLeu0h2kcaCEVlpYBNI8T88Tlhhb11MuZOjT+P5scABh10Ia3IbjbvIiwtJSRhpZHvIvnBx4iyfgTsg8FO11f9ASR+ASPnWw+OCTdwyCxHXQHJSIPAWF7E9gkCQQDeZF3ETy7md/tLUzgP+pYEts+JUsPA1LOXzKwzvdAwNPPo+VBHGlZhfEJpgr+rdhhrwEoU/+9kiyHh/CgsPEprAkEAz2M1Wlxvonw8636l5pV+yw0HjQ8YbL3GHmvqgbejJvhM2dZvhaCnIGGjF0/osX+bMl9Y78+d3eqO29h8tLcf8wJAE21XF5gHM9DVXe4mHpc4Va8WkBtvyD+MdL1HabmyHxPxHq/wyFVPqHJvZsIqNjM5zOfeUNlOs0zIJ/KcG8kkgwI/LI4j6EXztfT7IZ0UB3YWx4kFFkkn9jTPW7nTqArMApNV73cifpMFVO+lGl0QoRHJRgk2Ek+ImyTJjHH2WNz/AkB6Kvz09v7WL/dKlFqLBXRevcuu0WFMLZ0+Qdsl6D2o/dTay+HryMV7upxIDedkgU4EXWjzz6Kl0AZuy9AIYuIH"
	des, err := base64.StdEncoding.DecodeString(privateRsa)
	if err != nil {
		panic(err)
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(des)
	if err != nil {
		panic(err)
	}
	//base64.StdEncoding.Decode(text, cipherText)
	//解密 对称加密的密钥
	ciphertext, err := base64.StdEncoding.DecodeString(encrypt)
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	fmt.Println("plainText:", string(plainText))
	decodeString, err := base64.StdEncoding.DecodeString(accesskeySecret)
	// 用密钥解密
	text := encipherment.AesDecrypt(string(decodeString), string(plainText))
	fmt.Println(text)
}

func TestRSADecrypt(t *testing.T) {
	//加密
	data := []byte("Password@123")
	encrypt, err := RSAEncrypt(data)
	if err != nil {
		fmt.Println("File Encrypt error", err)
	}
	fmt.Printf("需要加密的内容：%v\n", encrypt)
	toString := base64.StdEncoding.EncodeToString(encrypt)
	fmt.Printf("加密后内容：%s\n", toString)
	// 解密
	decodeString, err := base64.StdEncoding.DecodeString(toString)
	if err != nil {
		fmt.Println("base64 Decode error", err)
	}
	fmt.Printf("需要解密的内容：%v\n", decodeString)
	decrypt, err := RSADecrypt(decodeString)
	if err != nil {
		fmt.Println("File Decode error", err)
	}
	fmt.Println(string(decrypt[:]))
}

func TestRead(t *testing.T) {
	//打开文件
	bufPrv, err := ioutil.ReadFile("private.pem")
	if err != nil {
		fmt.Println("读取文件失败", err)
	}
	fmt.Printf("私钥：%s\n", base64.StdEncoding.EncodeToString(bufPrv))
	//打开文件
	bufPub, err := ioutil.ReadFile("public.pem")
	if err != nil {
		fmt.Println("读取文件失败", err)
	}
	fmt.Printf("公钥：%s\n", base64.StdEncoding.EncodeToString(bufPub))
}

func TestSignVer(t *testing.T) {
	s := "Password@123"
	thas := sha1.New()
	salt := "12345"
	key, _ := scrypt.Key([]byte(s), []byte(salt), 16384, 8, 1, 32)
	fmt.Printf("加盐后: %x\n", key)
	thas.Write(key)
	sum := thas.Sum(nil)
	fmt.Printf("哈希值：%x\n", sum)
	fmt.Printf("哈希值：%x\n", string(sum))
	sign, err := RsaSignWithSha1Hex(string(sum), "basePrivate.pem")
	if err != nil {
		fmt.Println("签名失败", err)
		return
	}
	fmt.Println("签名后的信息：", sign)
	err = RsaVerySignWithSha1Base64(string(sum), sign)
	if err != nil {
		fmt.Println("签名校验失败", err)
		return
	}
	fmt.Println("签名校验成功")
}
