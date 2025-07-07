package config

import (
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load(".env")

type DbConfig struct {
	Host     string
	User     string
	Password string
	Name     string
}

func Db() DbConfig {
	return DbConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

var GatewaySecret = os.Getenv("GATEWAY_SECRET")
