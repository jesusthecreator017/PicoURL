package config

import (
	"fmt"
	"net/url"
	"time"

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
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	CacheTTL time.Duration // Cache Time-To-Live in seconds
}

func LoadRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
		Password: env.GetString("REDIS_PASSWORD", ""),
		DB:       env.GetInt("REDIS_DB", 0),
		CacheTTL: time.Duration(env.GetInt("REDIS_CACHE_TTL_MINUTES", 1440)) * time.Minute, // Default to 24 hours
	}
}

func LoadPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     env.GetString("POSTGRES_HOST", "localhost"),
		Port:     env.GetInt("POSTGRES_PORT", 5432),
		User:     env.GetString("POSTGRES_USER", "postgres"),
		Password: env.GetString("POSTGRES_PASSWORD", "password"),
		DB:       env.GetString("POSTGRES_DB", "pico_url"),
		SSLMode:  env.GetString("POSTGRES_SSLMODE", "disable"),
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

func DSN(cfg *PostgresConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, url.QueryEscape(cfg.Password), cfg.Host, cfg.Port, cfg.DB, cfg.SSLMode)
}
