package handler

import (
	"fmt"
	"net/http"

	"github.com/vandi37/TgLogger/internal/web/api"
	"github.com/vandi37/TgLogger/pkg/logger"
	"github.com/vandi37/TgLogger/pkg/service"
	"github.com/vandi37/vanerrors"
)

// Errors
const (
	NotFound = "this page is unavailable"
)

// The handler
type Handler struct {
	logger   *logger.Logger
	service  *service.Service
	mainPath string
}

// Created a new handler
func NewHandler(service *service.Service, logger *logger.Logger) *Handler {
	// Creating handler
	handler := Handler{
		service:  service,
		mainPath: "/send",
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
			return
		}

		return
	}

	// Runs the handler
	
}
