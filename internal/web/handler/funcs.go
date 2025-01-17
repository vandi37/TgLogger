package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vandi37/TgLogger/internal/service"
	"github.com/vandi37/TgLogger/internal/web/api"
	"github.com/vandi37/TgLogger/pkg/logger"
	"github.com/vandi37/vanerrors"
)

const (
	OkSend          = "successful sending message '%s' to chat %d"
	OkSendWithToken = OkSend + " from token %s"
	OkCheck         = "successful token checking, token %s exists"
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

	h.logger.Printf(OkSendWithToken, req.Text, req.Id, req.Token)
	err = api.Send(w, fmt.Sprintf(OkSend, req.Text, req.Id))
	if err != nil {
		h.logger.Errorln(err)
	}
}

func CheckToken(w http.ResponseWriter, token string, service *service.Service, logger *logger.Logger) bool {
	ok, err := service.CheckToken(context.TODO(), token)
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

func (h *Handler) CheckHandler(w http.ResponseWriter, r *http.Request) {
	var path = strings.Split(r.URL.Path, "/")
	if CheckToken(w, path[len(path)-1], h.service, h.logger) {

		h.logger.Printf(OkCheck, path[len(path)-1])

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
