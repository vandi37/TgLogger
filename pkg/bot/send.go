package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/vanerrors"
)

// errors
const (
	ErrorSending = "error sending"
	TextEmpty    = "text empty"
)

// Reply sends a message to the chat, splitting the text if it's too long.
func (b *Bot) Send(chat int64, text string) error {
	// Check for valid text input
	if text == "" {
		return vanerrors.NewSimple(TextEmpty)
	}

	// Split the message into parts if it exceeds 4096 characters
	messageParts := strings.Split(text, "\n")

	var sendText string

	// Send the message parts iteratively
	for i := 0; i < len(messageParts); i++ {
		// Checks the length of the text
		part := messageParts[i]
		length := len(part) + len(sendText)
		if length >= 4000 || len(messageParts)-1 == i {
			// Gets the message

			// Ads text
			if len(messageParts)-1 == i {
				sendText += "\n" + part
			}

			// Creates the message
			msg := tgbotapi.NewMessage(chat, sendText)
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true

			// Sends the message
			_, err := b.Bot.Send(msg)
			if err != nil {
				return vanerrors.NewWrap(ErrorSending, err, vanerrors.EmptyHandler)
			}

			sendText = part
		} else {
			// Adding text
			sendText += "\n" + part
		}
	}
	return nil
}
