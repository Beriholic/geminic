package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/viper"
)

const configFilePath = "$HOME/.config/geminic/config.toml"

type GeminicConfig struct {
	Key   string
	Model string
	Emoji bool
}

func Vertify() error {
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

func Create() {
	expandedPath := os.ExpandEnv(configFilePath)

	configDir := filepath.Dir(expandedPath)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		fmt.Printf("Failed to create config directory: %v\n", err)
		return
	}

	viper.SetConfigFile(expandedPath)
	_ = viper.ReadInConfig()

	currentKey := viper.GetString("key")
	currentEmoji := viper.GetBool("emoji")
	currentModel := viper.GetString("model")

	qs := []*survey.Question{
		{
			Name: "key",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("What is your Gemini API key? (cur: %s)", currentKey),
				Default: currentKey,
			},
		},
		{
			Name: "model",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("Which model do you want to use? (cur: %s)", currentModel),
				Default: currentModel,
			},
		},
		{
			Name: "emoji",
			Prompt: &survey.Confirm{
				Message: fmt.Sprintf("Do you want to enable emoji? (cur: %v)", currentEmoji),
				Default: currentEmoji,
			},
		},
	}

	answers := GeminicConfig{}

	if err := survey.Ask(qs, &answers); err != nil {
		fmt.Printf("Failed to get user input: %v\n", err)
		return
	}

	if answers.Key != "" {
		viper.Set("key", answers.Key)
	}
	if answers.Model != "" {
		viper.Set("model", answers.Model)
	}
	viper.Set("emoji", answers.Emoji)

	if err := viper.WriteConfigAs(expandedPath); err != nil {
		fmt.Printf("Failed to write config file: %v\n", err)
		return
	}

	fmt.Printf("Configuration saved to %s\n", expandedPath)
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
