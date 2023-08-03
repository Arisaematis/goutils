package rsa

import (
	"goutils/pkg/safety/base64"
	"io/ioutil"
	"os"
)

var (
	PemPublicData     []byte
	PemPrivateData    []byte
	BasePemPublicData []byte
	configFilePaths   = [2]string{
		"/etc/secret-volume/base-public.pem",
		"base-public.pem"}
)

func notExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	return err != nil
}
func InitPem() error {
	//public, _ := base64.StdEncoding.DecodeString(pubKey)
	for _, configPath := range configFilePaths {
		//从不同的层级目录初始化环境配置，直到有一次初始化成功后退出
		currentDir, _ := os.Getwd()
		klog.Infof("current directory: %s, load config: %s", currentDir, configPath)
		if notExists(configPath) {
			continue
		}
		basePublic, err := ioutil.ReadFile(configPath)
		if err != nil {
			return err
		}
		decodeString, err := base64.Decode(basePublic)
		if err != nil {
			return err
		}
		BasePemPublicData = decodeString
	}
	//打开文件
	public, err := ioutil.ReadFile("public.pem")
	if err != nil {
		return err
	}
	PemPublicData = public
	//打开文件
	private, err := ioutil.ReadFile("private.pem")
	if err != nil {
		return err
	}
	PemPrivateData = private
	return nil
}
