package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
	"movieexample.com/metadata/internal/controller/metadata"
	"movieexample.com/metadata/internal/repository"
)

type Handler struct {
	ctrl   *metadata.Controller
	logger *zap.Logger
}

func New(ctrl *metadata.Controller, logger *zap.Logger) *Handler {
	return &Handler{ctrl, logger}
}

func (h *Handler) GetMetadata(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("GET REQUEST",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.String("remote_addr", req.RemoteAddr),
	)
	id := req.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	m, err := h.ctrl.Get(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("repository get error", zap.Error(err))
		return
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		h.logger.Error("response encode error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
