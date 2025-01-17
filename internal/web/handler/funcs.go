package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/TgLogger/internal/web/api"
	"github.com/vandi37/vanerrors"
)

const (
	OkSend = "successful sending message '%s' to chat %d"
)

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
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

	if !h.CheckToken(w, r) {
		return
	}

	err = h.bot.SettingSend(req.Id, req.Text, req.Mode, !req.Notifications, !req.WebPreview)
	if err != nil {
		var code = http.StatusBadRequest
		if tg_err, ok := err.(tgbotapi.Error); ok {
			code = tg_err.Code
		}
		err = api.SendError(w, code, err)
		if err != nil {
			h.logger.Errorln(err)

		}
		return
	}

	err = api.Send(w, fmt.Sprintf(OkSend, req.Text, req.Id))
	if err != nil {
		h.logger.Errorln(err)
	}
}

func (h *Handler) CheckToken(w http.ResponseWriter, r *http.Request) bool {
	ok, err := h.service.CheckToken(r.Context(), r.FormValue("token"))
	if err != nil {
		err = api.SendError(w, http.StatusInternalServerError, vanerrors.NewWrap(ErrorCheckingToken, err, vanerrors.EmptyHandler))
		if err != nil {
			h.logger.Errorln(err)
		}
		return false
	}

	if !ok {
		err = api.SendError(w, http.StatusBadRequest, vanerrors.NewSimple(TokenNotFound))
		if err != nil {
			h.logger.Errorln(err)
		}
		return false
	}
	return true
}

func (h *Handler) CheckHandler(w http.ResponseWriter, r *http.Request) {
	if h.CheckToken(w, r) {
		err := api.Send(w, "token exists")
		if err != nil {
			h.logger.Errorln(err)
		}
	}
}

func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := api.SendError(w, http.StatusNotFound, vanerrors.NewSimple(NotFound))
	if err != nil {
		h.logger.Errorln(err)
	}
}
