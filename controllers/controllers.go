package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hipolito16/bot_telegram/middlewares"
	"github.com/hipolito16/bot_telegram/structs"
)

var (
	History *HistoryCommandController
	Admin   *AdminController
	User    *UserController
)

func StartControllers(bot *tgbotapi.BotAPI) {
	History = &HistoryCommandController{historiesCommand: &map[int64]structs.HistoryCommand{}}
	Admin = &AdminController{bot: bot}
	User = &UserController{bot: bot}
}

func Route(command string, update tgbotapi.Update) {
	if historyCommand, ok, _ := History.HasCommandHistory(update); ok {
		command = historyCommand
	}

	switch command {
	case "start":
		User.StartResponse(update)
	case "list":
		if middlewares.Auth.IsAdmin(update) {
			Admin.List(update)
		}
	case "add":
		if middlewares.Auth.IsAdmin(update) {
			Admin.Add(update)
		}
	case "remove":
		if middlewares.Auth.IsAdmin(update) {
			Admin.Remove(update)
		}
	case "id":
		User.Id(update)
	default:
		if middlewares.Auth.VerifyUser(update) {
			User.Question(update)
		}
	}
}
