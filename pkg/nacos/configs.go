package nacos

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
)

var Config = new(MyConfig)

type App struct {
	Port     int        `yaml:"port"`
	Name     string     `yaml:"name"`
	Services ServiceMap `yaml:"services"`
}
type ServiceMap map[string]string

type Log struct {
	ErrorPath string `yaml:"error_path" mapstructure:"error_path"`
	InfoPath  string `yaml:"info_path" mapstructure:"info_path"`
	MaxAge    int    `yaml:"max_age" mapstructure:"max_age"`
	Rotation  int    `yaml:"rotation" mapstructure:"rotation"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
}

type Jwt struct {
	AccessTokenExpiredTime  int64  `json:"access_token_timeout" mapstructure:"access_token_expired_time"`
	RefreshTokenExpiredTime int64  `json:"refresh_token_timeout" mapstructure:"refresh_token_expired_time"`
	Secret                  string `json:"secret"`
}

type MyConfig struct {
	*App
	*Log
	*Jwt
}

func InitConfig() {

	//// 加载配置
	//viper.SetConfigFile("./configs/configs.yaml")
	//
	//// 监听配置
	//viper.WatchConfig()
	//
	//// 监听是否更改配置文件
	//viper.OnConfigChange(func(e fsnotify.Event) {
	//	if err := viper.Unmarshal(&Config); err != nil {
	//		panic(err)
	//	}
	//})
	//
	//if err := viper.ReadInConfig(); err != nil {
	//	panic(fmt.Errorf("ReadInConfig failed, err: %v", err))
	//}
	//if err := viper.Unmarshal(&Config); err != nil {
	//	panic(fmt.Errorf("unmarshal failed, err: %v", err))
	//}

	// 初始化Nacos配置

	// 获取配置信息
	content, err := NacosClient.GetConfig()
	if err != nil {
		panic(fmt.Errorf("GetConfig failed, err: %v", err))
	}

	viper.SetConfigType("yaml")
	if err = viper.ReadConfig(bytes.NewBuffer([]byte(content))); err != nil {
		panic(fmt.Errorf("ReadConfig failed, err: %v", err))
	}

	if err = viper.Unmarshal(&Config); err != nil {
		panic(fmt.Errorf("unmarshal failed, err: %v", err))
	}
}
