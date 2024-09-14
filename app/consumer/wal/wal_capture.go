package wal

import (
	"context"
	"go-wal/config"
	"go-wal/pkg/kafka"
	"go-wal/pkg/logger"
	"go-wal/pkg/wal"
	"log"
)

func init() {
	config.Init()
	Init()
}

func CaptureListen() {
	listener := wal.NewPgWalListener(newProducer())

	ctx := context.Background()
	if err := listener.Start(ctx); err != nil {
		logger.Error(ctx).Err(err).Msg("Failed to start wal listener")
	}
}

func newProducer() *kafka.Producer {
	var producer *kafka.Producer
	if err := GetContainer().Invoke(func(p *kafka.Producer) {
		producer = p
	}); err != nil {
		log.Fatalln("Cannot get kafka producer", err)
	}
	return producer
}
