package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ConnectionString      string        `mapstructure:"connection_string"`
	MaxConnIdleTime       string        `mapstructure:"max_idle_connections"`
	MaxOpenConnections    string        `mapstructure:"max_open_connections"`
	ConnectionMaxLifetime string        `mapstructure:"connection_max_lifetime"`
	ApiSecret             string        `mapstructure:"token_hour_lifespan"`
	AccessTokenDuration   time.Duration `mapstructure:"api_secret"`
	HttpServerAddress     string        `mapstructure:"server_address"`
}

func LoadConfig(name string, path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	//viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("config: %v", err)
		return
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("config: %v", err)
		return
	}
	return
}
