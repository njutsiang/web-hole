package app

import (
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

var Config = ConfigYaml{}

type ConfigYaml struct {
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
		log.Println(err)
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}
	config := ConfigYaml{}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Println(err)
		return
	}
	Config = config
}