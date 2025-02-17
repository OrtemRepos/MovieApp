package http

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"go.uber.org/zap"
	"movieexample.com/movie/internal/gateway"
	"movieexample.com/pkg/discovery"
	"movieexample.com/rating/pkg/model"
)

type Gateway struct {
	registry discovery.Registry
	logger   *zap.Logger
}

func New(registry discovery.Registry, logger *zap.Logger) *Gateway {
	return &Gateway{registry, logger}
}

func (g *Gateway) GetAggregatedRating(
	ctx context.Context, recordID model.RecordID,
	recordType model.RecordType,
) (float64, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return 0, nil
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + "rating"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		g.logger.Error("Error when GET requesting the rating service", zap.String("url", url), zap.Error(err))
		return 0, err
	}
	g.logger.Info("Calling rating service.", zap.Any("request", req))
	query := req.URL.Query()
	query.Add("id", string(recordID))
	query.Add("type", string(recordType))
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", resp)
	}
	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (g *Gateway) PutRating(
	ctx context.Context, recordID model.RecordID,
	recordType model.RecordType, userID model.UserID,
	value model.RatingValue,
) error {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return err
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		g.logger.Error("Error when PUT requesting the rating service", zap.String("url", url), zap.Error(err))
		return err
	}
	g.logger.Info("Calling rating service")
	query := req.URL.Query()
	query.Add("id", string(recordID))
	query.Add("type", string(recordType))
	query.Add("userId", string(userID))
	query.Add("value", fmt.Sprintf("%v", value))
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}
	return nil
}
