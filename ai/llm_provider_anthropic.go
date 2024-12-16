// File: ai/llm_provider_anthropic.go

package ai

import (
	"context"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
)

// AnthropicLLMProvider implements the LLMProvider interface using Anthropic's official Go SDK.
// It provides access to Claude models through Anthropic's API.
type AnthropicLLMProvider struct {
	client AnthropicClient
	model  anthropic.Model
}

// AnthropicProviderConfig holds the configuration options for creating an Anthropic provider.
type AnthropicProviderConfig struct {
	// Client is the AnthropicClient implementation to use
	Client AnthropicClient

	// Model specifies which Anthropic model to use (e.g., "claude-3-opus-20240229", "claude-3-sonnet-20240229")
	Model anthropic.Model
}

// NewAnthropicLLMProvider creates a new Anthropic provider with the specified configuration.
// If no model is specified, it defaults to Claude 3.5 Sonnet.
//
// Example usage:
//
//	client := NewRealAnthropicClient("your-api-key")
//	provider := NewAnthropicLLMProvider(AnthropicProviderConfig{
//	    Client: client,
//	    Model:  anthropic.ModelClaude_3_5_Sonnet_20240620,
//	})
//
//	response, err := provider.GetResponse(messages, config)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewAnthropicLLMProvider(config AnthropicProviderConfig) *AnthropicLLMProvider {
	if config.Model == "" {
		config.Model = anthropic.ModelClaude_3_5_Sonnet_20240620
	}

	return &AnthropicLLMProvider{
		client: config.Client,
		model:  config.Model,
	}
}

// GetResponse generates a response using Anthropic's API for the given messages and configuration.
// It supports different message roles (user, assistant, system) and handles them appropriately.
// System messages are handled separately through Anthropic's system parameter.
func (p *AnthropicLLMProvider) GetResponse(messages []LLMMessage, config LLMRequestConfig) (LLMResponse, error) {
	startTime := time.Now()

	var anthropicMessages []anthropic.MessageParam
	for _, msg := range messages {
		switch msg.Role {
		case UserRole:
			anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Text)))
		case AssistantRole:
			anthropicMessages = append(anthropicMessages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Text)))
		case SystemRole:
			continue
		default:
			anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Text)))
		}
	}

	params := anthropic.MessageNewParams{
		Model:       anthropic.F(p.model),
		Messages:    anthropic.F(anthropicMessages),
		MaxTokens:   anthropic.F(config.MaxToken),
		TopP:        anthropic.Float(config.TopP),
		Temperature: anthropic.Float(config.Temperature),
	}

	// Add system message if present
	for _, msg := range messages {
		if msg.Role == SystemRole {
			params.System = anthropic.F([]anthropic.TextBlockParam{
				anthropic.NewTextBlock(msg.Text),
			})
			break
		}
	}

	message, err := p.client.CreateMessage(context.Background(), params)
	if err != nil {
		return LLMResponse{}, err
	}

	var responseText string
	for _, block := range message.Content {
		switch block := block.AsUnion().(type) {
		case anthropic.TextBlock:
			responseText += block.Text
		}
	}

	return LLMResponse{
		Text:             responseText,
		TotalInputToken:  int(message.Usage.InputTokens),
		TotalOutputToken: int(message.Usage.OutputTokens),
		CompletionTime:   time.Since(startTime).Seconds(),
	}, nil
}
