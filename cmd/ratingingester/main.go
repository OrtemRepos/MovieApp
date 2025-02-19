package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
	"movieexample.com/rating/pkg/model"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	logger.Info("Creating Kafka Client")

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(),
	)
	if err != nil {
		panic(err)
	}
	defer cl.Close()
	ctx := context.Background()

	const fileName = "ratingsdata.json"
	logger.Info("Read rating event from file", zap.String("file", fileName))

	ratingEvents, err := readRatingEvents(fileName)
	if err != nil {
		panic(err)
	}

	const topic = "ratings"
	if err := produceRatingEvents(ctx, topic, cl, ratingEvents); err != nil {
		panic(err)
	}

	const timeout = 10 * time.Second
	logger.Info("Waiting until all events get produced", zap.Duration("timeout", timeout))

	if err := cl.Flush(ctx); err != nil {
		panic(err)
	}
}

func readRatingEvents(fileName string) ([]model.RatingEvent, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var ratings []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func produceRatingEvents(ctx context.Context, topic string, producer *kgo.Client, events []model.RatingEvent) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(events))
	for _, ratingEvent := range events {
		encodedEvent, err := json.Marshal(ratingEvent)
		if err != nil {
			return err
		}
		record := kgo.Record{
			Topic: topic,
			Value: encodedEvent,
		}
		wg.Add(1)
		producer.Produce(
			ctx, &record,
			func(r *kgo.Record, err error) {
				defer wg.Done()
				if err != nil {
					log.Printf("record had a produce error: %v\n", err)
					errChan <- err
				}
			},
		)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
