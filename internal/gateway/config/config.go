package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type AutoGenerated struct {
	RPC RPC `yaml:"rpc"`
	Web Web `yaml:"web"`
	Log Log `yaml:"log"`
}
type RPC struct {
	User      string `yaml:"user"`
	VideoList string `yaml:"video_list"`
	Video     string `yaml:"video"`
	Social    string `yaml:"social"`
}
type Web struct {
	Addr string `yaml:"addr"`
}
type Log struct {
	Env        string `yaml:"env"`
	Path       string `yaml:"path"`
	Encoding   string `yaml:"encoding"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
}

//==============================================================================

// WrapperConfig
type WrapperConfig struct {
	AutoGenerated
}

var GlobalConfig *WrapperConfig

// 解析配置文件
func Init(path string) *WrapperConfig {
	// 读取配置文件解析 AutoGenerated 结构体
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var conf AutoGenerated
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		panic(err)
	}

	GlobalConfig = &WrapperConfig{
		AutoGenerated: conf,
	}

	return GlobalConfig
}
