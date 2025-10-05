package model

import (
	"fmt"
	"os"
	"path/filepath"

	value_utils "github.com/Beriholic/geminic/internal/utils"
	"github.com/spf13/viper"
)

const configFilePath = "$HOME/.config/geminic/config.toml"

var v *viper.Viper

type Config struct {
	Key           string `mapstructure:"key"`
	Model         string `mapstructure:"model"`
	Emoji         bool   `mapstructure:"emoji"`
	CustomURL     string `mapstructure:"custom_url"`
	I18n          string `mapstructure:"i18n"`
	ModelProvider string `mapstructure:"model_provider"`
}
type LocalConfig struct {
	Emoji bool   `mapstructure:"emoji"`
	I18n  string `mapstructure:"i18n"`
}

func (c *Config) UseCustom() bool {
	return c.CustomURL != ""
}

func (c *Config) Load() error {
	v, err := initViper()
	if err != nil {
		return err
	}

	c.Key = v.GetString("key")
	c.Model = v.GetString("model")
	c.Emoji = v.GetBool("emoji")
	c.CustomURL = v.GetString("custom_url")
	c.I18n = value_utils.GetStrngOrDefault(v.GetString("i18n"), "en_US")
	c.ModelProvider = v.GetString("model_provider")
	return nil
}

func (c *Config) Save() error {
	expandedPath := os.ExpandEnv(configFilePath)
	v, err := initViper()
	if err != nil {
		return err
	}

	v.Set("key", c.Key)
	v.Set("model", c.Model)
	v.Set("emoji", c.Emoji)
	v.Set("custom_url", c.CustomURL)
	v.Set("i18n", c.I18n)
	v.Set("model_provider", c.ModelProvider)

	if err := v.WriteConfigAs(expandedPath); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	fmt.Printf("Configuration saved to %s\n", expandedPath)
	return nil
}

func initViper() (*viper.Viper, error) {
	if v != nil {
		return v, nil
	}

	v := viper.New()
	expandedPath := os.ExpandEnv(configFilePath)

	configDir := filepath.Dir(expandedPath)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	v.SetConfigFile(expandedPath)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to red config file: %v", err)
	}
	return v, nil
}
