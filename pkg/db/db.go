package db

import (
	"context"
	"fmt"
	"go-wal/config"
	"go-wal/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewStorage(cfg *config.Database) (*gorm.DB, error) {
	// Init connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
	)
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   fmt.Sprintf("%s.", cfg.Schema),
			SingularTable: true,
		},
		CreateBatchSize: 7000,
	})
	if err != nil {
		return nil, err
	}

	db, err := client.DB()
	if err != nil {
		logger.Fatal(context.Background()).Err(err).Msg("Cannot connect to PostgreSQL")
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIdleConnection)
	db.SetMaxOpenConns(cfg.MaxActiveConnection)
	db.SetConnMaxIdleTime(cfg.MaxConnectionTimeout)
	if cfg.DebugLog {
		client = client.Debug()
	}

	return client, nil
}
