package main

import (
	"log"
	"nutriscan/internal/bot"
	"nutriscan/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	// Load configuration
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize and start the bot
	foodBot, err := bot.NewFoodAnalysisBot(cfg)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	foodBot.Start()
}
