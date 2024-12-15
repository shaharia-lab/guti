package ai

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
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

func TestOpenAILLMProvider_GetStreamingResponse(t *testing.T) {
	tests := []struct {
		name     string
		messages []LLMMessage
		timeout  time.Duration
		delay    time.Duration
		wantErr  bool
	}{
		{
			name:    "successful streaming",
			timeout: 100 * time.Millisecond,
			delay:   0,
		},
		{
			name:    "context cancellation",
			timeout: 5 * time.Millisecond,
			delay:   50 * time.Millisecond,
		},
	}

	responses := []string{
		`data: {"id":"123","choices":[{"delta":{"content":"Hello"}}]}`,
		`data: {"id":"123","choices":[{"delta":{"content":" world"}}]}`,
		`data: [DONE]`,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &OpenAILLMProvider{
				client: openai.NewClient(
					option.WithHTTPClient(&http.Client{
						Transport: &mockTransport{
							responses: responses,
							delay:     tt.delay,
						},
					}),
				),
				model: openai.ChatModelGPT3_5Turbo,
			}

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			stream, err := provider.GetStreamingResponse(ctx, []LLMMessage{{Role: UserRole, Text: "test"}}, LLMRequestConfig{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var gotCancel bool
			for resp := range stream {
				if resp.Error != nil {
					gotCancel = true
					break
				}
			}

			if tt.delay > tt.timeout && !gotCancel {
				t.Error("expected context cancellation")
			}
		})
	}
}

type mockTransport struct {
	responses []string
	delay     time.Duration
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		for _, resp := range m.responses {
			time.Sleep(10 * time.Millisecond) // Simulate streaming delay
			pw.Write([]byte(resp + "\n"))
		}
	}()

	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: pr,
	}, nil
}
