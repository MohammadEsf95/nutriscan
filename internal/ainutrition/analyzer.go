package ainutrition

import (
	"context"
	"encoding/base64"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type Analyzer struct {
	client *openai.Client
}

func NewAnalyzer(client *openai.Client) *Analyzer {
	return &Analyzer{client: client}
}

func (a *Analyzer) AnalyzeFood(ctx context.Context, imageData []byte, description string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a food nutrition calculator. Users will send pictures of their meals, text descriptions, or both. You must calculate the total calories, protein, carbohydrates, fat, fiber, and sugar in this exact order. Respond in Farsi only. If the input is irrelevant or unrecognizable, send a standard error message in Farsi and ask the user to upload the information again. Be approximate and avoid being overly exact. Do not ask for additional details.",
		},
	}

	// If image is provided, convert to base64
	var base64Image string
	if len(imageData) > 0 {
		base64Image = fmt.Sprintf("data:image/jpeg;base64,%s", toBase64(imageData))
		messages = append(messages, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser,
			MultiContent: []openai.ChatMessagePart{
				{
					Type: openai.ChatMessagePartTypeText,
					Text: "Here's an image of the food.",
				},
				{
					Type: openai.ChatMessagePartTypeImageURL,
					ImageURL: &openai.ChatMessageImageURL{
						URL: base64Image,
					},
				},
			},
		})
	}

	// Add text description
	if description != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Food description: %s", description),
		})
	}

	resp, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: messages,
	})

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func toBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
