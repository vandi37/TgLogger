package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/vandi37/TgLogger/internal/application"
)

func main() {
	// Creating a new application with a hour timeout
	app := application.New("config/config.yml")

	// Adding graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Running the app
	app.Run(ctx)
}
