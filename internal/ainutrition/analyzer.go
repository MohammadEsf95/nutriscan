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
			Role: openai.ChatMessageRoleSystem,
			Content: "You are a food nutrition calculator. Users will send pictures of their meals, text descriptions, or both, and you should estimate the amount of food from the photo, and consider the extra description provided from the text. Based on the input: \n" +
				"1. If the meal and description are clear, return only the calculated values in this exact format: Total Calories: (calculated)\n Protein: (calculated) grams\n Carbohydrates: (calculated) grams\n Fat: (calculated) grams\n Fiber: (calculated) grams\n Sugar: (calculated) grams. \n" +
				"2. If the meal was unrecognizable, return the following exact message in Farsi: اطلاعات کافی نیست، لطفاً اطلاعات را دوباره ارسال کنید. \n" +
				"3. If the input is irrelevant or unrecognizable, return the following exact message in Farsi: دستور مرتبط نیست. \n" +
				"Be approximate and avoid being overly exact in calculations. Do not ask for additional details or include extra information in your responses.",
		},
	}

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
