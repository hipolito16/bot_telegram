package middlewares

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hipolito16/bot_telegram/entities"
	"github.com/hipolito16/bot_telegram/structs"
)

type HistoryCommandMiddleware struct {
	historiesCommand *map[int64]structs.HistoryCommand
}

func (self *HistoryCommandMiddleware) HasCommandHistory(update tgbotapi.Update) (command string, ok bool, user entities.UserEntity) {
	if commandHistory, ok := (*self.historiesCommand)[update.Message.Chat.ID]; ok {
		return commandHistory.Command, true, commandHistory.User
	} else {
		return "", false, entities.UserEntity{}
	}
}

func (self *HistoryCommandMiddleware) AddCommandHistory(update tgbotapi.Update, user entities.UserEntity) {
	if commandHistory, ok := (*self.historiesCommand)[update.Message.Chat.ID]; ok {
		commandHistory.Command = update.Message.Command()
		commandHistory.User = user
	} else {
		(*self.historiesCommand)[update.Message.Chat.ID] = structs.HistoryCommand{
			Command: update.Message.Command(),
			User:    user,
		}
	}
}

func (self *HistoryCommandMiddleware) RemoveCommandHistory(update tgbotapi.Update) {
	if _, ok := (*self.historiesCommand)[update.Message.Chat.ID]; ok {
		delete(*self.historiesCommand, update.Message.Chat.ID)
	}
}
