package wal

import (
	"go-wal/config"
	"go-wal/pkg/kafka"
	"go.uber.org/dig"
	"log"
)

var container *dig.Container

func Init() {
	c := dig.New()
	// Provide dependencies. kafka, es
	provideKafka(c)
	container = c
}

func provideKafka(c *dig.Container) {
	// Provide kafka
	if err := c.Provide(func() *config.Kafka {
		return &config.Config.Kafka
	}); err != nil {
		log.Fatalln("Cannot provide kafka config", err)
	}

	if err := c.Provide(kafka.NewProducer); err != nil {
		log.Fatalln("Cannot provide kafka producer", err)
	}
}

func GetContainer() *dig.Container {
	return container
}
