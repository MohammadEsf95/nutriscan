package bot

import (
	"nutriscan/internal/ainutrition"
	"nutriscan/internal/config"
	"nutriscan/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	openai "github.com/sashabaranov/go-openai"
)

type FoodAnalysisBot struct {
	telegramBot  *tgbotapi.BotAPI
	openaiClient *openai.Client
	userStates   *storage.UserStateManager
	analyzer     *ainutrition.Analyzer
}

func NewFoodAnalysisBot(cfg *config.BotConfig) (*FoodAnalysisBot, error) {
	// Initialize Telegram Bot
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	// Initialize OpenAI Client
	openaiClient := openai.NewClient(cfg.OpenAIToken)

	return &FoodAnalysisBot{
		telegramBot:  bot,
		openaiClient: openaiClient,
		userStates:   storage.NewUserStateManager(),
		analyzer:     ainutrition.NewAnalyzer(openaiClient),
	}, nil
}

func (b *FoodAnalysisBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.telegramBot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go b.handleUpdate(update)
	}
}

func (b *FoodAnalysisBot) handleUpdate(update tgbotapi.Update) {
	userID := update.Message.From.ID

	// Handle photo and text inputs
	if update.Message.Photo != nil {
		b.handlePhoto(userID, update.Message)
	}

	if update.Message.Text != "" {
		b.handleText(userID, update.Message)
	}
}
