package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
)

const configFilePath = "$HOME/.config/geminic/config.toml"

type GeminicConfig struct {
	Key   string
	Model string
	Emoji bool
}

func Verify() error {
	cfg := Get()
	if cfg.Key == "" {
		return fmt.Errorf("api key must be set, use `geminic config` to set it")
	}
	if cfg.Model == "" {
		return fmt.Errorf("model must be set, use `geminic config` to set it")
	}
	return nil
}

var (
	geminicConfigOnce sync.Once
	geminicConfig     *GeminicConfig = nil
)

func Get() *GeminicConfig {
	geminicConfigOnce.Do(func() {
		geminicConfig = load()
	})

	return geminicConfig
}

func Create() error {
	expandedPath := os.ExpandEnv(configFilePath)

	configDir := filepath.Dir(expandedPath)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	viper.SetConfigFile(expandedPath)
	_ = viper.ReadInConfig()

	keyState := viper.GetString("key")
	modelState := viper.GetString("model")
	emojiState := viper.GetBool("emoji")

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is your Gemini API key?").
				Value(&keyState),
			huh.NewInput().
				Title("Which model do you want to use?").
				Value(&modelState),
			huh.NewSelect[bool]().
				Title("Do you want to enable emoji?").
				Options(
					huh.NewOption("Yes", true),
					huh.NewOption("No", false),
				).
				Value(&emojiState),
		).WithTheme(huh.ThemeBase()),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	viper.Set("key", keyState)
	viper.Set("model", modelState)
	viper.Set("emoji", emojiState)

	if err := viper.WriteConfigAs(expandedPath); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	fmt.Printf("Configuration saved to %s\n", expandedPath)
	return nil
}

func load() *GeminicConfig {
	expandedPath := os.ExpandEnv(configFilePath)

	viper.SetConfigFile(expandedPath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		return nil
	}

	return &GeminicConfig{
		Key:   viper.GetString("key"),
		Model: viper.GetString("model"),
		Emoji: viper.GetBool("emoji"),
	}
}

func SetModel(model string) error {
	expandedPath := os.ExpandEnv(configFilePath)

	viper.SetConfigFile(expandedPath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		return nil
	}

	viper.Set("model", model)
	return nil
}
