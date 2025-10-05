package service

import (
	"context"

	"github.com/Beriholic/geminic/internal/config"
	"github.com/Beriholic/geminic/internal/llm"
	"github.com/Beriholic/geminic/internal/llm/prompt"
	"github.com/Beriholic/geminic/internal/model/dto"
)

type LLMService struct {
	LLM llm.LLM
}

func init() {
}

func NewLLMServer(ctx context.Context) (*LLMService, error) {
	cfg := config.Get()
	llm, err := llm.GetLLM(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &LLMService{
		LLM: llm,
	}, nil
}

func (l *LLMService) Generate(ctx context.Context, dto *dto.CommitDTO) (*dto.GitCommit, error) {
	prompt := prompt.NewPrompt().Build(dto)
	return l.LLM.Generate(ctx, prompt)
}

func (l *LLMService) ModelList(ctx context.Context) ([]string, error) {
	return l.LLM.ModelList(ctx)
}
