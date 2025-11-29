package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DBConfig
	Auth AuthConfig
}

type AuthConfig struct {
	Secret           string
	AccessTokenMin   int64
	RefreshTokenDays int64
}

type DBConfig struct {
	DSN string
}

func LoadConfig() *Config {
	err := godotenv.Load()

	accessExpStr := os.Getenv("ACCESS_TOKEN_EXP_MIN")
	accessExpInt, err := strconv.Atoi(accessExpStr)

	if err != nil {
		accessExpInt = 15
	}

	refreshExpStr := os.Getenv("REFRESH_TOKEN_EXP_DAYS")
	refreshExtInt, err := strconv.Atoi(refreshExpStr)

	if err != nil {
		refreshExtInt = 7
	}
	if err != nil {
		log.Println("Error loading .env file")
	}
	cfg := &Config{
		DB: DBConfig{
			DSN: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret:           os.Getenv("JWT_SECRET"),
			AccessTokenMin:   int64(accessExpInt),
			RefreshTokenDays: int64(refreshExtInt),
		},
	}

	if cfg.DB.DSN == "" || cfg.Auth.Secret == "" {
		log.Fatal("Environment variables DSN or SECRET are not set")
	}

	return cfg
}
