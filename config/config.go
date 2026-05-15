package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const cfgPath = "./config"

type Config struct {
	Logger   `mapstructure:"logger"`
	Server   `mapstructure:"server"`
	Database `mapstructure:"database"`
}

type Server struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"ssl_mode"`

	PoolMaxConns              int           `mapstructure:"pool_max_conns"`
	PoolMinConns              int           `mapstructure:"pool_min_conns"`
	PoolMaxConnLifetime       time.Duration `mapstructure:"pool_max_conn_lifetime"`
	PoolMaxConnIdleTime       time.Duration `mapstructure:"pool_max_conn_idle_time"`
	PoolHealthCheckPeriod     time.Duration `mapstructure:"pool_health_check_period"`
	PoolMaxConnLifetimeJitter time.Duration `mapstructure:"pool_max_conn_lifetime_jitter"`
}

type Logger struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.AddConfigPath(cfgPath)
	v.SetConfigType("yaml")
	v.SetConfigName("config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %w", err)
	}

	return &cfg, nil
}
