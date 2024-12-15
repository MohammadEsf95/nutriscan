package main

import (
	"log"
	"nutriscan/internal/bot"
	"nutriscan/internal/config"
	"nutriscan/internal/database"
	"nutriscan/internal/users"

	"github.com/joho/godotenv"
)

func main() {
	// Load configuration
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize database
	dbConfig, err := database.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.AutoMigrate(&users.User{}); err != nil {
		log.Fatal(err)
	}

	userRepository := users.NewUserRepository(db)
	userService := users.NewUserService(userRepository)
	userHandler := users.NewUserHandler(userService)

	// Initialize bot configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize and start the bot
	foodBot, err := bot.NewFoodAnalysisBot(cfg, userHandler)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	foodBot.Start()
}
