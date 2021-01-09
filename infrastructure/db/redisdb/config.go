package redisdb

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Address string
}

func InitConfig() (*Config, error) {
	config := &Config{
		Address: viper.GetString("Redis.Address"),
	}
	if config.Address == "" {
		return nil, fmt.Errorf("Redis.Address is not found in config")
	}
	return config, nil
}
