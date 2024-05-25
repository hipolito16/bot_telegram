package gemini

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"os"
	"strings"
	"unicode/utf8"
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

	messages := resizeMessage(text)

	for i := range messages {
		reformatMessage(&messages[i])
	}

	return messages, nil
}

func resizeMessage(message string) []string {
	const limiteBytes = 4096

	if len(message) <= limiteBytes {
		return []string{message}
	}

	var result []string
	var start, current int

	for i := range message {
		current += utf8.RuneLen(rune(message[i]))
		if current > limiteBytes {
			result = append(result, message[start:i])
			start = i
			current = utf8.RuneLen(rune(message[i]))
		}
	}

	if start < len(message) {
		result = append(result, message[start:])
	}

	return result
}

func reformatMessage(message *string) {
	var builder strings.Builder
	runes := []rune(*message)
	length := len(runes)

	for i := 0; i < length; i++ {
		if runes[i] == '*' {
			before := i > 0 && runes[i-1] == '.'
			after := i < length-1 && runes[i+1] == '.'
			if !before && !after {
				continue
			}
		}
		builder.WriteRune(runes[i])
	}

	*message = builder.String()
}
