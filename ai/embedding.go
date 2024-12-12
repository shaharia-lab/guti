// Package ai provides artificial intelligence utilities including embedding generation capabilities.
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type EmbeddingModel string

const (
	EmbeddingModelAllMiniLML6V2                     EmbeddingModel = "all-MiniLM-L6-v2"
	EmbeddingModelAllMpnetBaseV2                    EmbeddingModel = "all-mpnet-base-v2"
	EmbeddingModelParaphraseMultilingualMiniLML12V2 EmbeddingModel = "paraphrase-multilingual-MiniLM-L12-v2"
)

// EmbeddingProvider defines the interface for services that can generate embeddings from text.
type EmbeddingProvider interface {
	GenerateEmbedding(ctx context.Context, input interface{}, model string) (*EmbeddingResponse, error)
}

// EmbeddingObject represents a single embedding result.
type EmbeddingObject struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// Usage represents token usage information.
type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// EmbeddingResponse represents the complete response from the embedding API.
type EmbeddingResponse struct {
	Object string            `json:"object"`
	Data   []EmbeddingObject `json:"data"`
	Model  EmbeddingModel    `json:"model"`
	Usage  Usage             `json:"usage"`
}

// EmbeddingService implements the EmbeddingProvider interface.
type EmbeddingService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewEmbeddingService creates a new EmbeddingService with the specified base URL.
func NewEmbeddingService(baseURL string, httpClient *http.Client) *EmbeddingService {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &EmbeddingService{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}
}

type embeddingRequest struct {
	Input          interface{}    `json:"input"`
	Model          EmbeddingModel `json:"model"`
	EncodingFormat string         `json:"encoding_format"`
}

// GenerateEmbedding implements the EmbeddingProvider interface.
func (s *EmbeddingService) GenerateEmbedding(ctx context.Context, input interface{}, model EmbeddingModel) (*EmbeddingResponse, error) {
	reqBody := embeddingRequest{
		Input:          input,
		Model:          model,
		EncodingFormat: "float",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.BaseURL+"/v1/embeddings", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var embResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &embResp, nil
}
