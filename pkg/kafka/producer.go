package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"go-wal/config"
	"go-wal/constant"
	"go-wal/pkg/logger"
)

type Producer struct {
	syncProducer sarama.SyncProducer
}

func newProducer(cnf *config.Kafka) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true

	cfg.Producer.Partitioner = sarama.NewHashPartitioner

	client, err := sarama.NewClient(cnf.Brokers, cfg)
	if err != nil {
		return nil, err
	}
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

// NewProducer create new kafka producer
func NewProducer(cnf *config.Kafka) (*Producer, error) {
	syncProducer, err := newProducer(cnf)
	if err != nil {
		return nil, err
	}
	return &Producer{
		syncProducer: syncProducer,
	}, nil
}

// Publish2KafkaMessages publish messages to kafka
func (k *Producer) Publish2KafkaMessages(ctx context.Context, messages []*sarama.ProducerMessage) error {
	if len(messages) == 0 {
		return nil
	}
	if err := k.syncProducer.SendMessages(messages); err != nil {
		logger.Error(ctx).Err(err).Msg("Publish2KafkaMessages.Failed")
		return err
	}
	logger.Info(ctx).Msg("Publish2KafkaMessages.SendMessages")
	return nil
}

func (k *Producer) BuildMessage(ctx context.Context, topic, key string, msgData []byte) sarama.ProducerMessage {
	return sarama.ProducerMessage{
		Topic:     topic,
		Partition: config.Config.Kafka.Partition,
		Value:     sarama.StringEncoder(msgData),
		Key:       sarama.StringEncoder(key),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte(constant.XRequestId),
				Value: []byte(logger.GetTraceIDFromContext(ctx)),
			},
		},
	}
}
