package llm

import (
	"context"

	"github.com/Beriholic/geminic/internal/model"
	"github.com/Beriholic/geminic/internal/model/model_provider"
)

func GetLLM(ctx context.Context, cfg *model.Config) (LLM, error) {
	if cfg.ModelProvider == model_provider.Gemini {
		return NewGeminiLLM(ctx)
	}
	return NewOpenAILLM(ctx)
}
