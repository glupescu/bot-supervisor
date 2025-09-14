package agent

import (
	"bot-supervisor/user"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	greetings = "hey"
)

func NewBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	return bot
}

func Run(ag Agent, userRoles map[int64]user.Identity) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	bot := ag.GetBot()
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.From.IsBot {
			fmt.Printf("bot %s tried to send message\n",
				update.Message.From.UserName)
			continue
		}
		userRole, err := user.GetRole(
			update.Message.From.ID, update.Message.From.FirstName, userRoles)
		if err != nil {
			fmt.Printf("Failed to get user role %v, bot %v\n",
				err, update.Message.From.IsBot)
			continue
		}
		if userRole == user.RestrictAccess {
			fmt.Printf("Disallow restrict user %v / ID %v, bot %v\n",
				update.Message.From.UserName, update.Message.From.ID, update.Message.From.IsBot)
			continue
		}
		var msg tgbotapi.MessageConfig
		var rsp string
		if len(update.Message.Text) < 3 {
			rsp = fmt.Sprintf(greetings)
		} else {
			rsp, err = ag.Serve(update.Message.Text,
				userRole, update.Message.Chat.ID)
			if err != nil {
				fmt.Printf("Got error while trying to do dispatch: %s\n", err)
				rsp = fmt.Sprintf(greetings)
			}
			if rsp == "" {
				rsp = "??? write 'help' for available commands"
			}
		}
		step := 1472
		for i := 0; i < len(rsp); i += step {
			end := min(i+step, len(rsp))

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, rsp[i:end])
			if _, err := bot.Send(msg); err != nil {
				fmt.Printf("Failed to send message %v: error %v\n", err)
				continue
			}
		}
	}
}
