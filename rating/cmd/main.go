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
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
	"movieexample.com/rating/internal/controller/rating"
	grpchandler "movieexample.com/rating/internal/handler/grpc"
	"movieexample.com/rating/internal/ingester/kafka"
	"movieexample.com/rating/internal/repository/postgresql"
)

const serviceName = "rating"

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	port := flag.Int("p", 0, "API handler port")
	flag.Parse()

	cfg := loadConfig()
	if *port == 0 {
		port = &cfg.APIConfig.Port
	}

	logger.Info("Starting the rating service", zap.Int("port", *port))

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
				logger.Info("Failed to report healthy state", zap.Error(err))
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

	repo, err := postgresql.New()
	if err != nil {
		panic(err)
	}
	ingestrer, err := kafka.NewIngester("localhost:9092", "ingester", "rating", logger)
	if err != nil {
		panic(err)
	}
	ctrl := rating.New(repo, ingestrer, logger)

	h := grpchandler.New(ctrl, logger)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Int("port", *port), zap.Error(err))
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		logger.Fatal("Failed to accepts incoming connections", zap.Error(err))
	}
}
