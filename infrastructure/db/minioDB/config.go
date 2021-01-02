package miniodb

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

func init() {
	// optionally look for config in the working directory
	viper.AddConfigPath("./env/")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.ReadInConfig()
}

func InitConfig() (*Config, error) {
	config := &Config{
		Endpoint:        viper.GetString("Minio.Endpoint"),
		AccessKeyID:     viper.GetString("Minio.AccessKeyID"),
		SecretAccessKey: viper.GetString("Minio.SecretAccessKey"),
		UseSSL:          viper.GetBool("Minio.UseSSL"),
		BucketName:      viper.GetString("Minio.BucketName"),
	}
	if config.Endpoint == "" {
		return nil, fmt.Errorf("Minio.Endpoint is not found in config")
	}
	if config.AccessKeyID == "" {
		return nil, fmt.Errorf("Minio.AccessKeyID is not found in config")
	}
	if config.SecretAccessKey == "" {
		return nil, fmt.Errorf("Minio.SecretAccessKey is not found in config")
	}
	if config.BucketName == "" {
		return nil, fmt.Errorf("Minio.BucketName is not found in config")
	}
	return config, nil
}
