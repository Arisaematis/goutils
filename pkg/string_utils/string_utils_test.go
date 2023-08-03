package string_utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

const (
	sec = "SDKbK1cChRisq9dcNFVc9NtwTUzO0iRlBxTAW41G3yg"
)

func Test(t *testing.T) {
	b := StringToByteSlice(sec)
	fmt.Println(b)
	s := ByteSliceToString(b)
	fmt.Println(s)

	// s2 := Sha256(b)
	s3 := hex.EncodeToString(b)
	fmt.Println(s3)
}
func Test1(t *testing.T) {
	kSecret := StringToByteSlice(sec)
	fmt.Println(kSecret)
	kDate := signAlgorithm("20220903", kSecret)
	b, err := hmacsha256(kSecret, "20220903")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)

	fmt.Println(kDate)
	kRegion := signAlgorithm("", kDate)
	fmt.Println(kRegion)
	kService := signAlgorithm("", kRegion)
	fmt.Println(kService)
	fmt.Println(signAlgorithm("sdk_request", kService))
}

func Test2(t *testing.T) {
	secret1 := "2131231@#42"
	message := "我加密一下"
	b, err := hmacsha256([]byte(message), secret1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
}

//Sha256加密
func Sha256(src []byte) string {
	m := sha256.New()
	m.Write(src)
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

func signAlgorithm(stringData string, key []byte) []byte {
	data := StringToByteSlice(stringData)
	//创建对应的sha256哈希加密算法,这里的key为加密密钥
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write(data)
	secretMessage := hash.Sum(nil)
	return secretMessage
}

func hmacsha256(key []byte, data string) ([]byte, error) {
	h := hmac.New(sha256.New, []byte(key))
	if _, err := h.Write([]byte(data)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func TestReplaceAtPosition(t *testing.T) {
	s1 := "/v1/{projectId}/ecs/{id}"
	path := []string{"123", "456"}
	for _, v := range path {
		start := strings.Index(s1, "{")
		end := strings.Index(s1, "}") + 1
		if start != -1 && end != -1 && start < end {
			s1 = ReplaceAtPosition(s1, v, start, end)
			fmt.Printf("ReplaceAtPosition(s1, v, start, end): %v\n", s1)
		}
	}
}

func TestDecimal(t *testing.T) {
	f := 7.54887983706721
	fmt.Printf("Decimal(f): %v\n", Decimal(f))
}

func TestContainsString(t *testing.T) {
	arr := []string{"apple", "banana", "orange"}
	fmt.Printf("ContainsString(arr, \"ana\"): %v\n", ContainsString(arr, "ana"))
	fmt.Printf("ContainsString(arr, \"banana\"): %v\n", ContainsString(arr, "banana"))
}
