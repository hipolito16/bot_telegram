package controllers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	Admin *AdminController
	User  *UserController
)

func StartControllers(bot *tgbotapi.BotAPI) {
	Admin = &AdminController{bot: bot}
	User = &UserController{bot: bot}
}

func extractCommand(update tgbotapi.Update) string {
	return update.Message.Command()
}
