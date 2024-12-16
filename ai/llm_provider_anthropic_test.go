package ai

import (
	"context"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/stretchr/testify/assert"
)

// MockAnthropicClient implements AnthropicClient interface for testing
type MockAnthropicClient struct {
	createMessageFunc func(ctx context.Context, params anthropic.MessageNewParams) (*anthropic.Message, error)
}

func (m *MockAnthropicClient) CreateMessage(ctx context.Context, params anthropic.MessageNewParams) (*anthropic.Message, error) {
	return m.createMessageFunc(ctx, params)
}

func TestAnthropicLLMProvider_NewAnthropicLLMProvider(t *testing.T) {
	tests := []struct {
		name          string
		config        AnthropicProviderConfig
		expectedModel anthropic.Model
	}{
		{
			name: "with specified model",
			config: AnthropicProviderConfig{
				Client: &MockAnthropicClient{},
				Model:  "claude-3-opus-20240229",
			},
			expectedModel: "claude-3-opus-20240229",
		},
		{
			name: "with default model",
			config: AnthropicProviderConfig{
				Client: &MockAnthropicClient{},
			},
			expectedModel: anthropic.ModelClaude_3_5_Sonnet_20240620,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewAnthropicLLMProvider(tt.config)

			assert.Equal(t, tt.expectedModel, provider.model, "unexpected model")
			assert.NotNil(t, provider.client, "expected client to be initialized")
		})
	}
}

func TestAnthropicLLMProvider_GetResponse(t *testing.T) {
	tests := []struct {
		name           string
		messages       []LLMMessage
		config         LLMRequestConfig
		expectedResult LLMResponse
		expectError    bool
	}{
		{
			name: "successful response with all message types",
			messages: []LLMMessage{
				{Role: SystemRole, Text: "You are a helpful assistant"},
				{Role: UserRole, Text: "Hello"},
				{Role: AssistantRole, Text: "Hi there"},
			},
			config: LLMRequestConfig{
				MaxToken:    100,
				TopP:        0.9,
				Temperature: 0.7,
			},
			expectedResult: LLMResponse{
				Text:             "Test response",
				TotalInputToken:  10,
				TotalOutputToken: 5,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAnthropicClient{
				createMessageFunc: func(ctx context.Context, params anthropic.MessageNewParams) (*anthropic.Message, error) {
					message := &anthropic.Message{
						Role:  anthropic.MessageRoleAssistant,
						Model: anthropic.ModelClaude_3_5_Sonnet_20240620,
						Usage: anthropic.Usage{
							InputTokens:  10,
							OutputTokens: 5,
						},
						Type: anthropic.MessageTypeMessage,
					}

					block := anthropic.ContentBlock{}
					if err := block.UnmarshalJSON([]byte(`{
						"type": "text",
						"text": "Test response"
					}`)); err != nil {
						t.Fatal(err)
					}

					message.Content = []anthropic.ContentBlock{block}
					return message, nil
				},
			}

			provider := NewAnthropicLLMProvider(AnthropicProviderConfig{
				Client: mockClient,
				Model:  anthropic.ModelClaude_3_5_Sonnet_20240620,
			})

			result, err := provider.GetResponse(tt.messages, tt.config)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult.Text, result.Text)
			assert.Equal(t, tt.expectedResult.TotalInputToken, result.TotalInputToken)
			assert.Equal(t, tt.expectedResult.TotalOutputToken, result.TotalOutputToken)
			assert.Greater(t, result.CompletionTime, float64(0), "completion time should be greater than 0")
		})
	}
}
