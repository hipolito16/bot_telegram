package entities

import (
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type GeminiChatHistoryEntity struct {
	IdGeminiChatHistory uint `gorm:"primarykey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	IdTelegram          int64          `gorm:"not null"`
	History             datatypes.JSON
}

func (GeminiChatHistoryEntity) TableName() string {
	return "gemini_chat_histories"
}

type Content struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

func (self *GeminiChatHistoryEntity) ToGenaiContent() []*genai.Content {
	var contentStructs []Content
	json.Unmarshal(self.History, &contentStructs)
	history := make([]*genai.Content, len(contentStructs))
	for i, contentStruct := range contentStructs {
		history[i] = &genai.Content{
			Role: contentStruct.Role,
			Parts: []genai.Part{
				genai.Text(contentStruct.Content),
			},
		}
	}
	return history
}

func (self *GeminiChatHistoryEntity) FromGenaiContent(history []*genai.Content) {
	if len(history) > 100 {
		history = history[len(history)-100:]
	}
	var contents []Content
	for _, content := range history {
		var text string
		for _, part := range content.Parts {
			text += fmt.Sprintf("%v", part)
		}
		contents = append(contents, Content{
			Role:    content.Role,
			Content: text,
		})
	}
	historyJson, _ := json.Marshal(&contents)
	self.History = historyJson
}
