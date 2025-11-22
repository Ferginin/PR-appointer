package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Env struct {
	DBName     string `env:"DB_NAME"`
	DBUsername string `env:"DB_USERNAME"`
	DBPassword string `env:"DB_PASSWORD"`
	DBPort     int    `env:"DB_PORT"`
	DBHost     string `env:"DB_HOST"`
	IPAddress  string `env:"IP_ADDRESS"`
	APIPort    int    `env:"API_PORT"`

	Environment string `env:"ENVIRONMENT"`
}

type Config struct {
	Env    Env
	Client *pgxpool.Pool
}

var config Config

func GetConfig() *Config {
	config.Env = *GetEnv()

	return &config
}

func GetEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file")
	}

	var cfg Env
	err = env.Parse(&cfg)
	if err != nil {
		slog.Error("Error parsing .env file:", err.Error(), nil)
		panic(err)
	}

	return &cfg
}
