package rating

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

type ratingRepository interface {
	Get(
		ctx context.Context, recordID model.RecordID,
		recordType model.RecordType,
	) ([]model.Rating, error)

	Put(
		ctx context.Context, recordID model.RecordID,
		recordType model.RecordType, record *model.Rating,
	) error
}

type Controleer struct {
	repo   ratingRepository
	logger *zap.Logger
}

func New(repo ratingRepository, logger *zap.Logger) *Controleer {
	return &Controleer{repo, logger}
}

func (c *Controleer) GetAggregatedRating(
	ctx context.Context, recordID model.RecordID,
	recordType model.RecordType,
) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordID, recordType)
	if errors.Is(err, repository.ErrNotFound) || errors.Is(err, repository.ErrWrongType) {
		return 0, err
	} else if err != nil {
		c.logger.Error("unexpected error in the controller", zap.Error(err))
		return 0, err
	}
	var sum float64
	for _, r := range ratings {
		sum += float64(r.Value)
	}
	return sum / float64(len(ratings)), nil
}

func (c *Controleer) PutRating(
	ctx context.Context, recordID model.RecordID,
	recordType model.RecordType, rating *model.Rating,
) error {
	return c.repo.Put(ctx, recordID, recordType, rating)
}
