package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/internal/grpcutil"
	"movieexample.com/metadata/pkg/model"
	"movieexample.com/movie/internal/gateway"
	"movieexample.com/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, gateway.ErrNotFound
		}
		return nil, err
	}
	return model.MetadataFromProto(resp.Metadata), nil
}

func (g *Gateway) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	client := gen.NewMetadataServiceClient(conn)
	_, err = client.PutMetadata(
		ctx,
		&gen.PutMetadataReuqest{
			Id:       id,
			Metadata: model.MetadataToProto(metadata),
		},
	)
	return err
}
