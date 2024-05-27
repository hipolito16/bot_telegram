package middlewares

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hipolito16/bot_telegram/database"
	"github.com/hipolito16/bot_telegram/entities"
	"github.com/hipolito16/bot_telegram/structs"
)

var (
	Auth    *AuthMiddleware
	History *HistoryCommandMiddleware
)

func NewMiddlewares(bot *tgbotapi.BotAPI) {
	var users *[]entities.UserEntity
	database.DB.Find(&users)
	Auth = &AuthMiddleware{bot: bot, users: users}
	History = &HistoryCommandMiddleware{historiesCommand: &map[int64]structs.HistoryCommand{}}
}
