package utils

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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

func GetConfig() (*Config, error) {

	s, _ := os.Getwd()

	configPath := filepath.Join(s, "config/application.yaml")

	dataBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil

}
