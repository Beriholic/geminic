package config

import (
	"fmt"
	"sync"

	"github.com/Beriholic/geminic/internal/model"
	"github.com/Beriholic/geminic/internal/model/model_provider"
	"github.com/charmbracelet/huh"
)

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
	configOnce sync.Once
	config     *model.Config = nil
)

func Get() *model.Config {
	var err error
	configOnce.Do(func() {
		config, err = load()
	})
	if err != nil {
		fmt.Printf("failed to load config: %v", err)
	}
	return config
}

func Create() error {
	var config model.Config
	config.Load()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is your Gemini API key?").
				Value(&config.Key),
			huh.NewInput().
				Title("Which model do you want to use?").
				Value(&config.Model),
			huh.NewSelect[bool]().
				Title("Do you want to enable emoji?").
				Options(
					huh.NewOption("Yes", true),
					huh.NewOption("No", false),
				).
				Value(&config.Emoji),
			huh.NewInput().
				Title("Custom backend connection (leave blank to disable)").
				Value(&config.CustomURL),
			huh.NewInput().
				Title("i18n").
				Value(&config.I18n),
			huh.NewSelect[string]().
				Title("Model Provider").
				Options(
					huh.NewOption(model_provider.Gemini, model_provider.Gemini),
					huh.NewOption(model_provider.OpenAI, model_provider.OpenAI),
				).
				Value(&config.ModelProvider),
		).WithTheme(huh.ThemeBase()),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	return config.Save()
}

func load() (*model.Config, error) {
	var config model.Config
	err := config.Load()
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func SetModel(model string) error {
	config := Get()
	config.Model = model
	return config.Save()
}
