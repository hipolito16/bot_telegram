package entities

import (
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"sync"
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
	genaiContentChan := make(chan *genai.Content, len(contentStructs))
	indexChange := make(chan int, len(contentStructs))
	wg := sync.WaitGroup{}
	genaiContents := make([]*genai.Content, len(contentStructs))
	for index, contentStruct := range contentStructs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			genaiContentChan <- &genai.Content{
				Role: contentStruct.Role,
				Parts: []genai.Part{
					genai.Text(contentStruct.Content),
				},
			}
			indexChange <- index
		}()
	}
	go func() {
		wg.Wait()
		close(genaiContentChan)
		close(indexChange)
	}()
	for index := range indexChange {
		genaiContents[index] = <-genaiContentChan
	}
	return genaiContents
}

func (self *GeminiChatHistoryEntity) FromGenaiContent(history []*genai.Content) {
	if len(history) > 100 {
		history = history[len(history)-100:]
	}
	contentChan := make(chan Content, len(history))
	indexChange := make(chan int, len(history))
	var wg sync.WaitGroup
	contents := make([]Content, len(history))
	for index, content := range history {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var text string
			for _, part := range content.Parts {
				text += fmt.Sprintf("%v", part)
			}
			contentChan <- Content{
				Content: text,
				Role:    content.Role,
			}
			indexChange <- index
		}()
	}
	go func() {
		wg.Wait()
		close(contentChan)
		close(indexChange)
	}()
	for index := range indexChange {
		contents[index] = <-contentChan
	}
	historyJson, _ := json.Marshal(&contents)
	self.History = historyJson
}
