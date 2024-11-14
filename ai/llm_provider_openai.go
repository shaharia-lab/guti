package ai

import (
	"context"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAILLMProvider implements the LLMProvider interface using the official OpenAI SDK
type OpenAILLMProvider struct {
	client *openai.Client
	model  string
}

// OpenAIProviderConfig holds configuration for OpenAI provider
type OpenAIProviderConfig struct {
	APIKey string
	Model  string
}

// NewOpenAILLMProvider creates a new OpenAI provider using the official SDK
func NewOpenAILLMProvider(config OpenAIProviderConfig) *OpenAILLMProvider {
	if config.Model == "" {
		config.Model = openai.ChatModelGPT3_5Turbo
	}

	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
	)

	return &OpenAILLMProvider{
		client: client,
		model:  config.Model,
	}
}

// GetResponse implements the LLMProvider interface
func (p *OpenAILLMProvider) GetResponse(question string, config LLMRequestConfig) (LLMResponse, error) {
	startTime := time.Now()

	params := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		}),
		Model:       openai.F(openai.ChatModel(p.model)),
		MaxTokens:   openai.Int(int64(config.MaxToken)),
		TopP:        openai.Float(config.TopP),
		Temperature: openai.Float(config.Temperature),
	}

	completion, err := p.client.Chat.Completions.New(context.Background(), params)
	if err != nil {
		return LLMResponse{}, err
	}

	if len(completion.Choices) == 0 {
		return LLMResponse{}, &LLMError{
			Code:    400,
			Message: "no choices in response",
		}
	}

	return LLMResponse{
		Text:             completion.Choices[0].Message.Content,
		TotalInputToken:  int(completion.Usage.PromptTokens),
		TotalOutputToken: int(completion.Usage.CompletionTokens),
		CompletionTime:   time.Since(startTime).Seconds(),
	}, nil
}
