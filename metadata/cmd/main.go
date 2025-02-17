package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"movieexample.com/metadata/internal/controller/metadata"
	httphandler "movieexample.com/metadata/internal/handler/http"
	"movieexample.com/metadata/internal/repository/memory"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
)

const serviceName = "metadata"

func main() {
	port := flag.Int("p", 8081, "API handler port")
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func(){ _ = logger.Sync() }()

	logger.Info("Starting the metadata service", zap.Int("port", *port))

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

	repo := memory.New()
	ctrl := metadata.New(repo, logger)

	h := httphandler.New(ctrl, logger)

	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic(err)
	}
}
