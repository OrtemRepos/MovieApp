package kafka

import (
	"context"
	"encoding/json"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
	"movieexample.com/rating/pkg/model"
)

type Ingester struct {
	consumer *kgo.Client
	logger   *zap.Logger
	topic    string
}

func NewIngester(addr string, groupID string, topic string, logger *zap.Logger) (*Ingester, error) {
	consumer, err := kgo.NewClient(
		kgo.SeedBrokers(addr),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
	)
	if err != nil {
		return nil, err
	}
	return &Ingester{consumer, logger, topic}, nil
}

func (i *Ingester) Ingest(ctx context.Context) (chan model.RatingEvent, error) {
	ch := make(chan model.RatingEvent, 1)
	go func() {
		defer close(ch)
		defer i.consumer.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				fetches := i.consumer.PollFetches(ctx)
				if fetches.IsClientClosed() {
					i.logger.Info("Kafka client closed")
					return
				}
				callBackErr := func(topic string, partition int32, err error) {
					i.logger.Error(
						"Error consuming",
						zap.String("topic", topic),
						zap.Int32("partition", partition),
						zap.Error(err),
					)
				}
				callBackRecord := func(record *kgo.Record) {
					var event model.RatingEvent
					if err := json.Unmarshal(record.Value, &event); err != nil {
						i.logger.Error("Failed to unmarshal event", zap.Any("event", event), zap.Error(err))
						return
					}
					ch <- event
				}
				fetches.EachError(callBackErr)
				fetches.EachRecord(callBackRecord)

			}
		}
	}()
	return ch, nil
}
