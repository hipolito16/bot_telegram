package controllers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hipolito16/bot_telegram/bot"
	"github.com/hipolito16/bot_telegram/database"
	"github.com/hipolito16/bot_telegram/entities"
	authmiddleware "github.com/hipolito16/bot_telegram/middlewares"
	"github.com/hipolito16/goutils"
	"strconv"
)

type AdminController struct {
	bot *tgbotapi.BotAPI
}

func (self *AdminController) List(update tgbotapi.Update) {
	var users []entities.UserEntity
	if tx := database.DB.Find(&users); tx.RowsAffected == 0 {
		msg := bot.NewMessage(update.Message.Chat.ID, "Nenhum usuário cadastrado.")
		self.bot.Send(msg)
		return
	}

	msg := bot.NewMessage(update.Message.Chat.ID, "Usuários cadastrados:")
	for _, user := range users {
		msg.Text += fmt.Sprintf("\nID: `%v` | Nome: `%v`", user.IdTelegram, user.Name)
	}
	self.bot.Send(msg)
}
func (self *AdminController) Add(update tgbotapi.Update) {
	if _, ok, user := History.HasCommandHistory(update); !ok {
		IdTelegramString := goutils.ExtractNumbers(update.Message.CommandArguments())
		if goutils.IsNilOrWhiteSpace(&IdTelegramString) {
			msg := bot.NewMessage(update.Message.Chat.ID, "Nenhum ID informado.")
			self.bot.Send(msg)
			return
		}

		if user.IdTelegram, _ = strconv.ParseInt(IdTelegramString, 10, 64); user.IdTelegram == 0 {
			msg := bot.NewMessage(update.Message.Chat.ID, "ID inválido.")
			self.bot.Send(msg)
			return
		}

		if tx := database.DB.First(&user, "id_telegram = ?", user.IdTelegram); tx.RowsAffected > 0 {
			msg := bot.NewMessage(update.Message.Chat.ID, "Usuário já cadastrado.")
			self.bot.Send(msg)
			return
		}

		History.AddCommandHistory(update, user)
		msg := bot.NewMessage(update.Message.Chat.ID, "Informe o nome para esse usuário.")
		self.bot.Send(msg)
	} else {
		if goutils.IsNilOrWhiteSpace(&update.Message.Text) {
			msg := bot.NewMessage(update.Message.Chat.ID, "Nenhum nome informado.")
			self.bot.Send(msg)
			return
		}

		user.Name = update.Message.Text
		database.DB.Create(&user)
		authmiddleware.Auth.RefreshUsers()
		History.RemoveCommandHistory(update)
		msg := bot.NewMessage(update.Message.Chat.ID, "Usuário cadastrado com sucesso.")
		self.bot.Send(msg)
	}
}

func (self *AdminController) Remove(update tgbotapi.Update) {
	IdTelegram := goutils.ExtractNumbers(update.Message.CommandArguments())
	if goutils.IsNilOrWhiteSpace(&IdTelegram) {
		msg := bot.NewMessage(update.Message.Chat.ID, "Informe o ID do usuário.")
		self.bot.Send(msg)
		return
	}

	var user entities.UserEntity
	tx := database.DB.First(&user, "id_telegram = ?", IdTelegram)
	if tx.RowsAffected == 0 {
		msg := bot.NewMessage(update.Message.Chat.ID, "Usuário não encontrado.")
		self.bot.Send(msg)
		return
	}

	database.DB.Delete(&user)
	authmiddleware.Auth.RefreshUsers()
	msg := bot.NewMessage(update.Message.Chat.ID, "Usuário removido com sucesso.")
	self.bot.Send(msg)
}
