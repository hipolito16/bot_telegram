package structs

import "github.com/hipolito16/bot_telegram/entities"

type HistoryCommand struct {
	Command string
	User    entities.UserEntity
}
