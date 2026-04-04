package config

import (
	"github.com/jesusthecreator017/PicoURL/internal/env"
)

type Config struct {
	Port       string
	CorsOrigin string
	Postgres   PostgresConfig
	Redis      RedisConfig
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func LoadRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
		Password: env.GetString("REDIS_PASSWORD", ""),
		DB:       env.GetInt("REDIS_DB", 0),
	}
}

func LoadPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     env.GetString("POSTGRES_HOST", "localhost"),
		Port:     env.GetInt("POSTGRES_PORT", 5432),
		User:     env.GetString("POSTGRES_USER", "postgres"),
		Password: env.GetString("POSTGRES_PASSWORD", "password"),
		DB:       env.GetString("POSTGRES_DB", "pico_url"),
	}
}

func LoadConfig() *Config {
	return &Config{
		Port:       env.GetString("PORT", "8080"),
		CorsOrigin: env.GetString("CORS_ORIGIN", "*"),
		Postgres:   *LoadPostgresConfig(),
		Redis:      *LoadRedisConfig(),
	}
}
