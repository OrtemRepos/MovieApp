package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"
	"movieexample.com/rating/internal/controller/rating"
	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

type Handler struct {
	ctrl   *rating.Controleer
	logger *zap.Logger
}

func New(ctrl *rating.Controleer, logger *zap.Logger) *Handler {
	return &Handler{ctrl, logger}
}

func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("GET REQUEST",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.String("remote_addr", req.RemoteAddr),
	)
	recordID := model.RecordID(req.FormValue("id"))
	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recordType := model.RecordType(req.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodGet:
		v, err := h.ctrl.GetAggregatedRating(req.Context(), recordID, recordType)
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if errors.Is(err, repository.ErrWrongType) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("unexpected error in the controller", zap.Error(err))
			return
		}
		if err := json.NewEncoder(w).Encode(v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error("response encode error", zap.Error(err))
		}
	case http.MethodPut:
		userID := model.UserID(req.FormValue("userId"))
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		v, err := strconv.ParseFloat(req.FormValue("value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := h.ctrl.PutRating(req.Context(), recordID, userID, recordType, model.RatingValue(v)); err != nil {
			h.logger.Error("repository put error", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
