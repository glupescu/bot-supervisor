package agent

import (
	"bot-supervisor/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Agent interface {
	Init(name string, bot *tgbotapi.BotAPI) error

	Serve(string, user.Type, int64) (string, error)

	GetBot() *tgbotapi.BotAPI
}
