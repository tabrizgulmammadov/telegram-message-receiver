package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	TelegramToken         string
	StoragePath           string
	BaseFileURL           string
	Debug                 bool
	SendAcknowledgment    bool
	AcknowledgmentMessage string
	MaxFileSize           int64
	AllowedFileTypes      []string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		TelegramToken:         os.Getenv("TELEGRAM_BOT_TOKEN"),
		StoragePath:           getEnvWithDefault("STORAGE_PATH", "messages"),
		BaseFileURL:           getEnvWithDefault("TELEGRAM_FILE_BASE_URL", "https://api.telegram.org/file/bot%s/%s"),
		Debug:                 os.Getenv("DEBUG") == "true",
		SendAcknowledgment:    getEnvWithDefault("SEND_ACKNOWLEDGMENT", "true") == "true",
		AcknowledgmentMessage: getEnvWithDefault("ACKNOWLEDGMENT_MESSAGE", "Message received!"),
		MaxFileSize:           getEnvAsInt64("MAX_FILE_SIZE", 20*1024*1024), // 20MB default
	}

	if config.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	return config, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
