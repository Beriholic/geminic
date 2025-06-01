package service

import (
	"context"
	"encoding/json"

	"github.com/Beriholic/geminic/internal/config"
	md "github.com/Beriholic/geminic/internal/model"
	"google.golang.org/genai"
)

type GeminiService struct {
	Prompt string
	client *genai.Client
}

func NewGeminiServer(
	ctx context.Context,
	commit string,
	diff string,
	files []string,
) (*GeminiService, error) {
	prompt := NewPrompt().Build(commit, diff, files)

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.Get().Key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &GeminiService{
		Prompt: prompt,
		client: client,
	}, nil
}

func NewGeminiServerBlank(
	ctx context.Context,
) (*GeminiService, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.Get().Key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}
	return &GeminiService{
		Prompt: "",
		client: client,
	}, nil
}

func (g *GeminiService) Generate(ctx context.Context) (*md.GitCommit, error) {
	geminiConfig := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   getGenerateStruct(),
	}

	result, err := g.client.Models.GenerateContent(
		ctx,
		config.Get().Model,
		genai.Text(g.Prompt),
		geminiConfig,
	)
	if err != nil {
		return nil, err
	}

	var gitCommit *md.GitCommit

	if err := json.Unmarshal([]byte(result.Text()), &gitCommit); err != nil {
		return nil, err
	}

	return gitCommit, nil
}

func (g *GeminiService) ListModels(ctx context.Context) []string {
	iter := g.client.Models.All(ctx)

	models := []string{}

	for model, err := range iter {
		if err != nil {
			continue
		}

		models = append(models, model.Name)
	}

	return models
}

func getGenerateStruct() *genai.Schema {
	if config.Get().Emoji {
		return &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"typ":   {Type: genai.TypeString, Description: "type of commit"},
				"emoji": {Type: genai.TypeString, Description: "emoji of commit"},
				"scope": {Type: genai.TypeString, Description: "scope of commit"},
				"msg":   {Type: genai.TypeString, Description: "message of commit"},
			},
		}
	}
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"typ":   {Type: genai.TypeString, Description: "type of commit"},
			"scope": {Type: genai.TypeString, Description: "scope of commit"},
			"msg":   {Type: genai.TypeString, Description: "message of commit"},
		},
	}
}
