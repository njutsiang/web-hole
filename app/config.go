package app

import (
	"github.com/njutsiang/web-hole/exception"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

// 配置
var Config = ConfigYaml{}

// 配置数据结构
type ConfigYaml struct {
	Log struct {
		Level         string `yaml:"Level"`
		ExportConsole int    `yaml:"ExportConsole"`
		ExportFile    struct {
			Path string `yaml:"Path"`
		} `yaml:"ExportFile"`
	} `yaml:"Log"`
	Frontend struct {
		HttpPort      int    `yaml:"HttpPort"`
		HttpTimeout   int    `yaml:"HttpTimeout"`
		WebsocketPort int    `yaml:"WebsocketPort"`
		WebsocketPath string `yaml:"WebsocketPath"`
		SecretKey     string `yaml:"SecretKey"`
	} `yaml:"Frontend"`
	Proxy struct {
		FrontendUrl string `yaml:"FrontendUrl"`
		BackendHost string `yaml:"BackendHost"`
		SecretKey   string `yaml:"SecretKey"`
	} `yaml:"Proxy"`
	Backend struct {
		HttpPort int `yaml:"HttpPort"`
	} `yaml:"Backend"`
}

// 读取配置文件
func InitConfig() {
	file, err := os.Open("./config.yaml")
	if err != nil {
		exception.Throw(exception.InitConfigFailed, "打开配置文件失败 " + err.Error())
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		exception.Throw(exception.InitConfigFailed, "读取配置文件失败 " + err.Error())
		return
	}
	config := ConfigYaml{}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		exception.Throw(exception.InitConfigFailed, "解析配置文件失败 " + err.Error())
		return
	}
	Config = config
}