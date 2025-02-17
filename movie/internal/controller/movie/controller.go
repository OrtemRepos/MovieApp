package movie

import (
	"context"
	"errors"

	"go.uber.org/zap"
	metadataModel "movieexample.com/metadata/pkg/model"
	"movieexample.com/movie/internal/gateway"
	"movieexample.com/movie/pkg/model"
	ratingModel "movieexample.com/rating/pkg/model"
)

type ratingGateway interface {
	GetAggregatedRating(
		ctx context.Context, recordID ratingModel.RecordID,
		recordType ratingModel.RecordType,
	) (float64, error)

	PutRating(
		ctx context.Context, recordID ratingModel.RecordID,
		recordType ratingModel.RecordType, userID ratingModel.UserID,
		value ratingModel.RatingValue,
	) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadataModel.Metadata, error)
}

type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
	logger          *zap.Logger
}

func New(
	ratingGateway ratingGateway,
	memetadataGateway metadataGateway, 
	logger *zap.Logger,
) *Controller {
	return &Controller{ratingGateway, memetadataGateway, logger}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if errors.Is(err, gateway.ErrNotFound) {
		return nil, err
	} else if err != nil {
		c.logger.Error("unexpected error in metadataGateway", zap.Error(err))
		return nil, err
	}
	details := &model.MovieDetails{Metadata: *metadata}
	rating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingModel.RecordID(id), ratingModel.RecordTypeMovie)
	if errors.Is(err, gateway.ErrNotFound) {
		return nil, err
	} else if err != nil {
		c.logger.Error("unexpected error in ratingGateway", zap.Error(err))
		return nil, err
	}
	details.Rating = rating
	return details, nil
}
