package wal

import (
	"context"
	"go-wal/config"
	"go-wal/domain/consumer/wal"
	"go-wal/pkg/kafka"
)

func init() {
	config.Init()
	Init()
}

func ConsumerLaunch() {
	p := wal.NewProcessor(wal.NewHandler())
	cfg := config.Config.Kafka
	kafka.NewConsumer(context.Background(), p, &cfg, wal.WalChangeTopic, cfg.WalDatabaseConsumerGroupId)
}
