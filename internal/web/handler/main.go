package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vandi37/TgLogger/internal/web/api"
	"github.com/vandi37/TgLogger/pkg/bot"
	"github.com/vandi37/TgLogger/pkg/logger"
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
	mainPath string
	bot      *bot.Bot
}

// Created a new handler
func New(bot *bot.Bot, service *service.Service, logger *logger.Logger) *Handler {
	// Creating handler
	handler := Handler{
		service:  service,
		mainPath: "/send",
		bot:      bot,
	}

	return &handler
}

// Serve
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// The header
	w.Header().Add("Content-Type", "application/json")

	// Gets the handler func
	if r.URL.Path != h.mainPath {

		// Not found
		err := api.SendError(w, http.StatusNotFound, vanerrors.NewSimple(NotFound, fmt.Sprintf("the only allowed page is %s", h.mainPath)))
		if err != nil {
			h.logger.Errorln(err)
		}

		return
	}

	if r.Method != http.MethodPost {
		// Not found
		err := api.SendError(w, http.StatusMethodNotAllowed, vanerrors.NewSimple(MethodNotAllowed, "allowed only method post"))
		if err != nil {
			h.logger.Errorln(err)
		}

		return
	}

	// Runs the handler
	defer r.Body.Close()

	var req api.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		err = api.SendError(w, http.StatusBadRequest, vanerrors.NewWrap(InvalidRequest, err, vanerrors.EmptyHandler))
		if err != nil {
			h.logger.Errorln(err)

		}
		return
	}

	ok, err := h.service.CheckToken(req.Token)
	if err != nil {
		err = api.SendError(w, http.StatusInternalServerError, vanerrors.NewWrap(ErrorCheckingToken, err, vanerrors.EmptyHandler))
		if err != nil {
			h.logger.Errorln(err)

		}
		return
	}

	if !ok {
		err = api.SendError(w, http.StatusBadRequest, vanerrors.NewSimple(TokenNotFound))
		if err != nil {
			h.logger.Errorln(err)

		}
		return
	}

	err = h.bot.Send(req.Id, req.Text)
	if err != nil {
		err = api.SendError(w, http.StatusInternalServerError, vanerrors.NewWrap(ErrorSendingMessage, err, vanerrors.EmptyHandler))
		if err != nil {
			h.logger.Errorln(err)

		}
		return
	}
	err = api.Send(w, fmt.Sprintf("successful sending message '%s' to chat %d", req.Text, req.Id))
	if err != nil {
		h.logger.Errorln(err)

	}
}
