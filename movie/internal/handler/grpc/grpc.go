package grpc

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/movie/internal/controller/movie"
	"movieexample.com/movie/internal/gateway"
	"movieexample.com/movie/pkg/model"
)

type Handler struct {
	gen.UnimplementedMovieServiceServer
	ctrl   *movie.Controller
	logger *zap.Logger
}

func New(ctrl *movie.Controller, logger *zap.Logger) *Handler {
	return &Handler{ctrl: ctrl, logger: logger}
}

func (h *Handler) GetMovieDetails(
	ctx context.Context,
	req *gen.GetMovieDetailsRequest,
) (*gen.GetMovieDetailsResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty id")
	}
	m, err := h.ctrl.Get(ctx, req.MovieId)
	if errors.Is(err, gateway.ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.GetMovieDetailsResponse{MovieDetails: model.MovieDetailsToProto(m)}, nil
}
