package main

import (
	"net/http"

	"go.uber.org/zap"
	"movieexample.com/movie/internal/controller/movie"
	metadatagateway "movieexample.com/movie/internal/gateway/metadata/http"
	moviegateway "movieexample.com/movie/internal/gateway/rating/http"
	httphandelr "movieexample.com/movie/internal/handler/http"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	
	logger.Info("Starting the movie service")

	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := moviegateway.New("localhos:8082")

	ctrl := movie.New(ratingGateway, metadataGateway, logger)

	handler := httphandelr.New(ctrl, logger)

	http.Handle("/movie", http.HandlerFunc(handler.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}