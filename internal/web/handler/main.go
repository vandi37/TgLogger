package handler

import (
	"net/http"

	"github.com/vandi37/TgLogger/pkg/db"
	"github.com/vandi37/TgLogger/pkg/logger"
)

// The handler
type Handler struct {
	logger *logger.Logger
	db     *db.DB
	funcs  map[string]http.HandlerFunc
}

// Created a new handler
func NewHandler(db *db.DB, logger *logger.Logger) *Handler {
	// Creating handler
	handler := Handler{
		logger: logger,
		db:     db,
	}

	// Adding functions
	handler.funcs = map[string]http.HandlerFunc{}

	return &handler
}

// Serve
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// The header
	w.Header().Add("Content-Type", "application/json")

	// Gets the handler func
	fn, ok := h.funcs[r.URL.Path]
	if !ok {

		// // Not found
		// err := api.SendErrorResponse(w, http.StatusNotFound, vanerrors.NewSimple(NotFound))
		// if err != nil {
		// 	h.logger.Errorln(err)
		// 	return
		// }

		return
	}

	// Runs the handler
	fn(w, r)
}
