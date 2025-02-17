package http

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"go.uber.org/zap"
	"movieexample.com/metadata/pkg/model"
	"movieexample.com/movie/internal/gateway"
	"movieexample.com/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
	logger   *zap.Logger
}

func New(registry discovery.Registry, logger *zap.Logger) *Gateway {
	return &Gateway{registry, logger}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + "metadata"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		g.logger.Error("Error when requesting the metadata service", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	g.logger.Info("Calling metadata service.", zap.Any("request", req))
	query := req.URL.Query()
	query.Add("id", id)
	req.URL.RawQuery = query.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}
	var v *model.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
