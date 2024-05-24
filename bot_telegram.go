package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hipolito16/bot_telegram/controllers"
	"github.com/hipolito16/bot_telegram/controllers/bot"
	"github.com/hipolito16/bot_telegram/database"
	"github.com/hipolito16/bot_telegram/gemini"
	"github.com/hipolito16/bot_telegram/middlewares"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
}

func main() {
	gemini.Start()
	database.Start()
	bot := bot.StartBot()
	middlewares.StartMiddlewares(bot)
	controllers.StartControllers(bot)

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Command() {
			case "add":
				if middlewares.Auth.IsAdmin(update) {
					controllers.Admin.Add(update)
				}
			case "remove":
				if middlewares.Auth.IsAdmin(update) {
					controllers.Admin.Remove(update)
				}
			case "id":
				controllers.User.Id(update)
			default:
				if middlewares.Auth.VerifyUser(update) {
					controllers.User.Question(update)
				}
			}
		}
	}
}
