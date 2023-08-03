package rsa

import "testing"

func TestGenerateRSAKey(t *testing.T) {
	//生成密钥对，保存到文件
	GenerateRSAKey(1024)
}
