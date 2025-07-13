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

func loadConfig() (config, error) {
	viper.SetConfigName("goon.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("GOON")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	var c config
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return config{}, fmt.Errorf("no configuration file found: %w", err)
		}
	}

	if err := viper.Unmarshal(&c); err != nil {
		return config{}, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return c, nil
}
