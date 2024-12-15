package application

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/vandi37/TgLogger/config"
	"github.com/vandi37/TgLogger/pkg/closer"
	"github.com/vandi37/TgLogger/pkg/db"
	"github.com/vandi37/TgLogger/pkg/logger"
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
	logger := logger.New(os.Stderr)

	// Loading config
	cfg, err := config.Get(a.Config)
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Println("got config")

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

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	closer.Close(ctx)

	// The program end
	logger.Println("application stopped")

	os.Exit(http.StatusTeapot)
}