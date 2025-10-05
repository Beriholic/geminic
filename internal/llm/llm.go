package llm

import (
	"context"

	"github.com/Beriholic/geminic/internal/model/dto"
)

type LLM interface {
	Generate(ctx context.Context, prompt string) (*dto.GitCommit, error)
	ModelList(ctx context.Context) ([]string, error)
}
