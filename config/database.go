package config

import "time"

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"database"`
	Schema   string `mapstructure:"schema"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`

	MaxIdleConnection    int           `mapstructure:"max_idle_conns"`
	MaxActiveConnection  int           `mapstructure:"max_active_conns"`
	MaxConnectionTimeout time.Duration `mapstructure:"max_conn_timeout"`

	DebugLog bool `mapstructure:"debug_log"`
}
