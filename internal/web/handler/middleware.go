package handler

import (
	"fmt"
	"net/http"

	"github.com/vandi37/TgLogger/internal/web/api"
	"github.com/vandi37/TgLogger/pkg/logger"
	"github.com/vandi37/vanerrors"
)

const (
	WrongMethod = "wrong method"
)

func CheckMethod(method string, next http.HandlerFunc, logger *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			err := api.SendError(w, http.StatusMethodNotAllowed, vanerrors.NewSimple(WrongMethod, fmt.Sprintf("method %s is not allowed, allowed method: %s", r.Method, method)))
			if err != nil {
				logger.Errorln(err)
				return
			}
			return
		}
		next(w, r)
	}
}
