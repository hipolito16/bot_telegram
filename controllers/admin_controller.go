package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hipolito16/bot_telegram/database"
	"github.com/hipolito16/bot_telegram/entities"
	authmiddleware "github.com/hipolito16/bot_telegram/middlewares"
	"github.com/hipolito16/goutils"
)

type AdminController struct {
	bot *tgbotapi.BotAPI
}

func (self *AdminController) Add(update tgbotapi.Update) {
	IdTelegram := goutils.ExtractNumbers(update.Message.CommandArguments())
	if goutils.IsNilOrWhiteSpace(&IdTelegram) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Informe o ID do usuário.")
		self.bot.Send(msg)
		return
	}

	user := entities.UserEntity{IdTelegram: update.Message.From.ID}
	database.DB.Create(&user)
	authmiddleware.Auth.RefreshUsers()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usuário adicionado com sucesso.")
	self.bot.Send(msg)
}

func (self *AdminController) Remove(update tgbotapi.Update) {
	IdTelegram := goutils.ExtractNumbers(update.Message.CommandArguments())
	if goutils.IsNilOrWhiteSpace(&IdTelegram) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Informe o ID do usuário.")
		self.bot.Send(msg)
		return
	}

	var user entities.UserEntity
	tx := database.DB.First(&user, "id_telegram = ?", IdTelegram)
	if tx.RowsAffected == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usuário não encontrado.")
		self.bot.Send(msg)
		return
	}

	database.DB.Delete(&user)
	authmiddleware.Auth.RefreshUsers()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usuário removido com sucesso.")
	self.bot.Send(msg)
}
