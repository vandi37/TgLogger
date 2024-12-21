package application

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/vandi37/TgLogger/config"
	"github.com/vandi37/TgLogger/internal/web/handler"
	"github.com/vandi37/TgLogger/internal/web/server"
	"github.com/vandi37/TgLogger/pkg/bot"
	"github.com/vandi37/TgLogger/pkg/closer"
	"github.com/vandi37/TgLogger/pkg/db"
	"github.com/vandi37/TgLogger/pkg/logger"
	"github.com/vandi37/TgLogger/pkg/service"
)

// Thr application program
type Application struct {
	Config string
}

// Creates a new application
func New(config string) *Application {
	return &Application{Config: config}
}

// Runs the application
func (a *Application) Run(ctx context.Context) {
	// Creates logger
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	logger := logger.New(file)

	// Loading config
	cfg, err := config.Get(a.Config)
	if err != nil {
		logger.Fatalln(err)
	}

	// The program

	// Creating closer
	closer := closer.New(logger)

	// connecting to the database
	db, err := db.New(cfg.DB)
	if err != nil {
		logger.Fatalln(err)
	}
	closer.Add(db.Close)

	// Creating tables
	err = db.Init()
	if err != nil {
		logger.Fatalln(err)
	}

	service := service.New(db, logger)

	bot, err := bot.New(cfg.Token, service, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	go bot.Run(ctx)

	handler := handler.New(bot, service, logger)

	server := server.New(handler, cfg.Port)

	go server.Run()
	closer.Add(server.Close)

	logger.Printf("application started on port %d", cfg.Port)

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	closer.Close(ctx)

	// The program end
	logger.Println("application stopped")

	os.Exit(http.StatusTeapot)
}
