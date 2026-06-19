package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func ReadConfig() Settings {
	viper.AddConfigPath("internal/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	var cfg Settings

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("Unable to decode into struct, %w \n", err))
	}

	cfg.Db.Host = os.Getenv("DB_HOST")
	cfg.Db.Port = os.Getenv("DB_PORT")
	cfg.Db.User = os.Getenv("DB_USER")
	cfg.Db.Password = os.Getenv("DB_PASSWORD")

	if cfg.Db.Host == "" {
		panic("DB_HOST is required")
	}

	if cfg.Db.Port == "" {
		panic("DB_PORT is required")
	}

	if cfg.Db.User == "" {
		panic("DB_USER is required")
	}

	if cfg.Db.Password == "" {
		panic("DB_PASSWORD is required")
	}

	return cfg
}
