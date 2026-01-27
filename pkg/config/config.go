package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	RedisAddr  string
	GRPCPort   string
	TickRate   time.Duration
	MaxMMRDiff int
}

func Load() Config {
	return Config{
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
		GRPCPort:   getEnv("GRPC_PORT", ":50051"),
		TickRate:   getEnvDuration("MATCH_INTERVAL", 200*time.Millisecond),
		MaxMMRDiff: getEnvInt("MAX_MMR_DIFF", 100),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}
