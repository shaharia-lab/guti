// Package ai provides a flexible interface for interacting with various Language Learning Models (LLMs).
package ai

import "fmt"

// LLMMessageRole represents the role of a message in a conversation.
type LLMMessageRole string

const (
	// UserRole represents a message from the user
	UserRole LLMMessageRole = "user"

	// AssistantRole represents a message from the assistant
	AssistantRole LLMMessageRole = "assistant"

	// SystemRole represents a message from the system
	SystemRole LLMMessageRole = "system"
)

// DefaultConfig holds the default values for LLMRequestConfig
var DefaultConfig = LLMRequestConfig{
	MaxToken:    1000, // Default max tokens
	TopP:        0.9,  // Default top-p value
	Temperature: 0.7,  // Default temperature
	TopK:        50,   // Default top-k value
}

// LLMRequestConfig defines configuration parameters for LLM requests.
type LLMRequestConfig struct {
	MaxToken    int64
	TopP        float64
	Temperature float64
	TopK        int64
}

// NewRequestConfig creates a new config with default values.
// Any non-zero values in the provided config will override the defaults.
func NewRequestConfig(opts ...RequestOption) LLMRequestConfig {
	config := DefaultConfig

	// Apply any provided options
	for _, opt := range opts {
		opt(&config)
	}

	return config
}

// RequestOption is a function that modifies the config
type RequestOption func(*LLMRequestConfig)

// WithMaxToken sets the max token value
func WithMaxToken(maxToken int64) RequestOption {
	return func(c *LLMRequestConfig) {
		if maxToken > 0 {
			c.MaxToken = maxToken
		}
	}
}

// WithTopP sets the top-p value
func WithTopP(topP float64) RequestOption {
	return func(c *LLMRequestConfig) {
		if topP > 0 {
			c.TopP = topP
		}
	}
}

// WithTemperature sets the temperature value
func WithTemperature(temp float64) RequestOption {
	return func(c *LLMRequestConfig) {
		if temp > 0 {
			c.Temperature = temp
		}
	}
}

// WithTopK sets the top-k value
func WithTopK(topK int64) RequestOption {
	return func(c *LLMRequestConfig) {
		if topK > 0 {
			c.TopK = topK
		}
	}
}

// LLMResponse encapsulates the response from an LLM provider.
// It includes both the generated text and metadata about the request.
type LLMResponse struct {
	// Text contains the generated response from the model
	Text string

	// TotalInputToken is the number of tokens in the input prompt
	TotalInputToken int

	// TotalOutputToken is the number of tokens in the generated response
	TotalOutputToken int

	// CompletionTime is the total time taken to generate the response in seconds
	CompletionTime float64
}

// LLMError represents errors that occur during LLM operations.
// It provides structured error information including an error code.
type LLMError struct {
	// Code represents the error code (usually HTTP status code for API errors)
	Code int

	// Message provides a detailed description of the error
	Message string
}

// Error implements the error interface for LLMError.
func (e *LLMError) Error() string {
	return fmt.Sprintf("LLMError %d: %s", e.Code, e.Message)
}

// LLMMessage represents a message in a conversation with an LLM.
// It includes the role of the speaker (user, assistant, etc.) and the text of the message.
type LLMMessage struct {
	Role LLMMessageRole
	Text string
}

// LLMProvider defines the interface that all LLM providers must implement.
// This allows for easy swapping between different LLM providers.
type LLMProvider interface {
	// GetResponse generates a response for the given question using the specified configuration.
	// Returns LLMResponse containing the generated text and metadata, or an error if the operation fails.
	GetResponse(messages []LLMMessage, config LLMRequestConfig) (LLMResponse, error)
}
