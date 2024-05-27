package main

import (
	_ "embed"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hipolito16/bot_telegram/bot"
	"github.com/hipolito16/bot_telegram/controllers"
	"github.com/hipolito16/bot_telegram/database"
	"github.com/hipolito16/bot_telegram/gemini"
	"github.com/hipolito16/bot_telegram/middlewares"
	"github.com/joho/godotenv"
	"os"
)

//go:embed .env
var envFile []byte

func init() {
	envMap, err := godotenv.UnmarshalBytes(envFile)
	if err != nil {
		panic(err)
	}
	for key, value := range envMap {
		if err = os.Setenv(key, value); err != nil {
			panic(err)
		}
	}
}

func main() {
	gemini.New()
	database.New()
	bot := bot.New()
	middlewares.NewMiddlewares(bot)
	controllers.NewControllers(bot)
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			controllers.Route(update.Message.Command(), update)
		}
	}
}
