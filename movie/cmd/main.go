package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"movieexample.com/gen"
	"movieexample.com/movie/internal/controller/movie"
	metadatagateway "movieexample.com/movie/internal/gateway/metadata/grpc"
	ratinggateway "movieexample.com/movie/internal/gateway/rating/grpc"
	grpchandler "movieexample.com/movie/internal/handler/grpc"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {
	port := flag.Int("p", 0, "API handler port")
	flag.Parse()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	cfg := loadConfig()
	if *port == 0 {
		port = &cfg.APIConfig.Port
	}

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

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)

	ctrl := movie.New(ratingGateway, metadataGateway, logger)
	h := grpchandler.New(ctrl, logger)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Int("port", *port), zap.Error(err))
	}
	srv := grpc.NewServer()
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		logger.Fatal("Failed to accepts incoming connections", zap.Error(err))
	}
}
