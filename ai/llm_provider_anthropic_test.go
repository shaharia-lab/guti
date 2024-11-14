package ai

import (
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

func TestAnthropicLLMProvider_NewAnthropicLLMProvider(t *testing.T) {
	tests := []struct {
		name          string
		config        AnthropicProviderConfig
		expectedModel string
	}{
		{
			name: "with specified model",
			config: AnthropicProviderConfig{
				APIKey: "test-key",
				Model:  "claude-3-opus-20240229",
			},
			expectedModel: "claude-3-opus-20240229",
		},
		{
			name: "with default model",
			config: AnthropicProviderConfig{
				APIKey: "test-key",
			},
			expectedModel: string(anthropic.ModelClaude_3_5_Sonnet_20240620),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewAnthropicLLMProvider(tt.config)

			if provider.model != tt.expectedModel {
				t.Errorf("expected model %q, got %q", tt.expectedModel, provider.model)
			}
			if provider.client == nil {
				t.Error("expected client to be initialized")
			}
		})
	}
}
