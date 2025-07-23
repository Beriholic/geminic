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
	Key       string
	Model     string
	Emoji     bool
	CustomUrl string
	I18n      string
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

func load() *GeminicConfig {
	expandedPath := os.ExpandEnv(configFilePath)

	viper.SetConfigFile(expandedPath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		return nil
	}

	return &GeminicConfig{
		Key:       viper.GetString("key"),
		Model:     viper.GetString("model"),
		Emoji:     viper.GetBool("emoji"),
		CustomUrl: viper.GetString("custom_url"),
		I18n:      viper.GetString("i18n"),
	}
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
