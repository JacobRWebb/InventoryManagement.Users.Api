package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GRPCAddr             string `mapstructure:"GRPC_ADDR"`
	GRPCPort             int    `mapstructure:"GRPC_PORT"`
	ConsulAddr           string `mapstructure:"CONSUL_ADDR"`
	ConsulDeregisterTime string `mapstructure:"CONSUL_DEREGISTER_TIME"`
	ConsulIntervalTime   string `mapstructure:"CONSUL_INTERVAL_TIME"`
	ServiceName          string `mapstructure:"SERVICE_NAME"`
}

func NewConfig() (config *Config, err error) {
	viper.AutomaticEnv()

	viper.SetDefault("GRPC_ADDR", "127.0.0.1")
	viper.SetDefault("GRPC_PORT", 50051)
	viper.SetDefault("CONSUL_ADDR", "http://localhost:8500")
	viper.SetDefault("CONSUL_DEREGISTER_TIME", "1m")
	viper.SetDefault("CONSUL_INTERVAL_TIME", "1m")
	viper.SetDefault("SERVICE_NAME", "User_Service")

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
