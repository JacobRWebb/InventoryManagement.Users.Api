package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Service ServiceConfig
	Consul  ConsulConfig
	DB      DBConfig
}

type ServiceConfig struct {
	Name     string `validate:"required"`
	GRPCAddr string `validate:"required,ip"`
	GRPCPort int    `validate:"required,min=1,max=65535"`
}

type ConsulConfig struct {
	Addr           string `validate:"required,url"`
	DeregisterTime string `validate:"required"`
	IntervalTime   string `validate:"required"`
}

type DBConfig struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required,min=1,max=65535"`
	Name     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DSN      string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found. Using environment variables.")
	} else {
		fmt.Println("Loaded .env file successfully")
	}

	config := &Config{
		Service: ServiceConfig{
			Name:     getEnv("SERVICE_NAME", "DEFAULT_SERVICE"),
			GRPCAddr: getEnv("GRPC_ADDR", "127.0.0.1"),
			GRPCPort: getEnvAsInt("GRPC_PORT", 50051),
		},
		Consul: ConsulConfig{
			Addr:           getEnv("CONSUL_ADDR", "http://localhost:8500"),
			DeregisterTime: getEnv("CONSUL_DEREGISTER_TIME", "10m"),
			IntervalTime:   getEnv("CONSUL_INTERVAL_TIME", "2m"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Name:     getEnv("DB_NAME", "DBName"),
			User:     getEnv("DB_USER", "DBUser"),
			Password: getEnv("DB_PASSWORD", "DBPass"),
			DSN:      getEnv("DATABASE_DSN", ""),
		},
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
