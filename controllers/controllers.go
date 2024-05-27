package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hipolito16/bot_telegram/middlewares"
)

var (
	Admin *AdminController
	User  *UserController
)

func NewControllers(bot *tgbotapi.BotAPI) {
	Admin = &AdminController{bot: bot}
	User = &UserController{bot: bot}
}

func Route(command string, update tgbotapi.Update) {
	if historyCommand, ok, _ := middlewares.History.HasCommandHistory(update); ok {
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
	case "limparchat":
		if middlewares.Auth.VerifyUser(update) {
			User.LimparChat(update)
		}
	default:
		if middlewares.Auth.VerifyUser(update) {
			User.Chat(update)
		}
	}
}
