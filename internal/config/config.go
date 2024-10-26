package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var GlobalConfig *Config

type Config struct {
	Run     runConfig     `yaml:"run"`
	Storage storageConfig `yaml:"storage"`
	Mysql   mysqlConfig   `yaml:"mysql"`
}
type runConfig struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
}

type storageConfig struct {
	Videos_data  string `yaml:"videos_data"`
	Videos_cover string `yaml:"videos_cover"`
	Images       string `yaml:"images"`
}

type mysqlConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Init() {
	s, _ := os.Getwd()
	configPath := filepath.Join(s, "config/application.yaml")

	dataBytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(err) // 使用panic来处理初始化失败的情况
	}
	config := Config{}
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		panic(err)
	}

	GlobalConfig = &config
	fmt.Println("已加载配置")
}

// GetConfig 返回全局配置的引用
func GetConfig() *Config {
	return GlobalConfig
}
