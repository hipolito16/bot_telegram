package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
)

func StartBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		panic(err)
	}

	setCommands(bot)
	return bot
}

func setCommands(bot *tgbotapi.BotAPI) {
	commands := tgbotapi.NewSetMyCommands([]tgbotapi.BotCommand{
		{Command: "add", Description: "Somente admin: Adiciona um usuário"},
		{Command: "remove", Description: "Somente admin: Remove um usuário"},
		{Command: "id", Description: "Qual meu ID?"},
	}...)

	if _, err := bot.Request(commands); err != nil {
		panic(err)
	}
}
