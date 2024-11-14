package ai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAILLMProvider_NewOpenAILLMProvider(t *testing.T) {
	tests := []struct {
		name     string
		config   OpenAIProviderConfig
		expected *OpenAILLMProvider
	}{
		{
			name: "with all values",
			config: OpenAIProviderConfig{
				APIKey:     "test-key",
				Model:      "gpt-4",
				APIVersion: "v2",
				BaseURL:    "https://custom.openai.com",
			},
			expected: &OpenAILLMProvider{
				apiKey:     "test-key",
				model:      "gpt-4",
				apiVersion: "v2",
				baseURL:    "https://custom.openai.com",
			},
		},
		{
			name: "with minimal values",
			config: OpenAIProviderConfig{
				APIKey: "test-key",
			},
			expected: &OpenAILLMProvider{
				apiKey:     "test-key",
				model:      "gpt-3.5-turbo",
				apiVersion: "v1",
				baseURL:    "https://api.openai.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewOpenAILLMProvider(tt.config)

			if provider.apiKey != tt.expected.apiKey {
				t.Errorf("expected apiKey %q, got %q", tt.expected.apiKey, provider.apiKey)
			}
			if provider.model != tt.expected.model {
				t.Errorf("expected model %q, got %q", tt.expected.model, provider.model)
			}
			if provider.apiVersion != tt.expected.apiVersion {
				t.Errorf("expected apiVersion %q, got %q", tt.expected.apiVersion, provider.apiVersion)
			}
			if provider.baseURL != tt.expected.baseURL {
				t.Errorf("expected baseURL %q, got %q", tt.expected.baseURL, provider.baseURL)
			}
		})
	}
}

func TestOpenAILLMProvider_GetResponse(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse *openAIResponse
		serverError    *LLMError
		serverStatus   int
		config         LLMRequestConfig
		expectedError  bool
	}{
		{
			name: "successful response",
			serverResponse: &openAIResponse{
				Choices: []choice{
					{
						Message: chatMessage{
							Content: "Hello, world!",
						},
					},
				},
				Usage: usage{
					PromptTokens:     5,
					CompletionTokens: 3,
					TotalTokens:      8,
				},
			},
			serverStatus: http.StatusOK,
			config: LLMRequestConfig{
				MaxToken:    100,
				TopP:        0.9,
				Temperature: 0.7,
			},
		},
		{
			name: "api error",
			serverError: &LLMError{
				Code:    400,
				Message: "invalid request",
			},
			serverStatus:  http.StatusBadRequest,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request headers
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected Content-Type header 'application/json', got %q", r.Header.Get("Content-Type"))
				}
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Errorf("expected Authorization header 'Bearer test-key', got %q", r.Header.Get("Authorization"))
				}

				w.WriteHeader(tt.serverStatus)
				if tt.serverError != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error": map[string]interface{}{
							"message": tt.serverError.Message,
							"type":    "error",
							"code":    tt.serverError.Code,
						},
					})
					return
				}

				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			provider := NewOpenAILLMProvider(OpenAIProviderConfig{
				APIKey:  "test-key",
				BaseURL: server.URL,
			})

			response, err := provider.GetResponse("test question", tt.config)

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

			if response.Text != tt.serverResponse.Choices[0].Message.Content {
				t.Errorf("expected text %q, got %q", tt.serverResponse.Choices[0].Message.Content, response.Text)
			}
			if response.TotalInputToken != tt.serverResponse.Usage.PromptTokens {
				t.Errorf("expected input tokens %d, got %d", tt.serverResponse.Usage.PromptTokens, response.TotalInputToken)
			}
			if response.TotalOutputToken != tt.serverResponse.Usage.CompletionTokens {
				t.Errorf("expected output tokens %d, got %d", tt.serverResponse.Usage.CompletionTokens, response.TotalOutputToken)
			}
			if response.CompletionTime <= 0 {
				t.Error("expected positive completion time")
			}
		})
	}
}
