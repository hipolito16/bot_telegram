package middlewares

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hipolito16/bot_telegram/bot"
	"github.com/hipolito16/bot_telegram/database"
	"github.com/hipolito16/bot_telegram/entities"
)

type AuthMiddleware struct {
	bot   *tgbotapi.BotAPI
	users *[]entities.UserEntity
}

func (auth *AuthMiddleware) IsAdmin(update tgbotapi.Update) bool {
	for _, user := range *auth.users {
		if user.IdTelegram == update.Message.From.ID && user.Admin {
			return true
		}
	}
	msg := bot.NewMessage(update.Message.Chat.ID, "Você não tem permissão para executar essa ação.")
	auth.bot.Send(msg)
	return false
}

func (auth *AuthMiddleware) VerifyUser(update tgbotapi.Update) bool {
	for _, user := range *auth.users {
		if user.IdTelegram == update.Message.From.ID {
			return true
		}
	}
	msg := bot.NewMessage(update.Message.Chat.ID, "Você não está autorizado a usar esse bot.")
	auth.bot.Send(msg)
	return false
}

func (auth *AuthMiddleware) RefreshUsers() {
	database.DB.Find(&auth.users)
}
