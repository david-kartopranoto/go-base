package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB            DBConfig            `mapstructure:"db"`
	MessageBroker MessageBrokerConfig `mapstructure:"message-broker"`
	Limiter       LimiterConfig       `mapstructure:"limiter"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`
}

type MessageBrokerConfig struct {
	Host     string                 `mapstructure:"host"`
	Port     string                 `mapstructure:"port"`
	User     string                 `mapstructure:"user"`
	Password string                 `mapstructure:"password"`
	Queue    map[string]QueueConfig `mapstructure:"queue"`
}

type QueueConfig struct {
	Durable      bool `mapstructure:"durable"`
	DeleteUnused bool `mapstructure:"delete-unused"`
	Exclusive    bool `mapstructure:"exclusive"`
	NoWait       bool `mapstructure:"no-wait"`
}

type LimiterConfig struct {
	MaxEventsPerSec int `mapstructure:"maxEventsPerSec"`
	MaxBurstSize    int `mapstructure:"maxBurstSize"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string, name string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
