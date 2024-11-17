package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
	// "github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	Run     runConfig     `yaml:"run"`
	Storage storageConfig `yaml:"storage"`
	Mysql   mysqlConfig   `yaml:"mysql"`
	JWT     jwtConfig     `yaml:"jwt"`
}
type runConfig struct {
	Name string `yaml:"name"`
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
}

type storageConfig struct {
	VideosData  string `yaml:"videos_data"`
	VideosCover string `yaml:"videos_cover"`
	VideosChunk string `yaml:"videos_chunk"`
	Images      string `yaml:"images"`
}

type mysqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type jwtConfig struct {
	AccessTokenSecret  string `yaml:"access_token_secret"`
	AccessTokenExpire  uint   `yaml:"access_token_expire"`
	RefreshTokenSecret string `yaml:"refresh_token_secret"`
	RefreshTokenExpire uint   `yaml:"refresh_token_expire"`
}

func InitConfig() {
	// viper 无法识别环境变量中的下划线
	// viper.SetConfigName("config")
	// viper.SetConfigType("yaml")
	// viper.AddConfigPath("./config")

	// if err := viper.ReadInConfig(); err != nil {
	// 	log.Fatalf("Error reading config file: %v", err)
	// }

	// AppConfig = &Config{}

	// if err := viper.Unmarshal(AppConfig); err != nil {
	// 	log.Fatalf("Unable to decode into struct: %v", err)
	// }

	dataBytes, err := os.ReadFile("./config/config.yaml")

	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	config := Config{}
	if err := yaml.Unmarshal(dataBytes, &config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
	
	AppConfig = &config
}
