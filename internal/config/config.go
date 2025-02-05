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
	Cot   bool
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
		return fmt.Errorf("Failed to create config directory: %v\n", err)
	}

	viper.SetConfigFile(expandedPath)
	_ = viper.ReadInConfig()

	keyState := viper.GetString("key")
	modelState := viper.GetString("model")
	emojiState := viper.GetBool("emoji")
	cotState := viper.GetBool("cot")

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(fmt.Sprintf("What is your Gemini API key?")).
				Value(&keyState),
			huh.NewInput().
				Title(fmt.Sprintf("Which model do you want to use?")).
				Value(&modelState),
			huh.NewSelect[bool]().
				Title(fmt.Sprintf("Do you want to enable emoji?")).
				Options(
					huh.NewOption("Yes", true),
					huh.NewOption("No", false),
				).
				Value(&emojiState),
			huh.NewSelect[bool]().
				Title(fmt.Sprintf("Do you want to show cot?")).
				Options(
					huh.NewOption("Yes", true),
					huh.NewOption("No", false),
				).
				Value(&cotState),
		).WithTheme(huh.ThemeBase()),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("Failed to get user input: %v\n", err)
	}

	viper.Set("key", keyState)
	viper.Set("model", modelState)
	viper.Set("emoji", emojiState)
	viper.Set("cot", cotState)

	if err := viper.WriteConfigAs(expandedPath); err != nil {
		return fmt.Errorf("Failed to write config file: %v\n", err)
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
		Cot:   viper.GetBool("cot"),
	}
}
