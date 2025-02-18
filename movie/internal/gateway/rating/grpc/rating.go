package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/internal/grpcutil"
	"movieexample.com/movie/internal/gateway"
	"movieexample.com/pkg/discovery"
	"movieexample.com/rating/pkg/model"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) GetAggregatedRating(
	ctx context.Context, recordID model.RecordID, recordType model.RecordType,
) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer func() { _ = conn.Close() }()
	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregateRating(
		ctx,
		&gen.GetAggregateRatingRequest{
			RecordId:   string(recordID),
			RecordType: string(recordType),
		},
	)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return 0, gateway.ErrNotFound
		}
		return 0, err
	}
	return resp.RecordValue, nil
}

func (g *Gateway) PutRating(
	ctx context.Context, recordID model.RecordID,
	recordType model.RecordType, userID model.UserID, value model.RatingValue,
) error {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	client := gen.NewRatingServiceClient(conn)
	_, err = client.PutRating(ctx,
		&gen.PutRatingRequest{
			RecordId:    string(recordID),
			RecordType:  string(recordType),
			UserId:      string(userID),
			RatingValue: int32(value),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
