package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) NewToken(u tgbotapi.Update) error {
	id := u.SentFrom().ID
	token, err := b.service.AddToken(id)
	if err != nil {
		return b.Send(id, fmt.Sprintf("❌ Error creating token: %v", err))
	}
	return b.Send(id, fmt.Sprintf("✅ Created token: `%s`", token))
}

func (b *Bot) DeleTeToken(u tgbotapi.Update) error {
	id := u.SentFrom().ID
	value := u.Message.CommandArguments()
	err := b.service.DeleteToken(value, id)
	if err != nil {
		return b.Send(id, fmt.Sprintf("❌ Error deleting token: %v", err))
	}
	return b.Send(id, fmt.Sprintf("✅ Deleted token: `%s`", value))
}

func (b *Bot) MyTokens(u tgbotapi.Update) error {
	id := u.SentFrom().ID
	tokens, err := b.service.GetTokens(id)
	if err != nil {
		return b.Send(id, fmt.Sprintf("❌ Error getting tokens: %v", err))
	}
	if len(tokens) <= 0 {
		return b.Send(id, "❌ You don't have any tokens")
	}
	var mes string = "✅ Your tokens:"
	for i, t := range tokens {
		mes += fmt.Sprintf("\n%d: `%s`", i+1, t)
	}
	return b.Send(id, mes)
}
