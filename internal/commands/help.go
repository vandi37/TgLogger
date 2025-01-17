package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/TgLogger/internal/service"
	"github.com/vandi37/TgLogger/pkg/bot"
)

func Help(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		return b.Send(update.FromChat().ID, `Bot commands:

/start - Start the bot
/help - View all commands

Token:

/new_token - Add token
/delete_token - Delete token

Other:
/cancel - Cancel any input`)
	}, "help"
}

func Start(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		return b.Send(update.FromChat().ID, `This is a logger to telegram

/help - View all the commands of the bot`)
	}, "start"
}
