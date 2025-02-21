package grpc

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/metadata/internal/controller/metadata"
	"movieexample.com/metadata/internal/repository"
	"movieexample.com/metadata/pkg/model"
)

type Handler struct {
	gen.UnimplementedMetadataServiceServer
	ctrl   *metadata.Controller
	logger *zap.Logger
}

func New(ctrl *metadata.Controller, logger *zap.Logger) *Handler {
	return &Handler{ctrl: ctrl, logger: logger}
}

func (h *Handler) GetMetadata(
	ctx context.Context,
	req *gen.GetMetadataRequest,
) (*gen.GetMetadataResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty id")
	}
	m, err := h.ctrl.Get(ctx, req.MovieId)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.GetMetadataResponse{Metadata: model.MetadataToProto(m)}, nil
}

func (h *Handler) PutMetadata(ctx context.Context, req *gen.PutMetadataReuqest) (*gen.PutMetadataResponse, error) {
	if req == nil || req.Id == "" || req.Metadata == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty id or metadata")
	}
	if err := h.ctrl.Put(ctx, req.Id, model.MetadataFromProto(req.Metadata)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.PutMetadataResponse{Id: req.Id}, nil
}
