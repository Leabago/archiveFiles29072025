package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const cfgFileName = "config/config.yaml"

type Config struct {
	App     App     `yaml:"app"`
	HTTP    HTTP    `yaml:"http"`
	Log     Log     `yaml:"log"`
	Swagger Swagger `yaml:"swagger"`
}

type App struct {
	Name         string `yaml:"APP_NAME" env-default:"my-app"`
	Version      string `yaml:"APP_VERSION" env-default:"1"`
	MaxNumLinks  int    `yaml:"APP_NUM_LINKS" env-default:"3"`
	ContentType  string `yaml:"APP_CONTENT_TYPE"`
	MaxBytesResp int64  `yaml:"APP_MAX_BYTES_RESP" env-default:"51200"`
}

type HTTP struct {
	Port string `yaml:"HTTP_PORT" env-default:":8080"`
}

type Log struct {
	Level string `yaml:"LOG_LEVEL" env-default:"debug"`
}

type Swagger struct {
	Enable bool `yaml:"SWAGGER_ENABLED" env-default:"false"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if _, err := os.Stat(cfgFileName); err != nil {
		return nil, fmt.Errorf("error opening config file: %s", err)
	}

	// Читаем конфиг-файл и заполняем нашу структуру
	err := cleanenv.ReadConfig(cfgFileName, cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}

	fmt.Println("cfg:", *cfg)

	return cfg, nil
}
