package main

import (
	"net/http"

	"go.uber.org/zap"
	"movieexample.com/rating/internal/controller/rating"
	handlerhttp "movieexample.com/rating/internal/handler/http"
	"movieexample.com/rating/internal/repository/memory"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("Starting the rating service")

	repo := memory.New()
	ctrl := rating.New(repo, logger)
	handler := handlerhttp.New(ctrl, logger)
	http.Handle("/rating", http.HandlerFunc(handler.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
