package controllers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hipolito16/bot_telegram/gemini"
)

type UserController struct {
	bot *tgbotapi.BotAPI
}

func (self *UserController) Question(update tgbotapi.Update) {
	response, err := gemini.Question(update.Message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Erro ao gerar resposta.")
		self.bot.Send(msg)
		return
	}

	for _, part := range response {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, part)
		msg.ParseMode = "markdown"
		self.bot.Send(msg)
	}
}

func (self *UserController) Id(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Seu ID é: `%v`", update.Message.From.ID))
	msg.ParseMode = "markdown"
	self.bot.Send(msg)
}

func (self *UserController) DefaultResponse(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Desculpe, não entendi o que você quis dizer.")
	self.bot.Send(msg)
}
