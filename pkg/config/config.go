package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName          string `mapstructure:"SERVICE_NAME"`
	GRPCAddr             string `mapstructure:"GRPC_ADDR"`
	GRPCPort             int    `mapstructure:"GRPC_PORT"`
	ConsulAddr           string `mapstructure:"CONSUL_ADDR"`
	ConsulDeregisterTime string `mapstructure:"CONSUL_DEREGISTER_TIME"`
	ConsulIntervalTime   string `mapstructure:"CONSUL_INTERVAL_TIME"`
	DBHost               string `mapstructure:"DB_HOST"`
	DBPort               int    `mapstructure:"DB_PORT"`
	DBName               string `mapstructure:"DB_NAME"`
	DBUser               string `mapstructure:"DB_USER"`
	DBPassword           string `mapstructure:"DB_PASSWORD"`
	DatabaseDSN          string `mapstructure:"DATABASE_DSN"`
}

func NewConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("..")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf(".env file not found: %w", err)
		} else {
			return nil, fmt.Errorf("error reading .env file: %w", err)
		}
	}

	v.AutomaticEnv()

	setDefaults(v)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func setDefaults(v *viper.Viper) {
	defaults := map[string]interface{}{
		"SERVICE_NAME":           "DEFAULT_SERVICE",
		"GRPC_ADDR":              "127.0.0.1",
		"GRPC_PORT":              50051,
		"CONSUL_ADDR":            "http://localhost:8500",
		"CONSUL_DEREGISTER_TIME": "10m",
		"CONSUL_INTERVAL_TIME":   "2m",
		"DB_HOST":                "localhost",
		"DB_PORT":                5432,
		"DB_NAME":                "DBName",
		"DB_USER":                "DBUser",
		"DB_PASSWORD":            "DBPass",
		"DATABASE_DSN":           "",
	}

	for key, value := range defaults {
		v.SetDefault(key, value)
	}
}
