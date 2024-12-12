package bot

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *FoodAnalysisBot) handlePhoto(userID int64, message *tgbotapi.Message) {
	photoSize := (message.Photo)[len(message.Photo)-1]
	fileConfig := tgbotapi.FileConfig{FileID: photoSize.FileID}
	file, err := b.telegramBot.GetFile(fileConfig)
	if err != nil {
		log.Println("Error getting file:", err)
		return
	}

	imageData, err := downloadFile(file.FilePath)
	if err != nil {
		log.Println("Error downloading file:", err)
		return
	}

	state := b.userStates.GetOrCreateState(userID)
	state.PendingImage = imageData
	b.userStates.UpdateState(userID, state)
}

func (b *FoodAnalysisBot) handleText(userID int64, message *tgbotapi.Message) {
	state := b.userStates.GetOrCreateState(userID)

	// Check if waiting for more details
	if state.WaitingForDetail {
		state.PendingText += " " + message.Text
		state.WaitingForDetail = false
		b.userStates.UpdateState(userID, state)
	} else {
		state.PendingText = message.Text
	}

	// Analyze food if we have enough information
	if len(state.PendingImage) > 0 || state.PendingText != "" {
		result, err := b.analyzer.AnalyzeFood(context.Background(), state.PendingImage, state.PendingText)

		if err != nil {
			b.sendMessage(userID, "Error analyzing food: "+err.Error())
		} else {
			b.sendMessage(userID, result)
		}

		// Reset state
		state.PendingImage = nil
		state.PendingText = ""
		b.userStates.UpdateState(userID, state)
	}
}

func (b *FoodAnalysisBot) sendMessage(userID int64, text string) {
	msg := tgbotapi.NewMessage(userID, text)
	_, err := b.telegramBot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message to user %d: %v", userID, err)
	}
}

func downloadFile(filePath string) ([]byte, error) {
	url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", os.Getenv("TELEGRAM_BOT_TOKEN"), filePath)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
