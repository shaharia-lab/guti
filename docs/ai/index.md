### AI Operations

[![Go Reference](https://pkg.go.dev/badge/github.com/shaharia-lab/guti/ai.svg)](https://pkg.go.dev/github.com/shaharia-lab/guti/ai)

The `ai` package provides a comprehensive interface for working with Language Learning Models (LLMs) and embedding models. It supports multiple providers (OpenAI, Anthropic), streaming responses, and various embedding models.

#### LLM Integration

Basic text generation with LLMs:

```go
import "github.com/shaharia-lab/guti/ai"

// Create OpenAI client and provider
client := ai.NewRealOpenAIClient("your-api-key",option.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}))
provider := ai.NewOpenAILLMProvider(ai.OpenAIProviderConfig{
    Client: client,
    Model:  "gpt-3.5-turbo", // Optional, defaults to gpt-3.5-turbo
})

// Create request with configuration
config := ai.NewRequestConfig(
    ai.WithMaxToken(2000),
    ai.WithTemperature(0.7),
)
request := ai.NewLLMRequest(config, provider)

// Generate response
response, err := request.Generate([]ai.LLMMessage{
    {Role: ai.SystemRole, Text: "You are a helpful assistant"},
    {Role: ai.UserRole, Text: "What is the capital of France?"},
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Response: %s\n", response.Text)
fmt.Printf("Tokens used: %d\n", response.TotalOutputToken)

```

#### Streaming Responses

Get realtime token-by-token responses:

```go
stream, err := request.GenerateStream(context.Background(), []ai.LLMMessage{
    {Role: ai.UserRole, Text: "Tell me a story"},
})
if err != nil {
    log.Fatal(err)
}

for response := range stream {
    if response.Error != nil {
        break
    }
    if response.Done {
        break
    }
    fmt.Print(response.Text)
}
```

#### Anthropic Integration

Use Claude models through Anthropic's API:

```go
// Create Anthropic client and provider
client := ai.NewRealAnthropicClient("your-api-key")
provider := ai.NewAnthropicLLMProvider(ai.AnthropicProviderConfig{
    Client: client,
    Model:  "claude-3-sonnet-20240229", // Optional, defaults to latest 3.5 Sonnet
})

request := ai.NewLLMRequest(config, provider)
```

#### Embedding Generation
[Embedding](/guti/ai/embedding.md)

#### Template Support

Create dynamic prompts using Go templates:

```go
template := &ai.LLMPromptTemplate{
    Template: "Hello {{.Name}}! Tell me about {{.Topic}}.",
    Data: map[string]interface{}{
        "Name":  "Alice",
        "Topic": "artificial intelligence",
    },
}

prompt, err := template.Parse()
if err != nil {
    log.Fatal(err)
}

response, err := request.Generate([]ai.LLMMessage{
    {Role: ai.UserRole, Text: prompt},
})
```

#### Configuration Options

| Option      | Default | Description                |
|-------------|---------|----------------------------|
| MaxToken    | 1000    | Maximum tokens to generate |
| TopP        | 0.9     | Nucleus sampling (0-1)     |
| Temperature | 0.7     | Output randomness (0-2)    |
| TopK        | 50      | Top-k sampling parameter   |

#### Custom Providers

Implement the provider interfaces to add support for additional services:

```go
// LLM Provider interface
type LLMProvider interface {
    GetResponse(messages []LLMMessage, config LLMRequestConfig) (LLMResponse, error)
    GetStreamingResponse(ctx context.Context, messages []LLMMessage, config LLMRequestConfig) (<-chan StreamingLLMResponse, error)
}

// OpenAI specific client interface
type OpenAIClient interface {
    CreateCompletion(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error)
    CreateStreamingCompletion(ctx context.Context, params openai.ChatCompletionNewParams) *ssestream.Stream[openai.ChatCompletionChunk]
}

// Embedding Provider interface
type EmbeddingProvider interface {
    GenerateEmbedding(ctx context.Context, input interface{}, model string) (*EmbeddingResponse, error)
}
```

The LLM providers now support dependency injection, allowing for better testability and configuration:

```go
// OpenAI with custom client
client := ai.NewRealOpenAIClient(
    "your-api-key",
    option.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
)
provider := ai.NewOpenAILLMProvider(ai.OpenAIProviderConfig{
    Client: client,
    Model:  "gpt-4",
})

// Anthropic with custom client
client := ai.NewRealAnthropicClient("your-api-key")
provider := ai.NewAnthropicLLMProvider(ai.AnthropicProviderConfig{
    Client: client,
    Model:  "claude-3-sonnet-20240229",
})
```