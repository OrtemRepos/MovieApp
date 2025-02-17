package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
	"movieexample.com/movie/internal/controller/movie"
	"movieexample.com/movie/internal/gateway"
)

type Handler struct {
	ctrl   *movie.Controller
	logger *zap.Logger
}

func New(ctrl *movie.Controller, logger *zap.Logger) *Handler {
	return &Handler{ctrl, logger}
}

func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("GET REQUEST",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.String("remote_addr", req.RemoteAddr),
	)
	id := req.FormValue("id")
	details, err := h.ctrl.Get(req.Context(), id)
	if errors.Is(err, gateway.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Error("unexpected error in Movie Handler", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&details); err != nil {
		h.logger.Error("response encode error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
} 