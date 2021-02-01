package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken   string
	MainChatID int64
	TechChatID int64
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		BotToken:   getEnvAsStr("BOT_TOKEN"),
		MainChatID: getEnvAsInt64("MAIN_CHAT_ID"),
		TechChatID: getEnvAsInt64("TECH_CHAT_ID"),
	}, nil
}

func getEnvAsStr(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return ""
}

func getEnvAsInt64(key string) int64 {
	valueStr := getEnvAsStr(key)
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}

	return 0
}
