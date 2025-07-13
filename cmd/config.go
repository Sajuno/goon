package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	AssistantID string `mapstructure:"assistant_id"`
	APIKey      string `mapstructure:"api_key"`
}

var cfg *config

func loadConfig() error {
	viper.SetConfigName("goon.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("GOON")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return fmt.Errorf("no configuration file found: %w", err)
		}
	}

	cfg = &config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return nil
}
