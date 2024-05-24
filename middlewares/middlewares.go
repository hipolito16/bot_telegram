package middlewares

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hipolito16/bot_telegram/database"
	"github.com/hipolito16/bot_telegram/entities"
)

var Auth *AuthMiddleware

func StartMiddlewares(bot *tgbotapi.BotAPI) {
	var users *[]entities.UserEntity
	database.DB.Find(&users)
	Auth = &AuthMiddleware{bot: bot, users: users}
}
