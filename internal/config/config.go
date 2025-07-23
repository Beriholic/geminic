package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
)

const configFilePath = "$HOME/.config/geminic/config.toml"

type GeminicConfig struct {
	Key       string `mapstructure:"key"`
	Model     string `mapstructure:"model"`
	Emoji     bool   `mapstructure:"emoji"`
	CustomUrl string `mapstructure:"custom_url"`
	I18n      string `mapstructure:"i18n"`
}
type GeminiLocalConfig struct {
	Emoji bool   `mapstructure:"emoji"`
	I18n  string `mapstructure:"i18n"`
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
	var err error
	geminicConfigOnce.Do(func() {
		geminicConfig, err = load()
	})
	if err != nil {
		fmt.Printf("failed to load config: %v", err)
	}
	return geminicConfig
}

func Create() error {
	expandedPath := os.ExpandEnv(configFilePath)

	configDir := filepath.Dir(expandedPath)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	viper.SetConfigFile(expandedPath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to red config file: %v", err)
	}

	keyState := viper.GetString("key")
	modelState := viper.GetString("model")
	emojiState := viper.GetBool("emoji")
	customUrl := viper.GetString("custom_url")

	i18n := viper.GetString("i18n")
	if i18n == "" {
		i18n = "en"
	}

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
			huh.NewInput().
				Title("Custom backend connection (leave blank to disable)").
				Value(&customUrl),
			huh.NewInput().
				Title("i18n").
				Value(&i18n),
		).WithTheme(huh.ThemeBase()),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	viper.Set("key", keyState)
	viper.Set("model", modelState)
	viper.Set("emoji", emojiState)
	viper.Set("custom_url", customUrl)
	viper.Set("i18n", i18n)

	if err := viper.WriteConfigAs(expandedPath); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	fmt.Printf("Configuration saved to %s\n", expandedPath)
	return nil
}

func CreateLocal() error {
	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	configPath := filepath.Join(currentPath, "geminic.toml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		_, err := os.Create(configPath)
		if err != nil {
			return err
		}
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	emojiState := viper.GetBool("emoji")
	i18n := viper.GetString("i18n")
	if i18n == "" {
		i18n = "en"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("Do you want to enable emoji?").
				Options(
					huh.NewOption("Yes", true),
					huh.NewOption("No", false),
				).
				Value(&emojiState),
			huh.NewInput().
				Title("i18n").
				Value(&i18n),
		).WithTheme(huh.ThemeBase()),
	)
	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	viper.Set("emoji", emojiState)
	viper.Set("i18n", i18n)

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	fmt.Printf("Configuration saved to %s\n", configPath)
	return nil
}

func load() (*GeminicConfig, error) {
	v := viper.New()

	expandedPath := os.ExpandEnv(configFilePath)
	v.SetConfigFile(expandedPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	currentPath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %v", err)
	}
	localConfigFile := filepath.Join(currentPath, "geminic.toml")

	if _, err := os.Stat(localConfigFile); err == nil {
		v.SetConfigFile(localConfigFile)

		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("failed to merge local config file: %v", err)
		}
	}

	var cfg GeminicConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return &cfg, nil
}

func SetModel(model string) error {
	expandedPath := os.ExpandEnv(configFilePath)
	viper.SetConfigFile(expandedPath)

	model = strings.TrimPrefix(model, "models/")

	viper.Set("model", model)

	if err := viper.WriteConfigAs(expandedPath); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	fmt.Printf("Configuration saved to %s\n", expandedPath)
	return nil
}

func (c *GeminicConfig) GetCustomModels() ([]string, error) {
	type Model struct {
		Id string `json:"id"`
	}

	type ResponseData struct {
		Data []Model `json:"data"`
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/models", c.CustomUrl), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Key))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData ResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, err
	}

	var ids []string
	for _, m := range responseData.Data {
		ids = append(ids, m.Id)
	}

	return ids, nil
}
