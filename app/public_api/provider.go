package publicapi

import (
	"go-wal/config"
	"go-wal/pkg/db"
	"go.uber.org/dig"
	"log"
)

var container *dig.Container

func InitProvider() {
	c := dig.New()
	provideInfrastructure(c)

	container = c
}

func provideInfrastructure(c *dig.Container) {
	// provide database config
	if err := c.Provide(func() *config.Database {
		return &config.Config.Database
	}); err != nil {
		log.Fatalln("Cannot provide database config", err)
	}
	if err := c.Provide(db.NewStorage); err != nil {
		log.Fatalln("Cannot provide database storage", err)
	}
}

func GetContainer() *dig.Container {
	return container
}
