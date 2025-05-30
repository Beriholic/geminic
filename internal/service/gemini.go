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

func (g *GeminiService) Generate(ctx context.Context) (*md.GitCommit, error) {
	geminiConfig := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"typ":   {Type: genai.TypeString},
				"emoji": {Type: genai.TypeString},
				"scope": {Type: genai.TypeString},
				"msg":   {Type: genai.TypeString},
			},
		},
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
