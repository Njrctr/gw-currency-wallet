package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type ConfigDB struct {
	Host     string
	Port     string
	DBName   string
	Username string
	SSLMode  string
	Password string
}

type ClientConfig struct {
	Address string
}

type ConfigApp struct {
	Port     string
	TokenTTL int
	CacheTTL int
}

type Config struct {
	DB     ConfigDB
	App    ConfigApp
	Client ClientConfig
}

func NewConfig() (Config, error) {

	mode := flag.String("mode", "debug", "")
	flag.Parse()
	confFile := "config_dev.env"
	if *mode != "debug" && *mode != "release" {
		logrus.Fatalf("Неверный режим запуска: %s", *mode)
	}
	if *mode == "release" {
		confFile = "config_docker.env"
	}

	if err := godotenv.Load(confFile); err != nil {
		logrus.Fatalf("Ошибка получения переменных окружения: %s", err.Error())
	}
	tokenTTL, _ := strconv.Atoi(os.Getenv("APP_TOKEN_TTL"))
	cacheTTL, _ := strconv.Atoi(os.Getenv("APP_CACHE_TTL"))
	config := Config{
		App: ConfigApp{
			Port:     os.Getenv("APP_PORT"),
			TokenTTL: tokenTTL,
			CacheTTL: cacheTTL,
		},
		DB: ConfigDB{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			DBName:   os.Getenv("DB_DBNAME"),
			Username: os.Getenv("DB_USERNAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Client: ClientConfig{
			Address: os.Getenv("GRPC_HOST"),
		},
	}
	return config, nil
}
