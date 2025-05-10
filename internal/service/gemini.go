package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Beriholic/geminic/internal/config"
	md "github.com/Beriholic/geminic/internal/model"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
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

	client, err := genai.NewClient(ctx,
		option.WithAPIKey(config.Get().Key),
	)
	if err != nil {
		return nil, err
	}

	return &GeminiService{
		Prompt: prompt,
		client: client,
	}, nil
}

func (g *GeminiService) Generate(ctx context.Context) (*md.GitCommit, error) {
	model := g.client.GenerativeModel(config.Get().Model)

	model.SetTemperature(1.02)
	model.ResponseMIMEType = "application/json"

	resp, err := model.GenerateContent(ctx, genai.Text(g.Prompt))
	if err != nil {
		return nil, err
	}

	jsonStr := fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])

	var gitCommit *md.GitCommit

	if err := json.Unmarshal([]byte(jsonStr), &gitCommit); err != nil {
		return nil, err
	}

	return gitCommit, nil
}

func (g *GeminiService) CloseClient() {
	g.client.Close()
}
