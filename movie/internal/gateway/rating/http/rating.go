package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"movieexample.com/movie/internal/gateway"
	"movieexample.com/rating/pkg/model"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr}
}

func (g *Gateway) GetAggregatedRating(
	ctx context.Context, recordID model.RecordID,
	recordType model.RecordType,
) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.addr+"/rating", nil)
	if err != nil {
		return 0, err
	}
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, g.addr+"/rating", nil)
	if err != nil {
		return err
	}
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
