// Package ai provides a flexible interface for interacting with various Language Learning Models (LLMs).
package ai

// LLMRequest handles the configuration and execution of LLM requests.
// It provides a consistent interface for interacting with different LLM providers.
type LLMRequest struct {
	requestConfig LLMRequestConfig
}

// NewLLMRequest creates a new LLMRequest with the specified configuration.
func NewLLMRequest(requestConfig LLMRequestConfig) *LLMRequest {
	return &LLMRequest{
		requestConfig: requestConfig,
	}
}

// Generate sends a prompt to the specified LLM provider and returns the response.
// Returns LLMResponse containing the generated text and metadata, or an error if the operation fails.
func (r *LLMRequest) Generate(messages []LLMMessage, llmProvider LLMProvider) (LLMResponse, error) {
	return llmProvider.GetResponse(messages, r.requestConfig)
}
