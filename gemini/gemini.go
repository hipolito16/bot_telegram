package gemini

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"os"
)

var apiKey string

func Start() {
	apiKey = os.Getenv("GEMINI_API_KEY")
}

func Question(question string) ([]string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	defer client.Close()

	model := client.GenerativeModel("gemini-pro")

	response, err := model.GenerateContent(ctx, genai.Text(question))
	if err != nil {
		return nil, err
	}

	var text string
	for _, candidate := range response.Candidates {
		for _, part := range candidate.Content.Parts {
			text += fmt.Sprintf("%v", part)
		}
	}

	return resizeMessage(text), nil
}

func resizeMessage(text string) []string {
	const maxMessageLength = 4096

	if len(text) <= maxMessageLength {
		return []string{text}
	}

	var messages []string
	for i := 0; i < len(text); i += maxMessageLength {
		end := min(i+maxMessageLength, len(text))
		messages = append(messages, text[i:end])
	}

	return messages
}
