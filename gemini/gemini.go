package gemini

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/generative-ai-go/genai"
	"github.com/hipolito16/bot_telegram/database"
	"github.com/hipolito16/bot_telegram/entities"
	"google.golang.org/api/option"
	"os"
	"strings"
	"unicode/utf8"
)

var apiKey string

func New() {
	apiKey = os.Getenv("GEMINI_API_KEY")
}

func Chat(update tgbotapi.Update) ([]string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-1.5-pro")
	cs := model.StartChat()
	var geminiChatHistory entities.GeminiChatHistoryEntity
	database.DB.FirstOrCreate(&geminiChatHistory, entities.GeminiChatHistoryEntity{IdTelegram: update.Message.From.ID})
	cs.History = geminiChatHistory.ToGenaiContent()
	response, err := cs.SendMessage(ctx, genai.Text(update.Message.Text))
	if err != nil {
		return nil, err
	}
	geminiChatHistory.FromGenaiContent(cs.History)
	database.DB.Save(&geminiChatHistory)
	var textBuilder strings.Builder
	for _, candidate := range response.Candidates {
		for _, part := range candidate.Content.Parts {
			textBuilder.WriteString(fmt.Sprintf("%v", part))
		}
	}
	messages := resizeMessage(textBuilder.String())
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
	*message = strings.ReplaceAll(*message, "*", "")
}
