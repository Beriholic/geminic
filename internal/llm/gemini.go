package llm

import (
	"context"
	"encoding/json"

	"github.com/Beriholic/geminic/internal/config"
	"github.com/Beriholic/geminic/internal/model/dto"
	"google.golang.org/genai"
)

type GeminiLLM struct {
	client *genai.Client
}

func NewGeminiLLM(ctx context.Context) (*GeminiLLM, error) {
	if config.Get().CustomURL != "" {
		genai.SetDefaultBaseURLs(genai.BaseURLParameters{
			GeminiURL: config.Get().CustomURL,
			VertexURL: config.Get().CustomURL,
		})
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.Get().Key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &GeminiLLM{client: client}, nil
}

func (g *GeminiLLM) Generate(ctx context.Context, pmt string) (*dto.GitCommit, error) {
	geminiConfig := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   dto.GitCommit{}.ToGeminiGenerateStruct(),
	}

	result, err := g.client.Models.GenerateContent(
		ctx,
		config.Get().Model,
		genai.Text(pmt),
		geminiConfig,
	)
	if err != nil {
		return nil, err
	}

	var gitCommit *dto.GitCommit

	if err := json.Unmarshal([]byte(result.Text()), &gitCommit); err != nil {
		return nil, err
	}

	return gitCommit, nil
}

func (g *GeminiLLM) ModelList(ctx context.Context) ([]string, error) {
	var models []string
	iter := g.client.Models.All(ctx)
	for model, err := range iter {
		if err != nil {
			return nil, err
		}
		models = append(models, model.Name)
	}
	return models, nil
}
