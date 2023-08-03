package setting

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	_ "github.com/lstack-org/go-web-framework/pkg/version"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

func init() {
	setup("")
}

var (
	configFilePaths = [5]string{
		"conf/config.yaml",
		"../conf/config.yaml",
		"../../conf/config.yaml",
		"../../../conf/config.yaml",
		"../../../../conf/config.yaml"}
	//EnvConfig 环境配置，指定是什么环境，如开发环境，测试环境，本地等
	EnvConfig *envConfig
	//HTTPSetting http相关配置
	HTTPSetting *httpSetting
	//RemoteSetting 远程调用配置
	RemoteSetting *remoteSetting
	//AppConfig 应用配置
	AppConfig *appConfig
	//NotifySetting 邮件通知配置
	//NotifySetting *notifySetting
	// redis 配置信息
	RedisSetting *redisSetting

	maxLogLevel int32 = 7
)

type envConfig struct {
	Title      string                  `yaml:"title"`
	ReleaseEnv string                  `yaml:"releaseEnv"`
	Version    string                  `yaml:"version"`
	ClientType string                  `yaml:"clientType"`
	Server     map[string]ServerConfig `yaml:"server"`
}

//ServerConfig 服务配置
type ServerConfig struct {
	AppConfig     appConfig     `yaml:"appConfig"`
	HTTPSetting   httpSetting   `yaml:"httpSetting"`
	RemoteSetting remoteSetting `yaml:"remoteSetting"`
	//NotifySetting notifySetting `yaml:"notifySetting"`
	RedisSetting redisSetting `yaml:"redisSetting"`
}

type appConfig struct {
	WorkId   int64 `yaml:"workId"`
	LogLevel int32 `yaml:"logLevel"`
	StarCron int64 `yaml:"starCron"`
}

type remoteSetting struct {
	MongoAddr     string `yaml:"mongoAddr"`
	MongoDatabase string `yaml:"mongoDatabase"`
}

//type notifySetting struct {
//	MailHost string `yaml:"mailHost"`
//	MailPort int    `yaml:"mailPort"`
//	MailUser string `yaml:"mailUser"` // 发件人
//	MailPass string `yaml:"mailPass"` // 发件人密码
//	MailTo   string `yaml:"mailTo"`   // 收件人 多个用,分割
//	DingHook string `yaml:"dingHook"`
//}

type httpSetting struct {
	HTTPPort     string        `yaml:"httpPort"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	RunMode      string        `yaml:"runMode"`
}
type redisSetting struct {
	RedisHost string `yaml:"redisHost"`
	RedisPwd  string `yaml:"redisPwd"`
	RedisDB   int    `yaml:"redisDb"`
}

func setup(path string) {
	loadConfig(path)
}

func loadConfig(path string) {
	if path != "" {
		// 根据传参初始化环境配置
		if notExists(path) {
			panic(fmt.Errorf("config file path: %s not exists", path))
		}

		loadConfigWithPath(path)
		return
	}
	for _, configPath := range configFilePaths {
		//从不同的层级目录初始化环境配置，直到有一次初始化成功后退出
		currentDir, _ := os.Getwd()
		klog.Infof("current directory: %s, load config: %s", currentDir, configPath)
		if notExists(configPath) {
			continue
		}
		loadConfigWithPath(configPath)
		break
	}
	if EnvConfig == nil {
		panic("envConfig init fail")
	}
}

func loadConfigWithPath(configPath string) {
	viper.SetConfigFile(configPath) // 指定配置文件路径
	err := viper.ReadInConfig()     // 读取配置信息
	if err != nil {                 // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(&EnvConfig); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}

	// 监控配置文件变化
	viper.WatchConfig()
	// 注意！！！配置文件发生变化后要同步到全局变量Conf
	viper.OnConfigChange(func(in fsnotify.Event) {
		klog.Infof("config file changed:", in.Name)
		if err := viper.Unmarshal(&EnvConfig); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
		}
	})

	klog.Infof("read Config: %s", configPath)

	releaseEnv := os.Getenv("releaseEnv")
	if releaseEnv == "" {
		releaseEnv = EnvConfig.ReleaseEnv
	}
	klog.Info("releaseEnv: ", releaseEnv)
	serverConfig := EnvConfig.Server[releaseEnv]
	HTTPSetting = &serverConfig.HTTPSetting
	HTTPSetting.ReadTimeout = HTTPSetting.ReadTimeout * time.Second
	HTTPSetting.WriteTimeout = HTTPSetting.WriteTimeout * time.Second
	RemoteSetting = &serverConfig.RemoteSetting
	AppConfig = &serverConfig.AppConfig
	RedisSetting = &serverConfig.RedisSetting
	if AppConfig.LogLevel > maxLogLevel {
		klog.Warningf("log level：%v > maxLogLevel(%v), max value %v used.", AppConfig.LogLevel, maxLogLevel, maxLogLevel)
		AppConfig.LogLevel = maxLogLevel
	}
	//NotifySetting = &serverConfig.NotifySetting
	klog.Infof("config: %+v", prettifyJSON(EnvConfig.Server[releaseEnv]))
}

func prettifyJSON(i interface{}) string {
	var str []byte
	var err error
	str, err = json.MarshalIndent(i, "", "    ")
	if err != nil {
		klog.Fatal(err)
	}

	return string(str)
}

func notExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	return err != nil
}
