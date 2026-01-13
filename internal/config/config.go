package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort         string
	DatabaseURL        string
	JWTSecret          string
	JWTRefreshSecret   string
	JWTExpiresIn       time.Duration
	JWTRefreshExpires  time.Duration
	RedisURL           string
	GoogleClientID     string
	GoogleClientSecret string
	AWSS3Bucket        string
	SMTPHost           string
	SMTPPort           int
	SMTPUser           string
	SMTPPassword       string
	Environment        string
}

func Load() *Config {
	return &Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/crm?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "your-256-bit-secret"),
		JWTRefreshSecret:   getEnv("JWT_REFRESH_SECRET", "your-refresh-256-bit-secret"),
		JWTExpiresIn:       getDurationEnv("JWT_EXPIRES_IN", 15*time.Minute),
		JWTRefreshExpires:  getDurationEnv("JWT_REFRESH_EXPIRES", 7*24*time.Hour),
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		AWSS3Bucket:        getEnv("AWS_S3_BUCKET", ""),
		SMTPHost:           getEnv("SMTP_HOST", "localhost"),
		SMTPPort:           getIntEnv("SMTP_PORT", 587),
		SMTPUser:           getEnv("SMTP_USER", ""),
		SMTPPassword:       getEnv("SMTP_PASSWORD", ""),
		Environment:        getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
