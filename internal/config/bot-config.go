package config

import (
	"errors"
	"os"
)

type BotConfig struct {
	TelegramToken string
	OpenAIToken   string
}

func Load() (*BotConfig, error) {
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	openAIToken := os.Getenv("OPENAI_API_KEY")

	if telegramToken == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN is not set")
	}

	if openAIToken == "" {
		return nil, errors.New("OPENAI_API_KEY is not set")
	}

	return &BotConfig{
		TelegramToken: telegramToken,
		OpenAIToken:   openAIToken,
	}, nil
}
