package commands

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/TgLogger/internal/service"
	"github.com/vandi37/TgLogger/pkg/bot"
)

func NewToken(b *bot.Bot, s *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		id := update.SentFrom().ID
		token, err := s.AddToken(ctx, id)
		if err != nil {
			return b.Send(id, fmt.Sprintf("❌ Error creating token: %v", err))
		}
		return b.Send(id, fmt.Sprintf("✅ Created token: `%s`", token))
	}, "new_token"
}

func DeleTeToken(b *bot.Bot, s *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		id := update.SentFrom().ID

		err := b.Send(id, "‼️ Please send token to delete")
		if err != nil {
			return err
		}
		wait, cancel := b.Waiter.Add(update.SentFrom().ID)
		defer b.Waiter.Remove(update.SentFrom().ID)

		select {
		case <-cancel.Canceled():
			return nil
		case <-ctx.Done():
			return nil
		case answer := <-wait:
			err := s.DeleteToken(ctx, answer.Text, id)
			if err != nil {
				return b.Send(id, fmt.Sprintf("❌ Error deleting token: %v", err))
			}
			return b.Send(id, fmt.Sprintf("✅ Deleted token: `%s`", answer.Text))
		}
	}, "delete_token"

}

func MyTokens(b *bot.Bot, s *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		id := update.SentFrom().ID
		tokens, err := s.GetTokens(ctx, id)
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
	}, "my_tokens"

}
