package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env    string
	Server struct {
		Host string
		Port int
	}
	DB  DBConfig
	App AppConfig
	WS  WSConfig
}

type AppConfig struct {
	MessageMaxLen int
	TTL           time.Duration
	SleepTime     time.Duration
	LogLVL        string
}

type DBConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     int
	SSLMode  string
}

type WSConfig struct {
	ReadDeadline  time.Duration
	WriteDeadline time.Duration
	PingPeriod    time.Duration
	BuffSize      int64
}

func SetConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	configName := os.Getenv("CONFIG_NAME")
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return &cfg
}
