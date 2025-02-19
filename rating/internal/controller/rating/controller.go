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

type ratingIngester interface {
	Ingest(ctx context.Context) (chan model.RatingEvent, error)
}

type Controleer struct {
	repo     ratingRepository
	ingester ratingIngester
	logger   *zap.Logger
}

func New(repo ratingRepository, ingester ratingIngester, logger *zap.Logger) *Controleer {
	return &Controleer{repo, ingester, logger}
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
	ctx context.Context, recordID model.RecordID, userID model.UserID,
	recordType model.RecordType, value model.RatingValue,
) error {
	rating := &model.Rating{
		RecordID:   recordID,
		RecordType: recordType,
		UserID:     userID,
		Value:      value,
	}
	return c.repo.Put(ctx, recordID, recordType, rating)
}

func (c *Controleer) StartIngestion(ctx context.Context) error {
	ch, err := c.ingester.Ingest(ctx)
	if err != nil {
		return err
	}
	for e := range ch {
		if err := c.PutRating(ctx, e.RecordID, e.UserID, e.RecordType, e.Value); err != nil {
			return err
		}

	}
	return nil
}