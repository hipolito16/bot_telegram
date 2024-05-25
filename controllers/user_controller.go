package controllers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hipolito16/bot_telegram/bot"
	"github.com/hipolito16/bot_telegram/gemini"
	"os"
)

type UserController struct {
	bot *tgbotapi.BotAPI
}

func (self *UserController) StartResponse(update tgbotapi.Update) {
	adminUserName := os.Getenv("ADMIN_USER_NAME")
	startText := fmt.Sprintf("Olá! Eu sou um bot integrado com a API do Gemini. Para começar digite /id para saber seu ID no Telegram.\n\nApós isso, envie o ID para o meu criador @%v para que ele possa te adicionar.", adminUserName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, startText)
	self.bot.Send(msg)
}

func (self *UserController) Question(update tgbotapi.Update) {
	messages, err := gemini.Question(update.Message.Text)
	if err != nil {
		msg := bot.NewMessage(update.Message.Chat.ID, "Erro ao gerar resposta.")
		self.bot.Send(msg)
		return
	}

	for _, part := range messages {
		msg := bot.NewMessage(update.Message.Chat.ID, part)
		self.bot.Send(msg)
	}
}

func (self *UserController) Id(update tgbotapi.Update) {
	msg := bot.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Seu ID é: `%v`", update.Message.From.ID))
	self.bot.Send(msg)
}
