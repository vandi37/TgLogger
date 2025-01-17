package bot

import (
	"context"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/TgLogger/internal/service"
	"github.com/vandi37/TgLogger/pkg/commands"
	"github.com/vandi37/TgLogger/pkg/logger"
	"github.com/vandi37/vanerrors"
)

// Errors
const (
	ErrorGettingBot = "error getting bot"
	ContextExit     = "context exit"
)

type Command func(u *tgbotapi.Update) error

// The telegram bot
type Bot struct {
	Bot      *tgbotapi.BotAPI
	logger   *logger.Logger
	mu       sync.Mutex
	upd      tgbotapi.UpdateConfig
	service  *service.Service
	commands commands.Commands
	// rootHandler *handler.Handler
}

// Creates a new bot
func New(token string, service *service.Service, logger *logger.Logger) (*Bot, error) {
	// New bor api
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorGettingBot, err, vanerrors.EmptyHandler)
	}

	// Creates an update
	u := tgbotapi.NewUpdate(60)

	res := Bot{
		Bot:     bot,
		logger:  logger,
		upd:     u,
		service: service,
	}
	res.commands = commands.Commands{
		"new_token":    res.NewToken,
		"delete_token": res.DeleTeToken,
		"my_tokens":    res.MyTokens,
	}

	return &res, nil
}

// Runs the bot
func (b *Bot) Run(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Gets updates
	updates := b.Bot.GetUpdatesChan(b.upd)
	updates.Clear()

	for {
		select {
		case <-ctx.Done():
			return vanerrors.NewSimple(ContextExit)
		case update := <-updates:
			if update.Message == nil || !update.Message.IsCommand() {
				continue
			}

			err := b.service.NewUser(ctx, update.SentFrom().ID)
			if err != nil {
				b.logger.Errorln(err)
				continue
			}

			err = b.commands.Run(update)
			if err != nil {
				b.logger.Errorln(err)
				continue
			}
		}
	}
}
