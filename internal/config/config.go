package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"go-rest-api/pkg/logging"
	"log"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type"`
		BindIp string `yaml:"bind_ip"`
		Port   string `yaml:"port" env:"PORT"`
	} `yaml:"listen"`
	MongoDB struct {
		Port       string `yaml:"port" env:"DB_PORT"`
		Host       string `yaml:"host" env:"DB_HOST"`
		DB         string `yaml:"database"`
		Username   string `yaml:"username" env:"DB_USERNAME"`
		Password   string `yaml:"password" env:"DB_PASSWORD"`
		Collection string `yaml:"collection"`
		AuthDB     string `yaml:"auth_db"`
	} `yaml:"mongodb"`
}

var instance *Config

var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
