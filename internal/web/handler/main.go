package handler

import (
	"net/http"

	"github.com/vandi37/TgLogger/pkg/bot"
	"github.com/vandi37/TgLogger/pkg/logger"
	"github.com/vandi37/TgLogger/pkg/service"
)

// Errors
const (
	NotFound            = "this page is unavailable"
	InvalidRequest      = "invalid request"
	ErrorCheckingToken  = "error checking token"
	TokenNotFound       = "token does not exist"
	ErrorSendingMessage = "error sending message"
	MethodNotAllowed    = "method not allowed"
)

// The handler
type Handler struct {
	*http.ServeMux
	logger  *logger.Logger
	service *service.Service
	bot     *bot.Bot
}

// Created a new handler
func New(bot *bot.Bot, service *service.Service, logger *logger.Logger) *Handler {
	// Creating handler
	handler := Handler{
		ServeMux: http.NewServeMux(),
		service:  service,
		logger:   logger,
		bot:      bot,
	}
	handler.HandleFunc("/api/send", ContentType(handler.CheckMethod(http.MethodPost, handler.Send)))
	handler.HandleFunc("/api/check/{x}", ContentType(handler.CheckMethod(http.MethodGet, handler.CheckHandler)))
	handler.HandleFunc("/", ContentType(handler.NotFoundHandler))

	return &handler
}
