package llm

import (
	"context"
	"fmt"

	"github.com/Beriholic/geminic/internal/config"
	"github.com/Beriholic/geminic/internal/model/dto"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type OpenAILLM struct {
	client *openai.Client
}

func NewOpenAILLM(ctx context.Context) (*OpenAILLM, error) {
	apiConfig := openai.DefaultConfig(config.Get().Key)

	if config.Get().CustomURL != "" {
		apiConfig.BaseURL = config.Get().CustomURL
	}

	client := openai.NewClientWithConfig(apiConfig)
	return &OpenAILLM{client: client}, nil
}

func (o *OpenAILLM) Generate(ctx context.Context, prompt string) (*dto.GitCommit, error) {
	var gitCommit dto.GitCommit

	schema, err := jsonschema.GenerateSchemaForType(gitCommit)
	if err != nil {
		return nil, err
	}

	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       config.Get().Model,
		Temperature: 0.75,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Start writing a Git commit",
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "result of git commit",
				Schema: schema,
				Strict: true,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if resp.Choices == nil {
		return nil, fmt.Errorf("Blank repley")
	}

	content := resp.Choices[0].Message.Content
	err = schema.Unmarshal(content, &gitCommit)
	if err != nil {
		return nil, fmt.Errorf("json: %v err: %v", content, err)
	}
	return &gitCommit, nil
}

func (o *OpenAILLM) ModelList(ctx context.Context) ([]string, error) {
	var models []string
	_models, err := o.client.ListModels(ctx)
	if err != nil {
		return nil, err
	}

	for _, model := range _models.Models {
		models = append(models, model.ID)
	}

	return models, nil
}
