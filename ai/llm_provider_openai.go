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

// NewOpenAILLMProvider creates a new OpenAI provider with the specified configuration.
// If no model is specified, it defaults to GPT-3.5-turbo.
func NewOpenAILLMProvider(config OpenAIProviderConfig) *OpenAILLMProvider {
	if config.Model == "" {
		config.Model = string(openai.ChatModelGPT3_5Turbo)
	}

	return &OpenAILLMProvider{
		client: openai.NewClient(option.WithAPIKey(config.APIKey)),
		model:  config.Model,
	}
}

// GetResponse generates a response using OpenAI's API for the given messages and configuration.
// It supports different message roles (user, assistant, system) and handles them appropriately.
func (p *OpenAILLMProvider) GetResponse(messages []LLMMessage, config LLMRequestConfig) (LLMResponse, error) {
	startTime := time.Now()

	var openAIMessages []openai.ChatCompletionMessageParamUnion
	for _, msg := range messages {
		switch msg.Role {
		case "user":
			openAIMessages = append(openAIMessages, openai.UserMessage(msg.Text))
		case "assistant":
			openAIMessages = append(openAIMessages, openai.AssistantMessage(msg.Text))
		case "system":
			openAIMessages = append(openAIMessages, openai.SystemMessage(msg.Text))
		default:
			openAIMessages = append(openAIMessages, openai.UserMessage(msg.Text))
		}
	}

	params := openai.ChatCompletionNewParams{
		Messages:    openai.F(openAIMessages),
		Model:       openai.F(openai.ChatModel(p.model)),
		MaxTokens:   openai.Int(config.MaxToken),
		TopP:        openai.Float(config.TopP),
		Temperature: openai.Float(config.Temperature),
	}

	completion, err := p.client.Chat.Completions.New(context.Background(), params)
	if err != nil {
		return LLMResponse{}, err
	}

	if len(completion.Choices) == 0 {
		return LLMResponse{}, &LLMError{Code: 400, Message: "no choices in response"}
	}

	return LLMResponse{
		Text:             completion.Choices[0].Message.Content,
		TotalInputToken:  int(completion.Usage.PromptTokens),
		TotalOutputToken: int(completion.Usage.CompletionTokens),
		CompletionTime:   time.Since(startTime).Seconds(),
	}, nil
}
