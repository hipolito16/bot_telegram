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

func NewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "markdown"
	return msg
}

func setCommands(bot *tgbotapi.BotAPI) {
	commands := tgbotapi.NewSetMyCommands([]tgbotapi.BotCommand{
		{Command: "start", Description: "Olá! Eu sou um bot integrado com a API do Gemini. Como posso te ajudar?"},
		{Command: "list", Description: "Lista todos usuários."},
		{Command: "add", Description: "Adiciona um usuário."},
		{Command: "remove", Description: "Remove um usuário."},
		{Command: "id", Description: "Retorna seu ID no Telegram."},
	}...)

	if _, err := bot.Request(commands); err != nil {
		panic(err)
	}
}
