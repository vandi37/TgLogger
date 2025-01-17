package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/vanerrors"
)

// errors
const (
	ErrorSending = "error sending"
	TextEmpty    = "text empty"
)

func (b *Bot) Send(chat int64, text string) error {
	return b.SettingSend(chat, text, "", false, true)
}

func (b *Bot) SettingSend(chat int64, text string, mode string, disable_notifications, disable_preview bool) error {
	if text == "" {
		return vanerrors.NewSimple(TextEmpty)
	}

	msg := tgbotapi.NewMessage(chat, text)
	msg.DisableNotification = disable_notifications
	msg.DisableWebPagePreview = disable_preview
	msg.ParseMode = mode

	_, err := b.Bot.Send(msg)
	return err
}
