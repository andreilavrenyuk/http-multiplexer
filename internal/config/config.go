package config

import (
	"os"
	"strconv"
	"time"
)

type appConfig struct {
	Port           string
	MaxUrls        int
	MaxRequests    int
	MaxOutRequests int
	RequestTimeout time.Duration
}

var config = &appConfig{
	Port:           getEnvString("PORT", "3000"),
	MaxUrls:        getEnvInt("MAX_URLS", 20),
	MaxRequests:    getEnvInt("MAX_REQUESTS", 1),
	MaxOutRequests: getEnvInt("MAXOUT_REQUESTS", 4),
	RequestTimeout: getEnvDuration("REQUEST_TIMEOUT", time.Second),
}

func Port() string {
	return config.Port
}

func MaxUrls() int {
	return config.MaxUrls
}

func MaxRequests() int {
	return config.MaxRequests
}

func MaxOutRequests() int {
	return config.MaxOutRequests
}

func RequestTimeout() time.Duration {
	return config.RequestTimeout
}

func getEnvString(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if result, err := strconv.Atoi(value); err == nil {
			return result
		}
	}

	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if result, err := time.ParseDuration(value); err == nil {
			return result
		}
	}

	return defaultValue
}
