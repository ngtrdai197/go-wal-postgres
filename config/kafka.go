package config

type Kafka struct {
	Brokers                    []string `mapstructure:"broker_list"`
	Partition                  int32    `mapstructure:"partition"`
	Partitioner                string   `mapstructure:"partitioner"` // "The partitioning scheme to use. Can be `hash`, `manual`, or `random`")
	WalDatabaseConsumerGroupId string   `mapstructure:"wal_database_group_id"`
}
