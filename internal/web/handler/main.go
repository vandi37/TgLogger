package handler

import (
	"net/http"

	"github.com/vandi37/TgLogger/internal/web/api"
	"github.com/vandi37/TgLogger/pkg/bot"
	"github.com/vandi37/TgLogger/pkg/logger"
	"github.com/vandi37/TgLogger/pkg/maps"
	"github.com/vandi37/TgLogger/pkg/service"
	"github.com/vandi37/vanerrors"
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
	logger   *logger.Logger
	service  *service.Service
	funcs    map[string]http.HandlerFunc
	handlers map[string]http.Handler
	bot      *bot.Bot
}

// Created a new handler
func New(bot *bot.Bot, service *service.Service, logger *logger.Logger) *Handler {
	// Creating handler
	handler := Handler{
		service: service,
		bot:     bot,
	}
	handler.funcs = map[string]http.HandlerFunc{
		"/api/send": CheckMethod(http.MethodPost, handler.Send, logger),
	}
	handler.handlers = map[string]http.Handler{
		"/api/check": newSendHandler(service, logger),
	}

	return &handler
}

// Serve
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	fn, ok := h.funcs[r.URL.Path]
	if ok && fn != nil {
		fn(w, r)
		return
	}

	hand, ok := maps.Find(h.handlers, r.URL.Path, "/")
	if ok && hand != nil {
		hand.ServeHTTP(w, r)
		return
	}

	err := api.SendError(w, http.StatusNotFound, vanerrors.NewSimple(NotFound))
	if err != nil {
		h.logger.Errorln(err)
		return
	}

}
