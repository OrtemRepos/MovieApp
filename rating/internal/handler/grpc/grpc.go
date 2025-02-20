package grpc

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/rating/internal/controller/rating"
	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

type Handler struct {
	gen.UnimplementedRatingServiceServer
	ctrl   *rating.Controleer
	logger *zap.Logger
}

func New(ctrl *rating.Controleer, logger *zap.Logger) *Handler {
	return &Handler{ctrl: ctrl, logger: logger}
}

func (h *Handler) GetAggregateRating(
	ctx context.Context,
	req *gen.GetAggregateRatingRequest,
) (*gen.GetAggregateRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty id")
	}
	v, err := h.ctrl.GetAggregatedRating(
		ctx, model.RecordID(req.RecordId),
		model.RecordType(req.RecordType),
	)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, repository.ErrWrongType) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.GetAggregateRatingResponse{RecordValue: v}, nil
}

func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.UserId == "" || req.RecordId == "" {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty user id or record id")
	}
	if err := h.ctrl.PutRating(
		ctx, model.RecordID(req.RecordId),
		model.UserID(req.UserId), model.RecordType(req.RecordType),
		model.RatingValue(req.RatingValue)); err != nil {

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.PutRatingResponse{RecordId: req.RecordId, RatingValue: req.RatingValue}, nil
}
