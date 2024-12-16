package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// A command
type Command func(u tgbotapi.Update) error

// A map of commands
type Commands map[string]Command

func (c Commands) Run(u tgbotapi.Update) error {
	command, ok := c[u.Message.Command()]
	if !ok {
		return nil
	}
	return command(u)
}
