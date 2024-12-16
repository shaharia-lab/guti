// File: ai/llm_provider_anthropic_client.go

package ai

import (
	"context"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// AnthropicClient defines the interface for interacting with Anthropic's API.
// This interface abstracts the essential message-related operations used by AnthropicLLMProvider.
type AnthropicClient interface {
	// CreateMessage creates a new message using Anthropic's API.
	// The method takes a context and MessageNewParams and returns a Message response or an error.
	CreateMessage(ctx context.Context, params anthropic.MessageNewParams) (*anthropic.Message, error)
}

// RealAnthropicClient implements the AnthropicClient interface using Anthropic's official SDK.
type RealAnthropicClient struct {
	messages *anthropic.MessageService
}

// NewRealAnthropicClient creates a new instance of RealAnthropicClient with the provided API key.
//
// Example usage:
//
//	client := NewRealAnthropicClient("your-api-key")
//	provider := NewAnthropicLLMProvider(AnthropicProviderConfig{
//	    Client: client,
//	    Model:  "claude-3-sonnet-20240229",
//	})
func NewRealAnthropicClient(apiKey string) *RealAnthropicClient {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &RealAnthropicClient{
		messages: client.Messages,
	}
}

// CreateMessage implements the AnthropicClient interface using the real Anthropic client.
func (c *RealAnthropicClient) CreateMessage(ctx context.Context, params anthropic.MessageNewParams) (*anthropic.Message, error) {
	return c.messages.New(ctx, params)
}
