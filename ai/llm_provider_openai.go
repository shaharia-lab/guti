package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OpenAILLMProvider implements the LLMProvider interface for OpenAI's API.
type OpenAILLMProvider struct {
	apiKey     string
	model      string
	apiVersion string
	baseURL    string
}

// OpenAIProviderConfig holds the configuration options for the OpenAI provider.
type OpenAIProviderConfig struct {
	// APIKey is the authentication key for OpenAI's API
	APIKey string
	// Model specifies which OpenAI model to use (e.g., "gpt-3.5-turbo", "gpt-4")
	Model string
	// APIVersion specifies the API version to use (defaults to "v1")
	APIVersion string
	// BaseURL specifies the API endpoint (defaults to "https://api.openai.com")
	BaseURL string
}

// NewOpenAILLMProvider creates a new OpenAI provider with the specified configuration.
func NewOpenAILLMProvider(config OpenAIProviderConfig) *OpenAILLMProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.openai.com"
	}
	if config.APIVersion == "" {
		config.APIVersion = "v1"
	}
	if config.Model == "" {
		config.Model = "gpt-3.5-turbo"
	}

	return &OpenAILLMProvider{
		apiKey:     config.APIKey,
		model:      config.Model,
		apiVersion: config.APIVersion,
		baseURL:    config.BaseURL,
	}
}

type openAIRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	TopP        float64       `json:"top_p,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []choice `json:"choices"`
	Usage   usage    `json:"usage"`
}

type choice struct {
	Index        int         `json:"index"`
	Message      chatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// GetResponse implements the LLMProvider interface for OpenAI.
func (p *OpenAILLMProvider) GetResponse(question string, config LLMRequestConfig) (LLMResponse, error) {
	startTime := time.Now()

	messages := []chatMessage{
		{
			Role:    "user",
			Content: question,
		},
	}

	requestBody := openAIRequest{
		Model:       p.model,
		Messages:    messages,
		MaxTokens:   config.MaxToken,
		Temperature: config.Temperature,
		TopP:        config.TopP,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return LLMResponse{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	endpoint := fmt.Sprintf("%s/%s/chat/completions", p.baseURL, p.apiVersion)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return LLMResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return LLMResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    int    `json:"code"`
			} `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return LLMResponse{}, &LLMError{
				Code:    resp.StatusCode,
				Message: fmt.Sprintf("API request failed with status %d", resp.StatusCode),
			}
		}
		return LLMResponse{}, &LLMError{
			Code:    errorResponse.Error.Code,
			Message: fmt.Sprintf("%s (type: %s)", errorResponse.Error.Message, errorResponse.Error.Type),
		}
	}

	var openAIResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return LLMResponse{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return LLMResponse{}, fmt.Errorf("no choices in response")
	}

	return LLMResponse{
		Text:             openAIResp.Choices[0].Message.Content,
		TotalInputToken:  openAIResp.Usage.PromptTokens,
		TotalOutputToken: openAIResp.Usage.CompletionTokens,
		CompletionTime:   time.Since(startTime).Seconds(),
	}, nil
}
