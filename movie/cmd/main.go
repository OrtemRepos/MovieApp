package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"movieexample.com/movie/internal/controller/movie"
	metadatagateway "movieexample.com/movie/internal/gateway/metadata/http"
	moviegateway "movieexample.com/movie/internal/gateway/rating/http"
	httphandelr "movieexample.com/movie/internal/handler/http"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {
	port := flag.Int("p", 8083, "API handler port")
	flag.Parse()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func(){ _ = logger.Sync() }()

	logger.Info("Starting the movie service", zap.Int("port", *port))

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", *port)); err != nil {
		panic(err)
	}
	go func(ctx context.Context) {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				logger.Error("Failed to report healthy state", zap.Error(err))
			}
			time.Sleep(1 * time.Second)
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}(ctx)

	defer func() {
		if err := registry.Deregister(ctx, instanceID, serviceName); err != nil {
			logger.Error("Error when executing Deregister", zap.Error(err))
		}
	}()
	defer cancelFunc()

	metadataGateway := metadatagateway.New(registry, logger)
	ratingGateway := moviegateway.New(registry, logger)

	ctrl := movie.New(ratingGateway, metadataGateway, logger)

	handler := httphandelr.New(ctrl, logger)

	http.Handle("/movie", http.HandlerFunc(handler.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic(err)
	}
}
