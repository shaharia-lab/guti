package ai

import (
	"errors"
	"testing"
)

type mockProvider struct {
	response LLMResponse
	err      error
}

func (m *mockProvider) GetResponse(messages []LLMMessage, _ LLMRequestConfig) (LLMResponse, error) {
	return m.response, m.err
}

func TestLLMRequest_Generate(t *testing.T) {
	tests := []struct {
		name           string
		config         LLMRequestConfig
		mockResponse   LLMResponse
		mockError      error
		expectedError  bool
		expectedOutput LLMResponse
	}{
		{
			name: "successful generation",
			config: LLMRequestConfig{
				MaxToken:    100,
				TopP:        0.9,
				Temperature: 0.7,
			},
			mockResponse: LLMResponse{
				Text:             "Hello, world!",
				TotalInputToken:  5,
				TotalOutputToken: 3,
				CompletionTime:   0.5,
			},
			expectedOutput: LLMResponse{
				Text:             "Hello, world!",
				TotalInputToken:  5,
				TotalOutputToken: 3,
				CompletionTime:   0.5,
			},
		},
		{
			name: "provider error",
			config: LLMRequestConfig{
				MaxToken:    100,
				TopP:        0.9,
				Temperature: 0.7,
			},
			mockError:     errors.New("provider error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &mockProvider{
				response: tt.mockResponse,
				err:      tt.mockError,
			}

			request := NewLLMRequest(tt.config)
			response, err := request.Generate([]LLMMessage{{
				Role: "user",
				Text: "test prompt",
			}}, provider)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if response.Text != tt.expectedOutput.Text {
				t.Errorf("expected text %q, got %q", tt.expectedOutput.Text, response.Text)
			}
			if response.TotalInputToken != tt.expectedOutput.TotalInputToken {
				t.Errorf("expected input tokens %d, got %d", tt.expectedOutput.TotalInputToken, response.TotalInputToken)
			}
			if response.TotalOutputToken != tt.expectedOutput.TotalOutputToken {
				t.Errorf("expected output tokens %d, got %d", tt.expectedOutput.TotalOutputToken, response.TotalOutputToken)
			}
		})
	}
}
