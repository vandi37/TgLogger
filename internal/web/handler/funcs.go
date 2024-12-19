package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vandi37/TgLogger/internal/web/api"
	"github.com/vandi37/TgLogger/pkg/logger"
	"github.com/vandi37/TgLogger/pkg/service"
	"github.com/vandi37/vanerrors"
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

	if !CheckToken(w, req.Token, h.service, h.logger) {
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

func CheckToken(w http.ResponseWriter, token string, service *service.Service, logger *logger.Logger) bool {
	ok, err := service.CheckToken(token)
	if err != nil {
		err = api.SendError(w, http.StatusInternalServerError, vanerrors.NewWrap(ErrorCheckingToken, err, vanerrors.EmptyHandler))
		if err != nil {
			logger.Errorln(err)
		}
		return false
	}

	if !ok {
		err = api.SendError(w, http.StatusBadRequest, vanerrors.NewSimple(TokenNotFound))
		if err != nil {
			logger.Errorln(err)

		}
		return false
	}
	return true
}

type sendHandler struct {
	logger  *logger.Logger
	service *service.Service
}

func newSendHandler(service *service.Service, logger *logger.Logger) *sendHandler {
	return &sendHandler{logger, service}
}

func (h *sendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	CheckMethod(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/check/") {
			err := api.SendError(w, http.StatusBadRequest, vanerrors.NewSimple(TokenNotFound))
			if err != nil {
				h.logger.Errorln(err)
			}
			return
		}

		if CheckToken(w, r.URL.Path[11:], h.service, h.logger) {
			err := api.Send(w, "token exists")
			if err != nil {
				h.logger.Errorln(err)
			}
		}
	}, h.logger)(w, r)
}
