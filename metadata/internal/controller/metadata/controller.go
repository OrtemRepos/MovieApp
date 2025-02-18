package metadata

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"movieexample.com/metadata/internal/repository"
	"movieexample.com/metadata/pkg/model"
)

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
	Put(ctx context.Context, id string, metadata *model.Metadata) error
}

type Controller struct {
	repo   metadataRepository
	logger *zap.Logger
}

func New(repo metadataRepository, logger *zap.Logger) *Controller {
	return &Controller{repo, logger}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, err
	} else if err != nil {
		c.logger.Error("error when retrieving from the repository", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (c *Controller) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	if err := c.repo.Put(ctx, id, metadata); err != nil {
		return err
	}
	return nil
}
