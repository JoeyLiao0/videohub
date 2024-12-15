package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var AppConfig *Config

type Config struct {
	Run     runConfig     `yaml:"run"`
	Storage storageConfig `yaml:"storage"`
	Static  staticConfig  `yaml:"static"`
	Mysql   mysqlConfig   `yaml:"mysql"`
	Redis   redisConfig   `yaml:"redis"`
	JWT     jwtConfig     `yaml:"jwt"`
	CORS    corsConfig    `yaml:"cors"`
	Email   emailConfig   `yaml:"email"`
	Log     logConfig     `yaml:"log"`
	Video   videoConfig   `yaml:"video"`
}

type runConfig struct {
	Name  string `yaml:"name"`
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

type storageConfig struct {
	Base        string `yaml:"base"`
	VideosData  string `yaml:"videos_data"`
	VideosCover string `yaml:"videos_cover"`
	VideosChunk string `yaml:"videos_chunk"`
	Images      string `yaml:"images"`
}

type staticConfig struct {
	Base   string `yaml:"base"`
	Video  string `yaml:"video"`
	Cover  string `yaml:"cover"`
	Avatar string `yaml:"avatar"`
}

type mysqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type redisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type jwtConfig struct {
	AccessTokenSecret  string `yaml:"access_token_secret"`
	AccessTokenExpire  uint   `yaml:"access_token_expire"`
	RefreshTokenSecret string `yaml:"refresh_token_secret"`
	RefreshTokenExpire uint   `yaml:"refresh_token_expire"`
}

type corsConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           uint     `yaml:"max_age"`
}

type emailConfig struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Expiration int    `yaml:"expiration"`
}

type logConfig struct {
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

type videoConfig struct {
	DefaultStatus int `yaml:"default_status"`
	DefaultPage   int `yaml:"default_page"`
	DefaultLimit  int `yaml:"default_limit"`
}

func InitConfig() {
	dataBytes, err := os.ReadFile("config/config.yaml")
	if err != nil {
		logrus.Fatalf("Error reading config file: %v", err)
	}

	config := Config{}
	if err := yaml.Unmarshal(dataBytes, &config); err != nil {
		logrus.Fatalf("Unable to decode into struct: %v", err)
	}

	AppConfig = &config
	logrus.Info("Config loaded successfully")
}
