package handler

import (
	"fmt"
	"net/http"

	"github.com/vandi37/TgLogger/internal/web/api"
	"github.com/vandi37/vanerrors"
)

const (
	WrongMethod = "wrong method"
)

func (h *Handler) CheckMethod(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			err := api.SendError(w, http.StatusMethodNotAllowed, vanerrors.NewSimple(WrongMethod, fmt.Sprintf("method %s is not allowed, allowed method: %s", r.Method, method)))
			if err != nil {
				h.logger.Errorln(err)
			}
			return
		}
		next(w, r)
	}
}

func ContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next(w, r)
	}
}
