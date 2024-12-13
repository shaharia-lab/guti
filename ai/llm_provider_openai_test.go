package ai

import (
	"testing"

	"github.com/openai/openai-go"
)

func TestOpenAILLMProvider_NewOpenAILLMProvider(t *testing.T) {
	tests := []struct {
		name          string
		config        OpenAIProviderConfig
		expectedModel string
	}{
		{
			name: "with specified model",
			config: OpenAIProviderConfig{
				APIKey: "test-key",
				Model:  "gpt-4",
			},
			expectedModel: "gpt-4",
		},
		{
			name: "with default model",
			config: OpenAIProviderConfig{
				APIKey: "test-key",
			},
			expectedModel: string(openai.ChatModelGPT3_5Turbo),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewOpenAILLMProvider(tt.config)

			if provider.model != tt.expectedModel {
				t.Errorf("expected model %q, got %q", tt.expectedModel, provider.model)
			}
			if provider.client == nil {
				t.Error("expected client to be initialized")
			}
		})
	}
}
