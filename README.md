<h1 align="center">Guti</h1>
<p align="center">A Utility Library for Golang Developer</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/shaharia-lab/guti"><img src="https://pkg.go.dev/badge/github.com/shaharia-lab/guti.svg" height="20"/></a>
</p><br/><br/>

<p align="center">
  <a href="https://github.com/shaharia-lab/guti/actions/workflows/CI.yaml"><img src="https://github.com/shaharia-lab/guti/actions/workflows/CI.yaml/badge.svg" height="20"/></a>
  <a href="https://codecov.io/gh/shaharia-lab/guti"><img src="https://codecov.io/gh/shaharia-lab/guti/branch/master/graph/badge.svg?token=NKTKQ45HDN" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shaharia-lab_guti"><img src="https://sonarcloud.io/api/project_badges/measure?project=shaharia-lab_guti&metric=reliability_rating" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shaharia-lab_guti"><img src="https://sonarcloud.io/api/project_badges/measure?project=shaharia-lab_guti&metric=vulnerabilities" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shaharia-lab_guti"><img src="https://sonarcloud.io/api/project_badges/measure?project=shaharia-lab_guti&metric=security_rating" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shaharia-lab_guti"><img src="https://sonarcloud.io/api/project_badges/measure?project=shaharia-lab_guti&metric=sqale_rating" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shaharia-lab_guti"><img src="https://sonarcloud.io/api/project_badges/measure?project=shaharia-lab_guti&metric=code_smells" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shaharia-lab_guti"><img src="https://sonarcloud.io/api/project_badges/measure?project=shaharia-lab_guti&metric=ncloc" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shaharia-lab_guti"><img src="https://sonarcloud.io/api/project_badges/measure?project=shaharia-lab_guti&metric=alert_status" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shaharia-lab_guti"><img src="https://sonarcloud.io/api/project_badges/measure?project=shaharia-lab_guti&metric=duplicated_lines_density" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shaharia-lab_guti"><img src="https://sonarcloud.io/api/project_badges/measure?project=shaharia-lab_guti&metric=bugs" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shaharia-lab_guti"><img src="https://sonarcloud.io/api/project_badges/measure?project=shaharia-lab_guti&metric=sqale_index" height="20"/></a>
</p><br/><br/>

## Usage

```shell
go get github.com/shaharia-lab/guti
```

And start using like -

```golang
import (
    "github.com/shaharia-lab/guti"
)

guti.ContainsAll()
```

### AI Operations

The `ai` package provides a comprehensive interface for working with Language Learning Models (LLMs) and embedding models. It supports multiple providers (OpenAI, Anthropic), streaming responses, and various embedding models.

#### LLM Integration

Basic text generation with LLMs:

```go
import "github.com/shaharia-lab/guti/ai"

// Create an OpenAI provider
provider := ai.NewOpenAILLMProvider(ai.OpenAIProviderConfig{
    APIKey: "your-api-key",
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

Generate vector embeddings for text:

```go
provider := ai.NewEmbeddingService("http://api.example.com", nil)

embedding, err := provider.GenerateEmbedding(
    context.Background(),
    "Hello world",
    ai.EmbeddingModelAllMiniLML6V2,
)
if err != nil {
    log.Fatal(err)
}
```

Supported embedding models:
- `EmbeddingModelAllMiniLML6V2`: Lightweight, general-purpose model
- `EmbeddingModelAllMpnetBaseV2`: Higher quality, more compute intensive
- `EmbeddingModelParaphraseMultilingualMiniLML12V2`: Optimized for multilingual text

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
type LLMProvider interface {
    GetResponse(messages []LLMMessage, config LLMRequestConfig) (LLMResponse, error)
    GetStreamingResponse(ctx context.Context, messages []LLMMessage, config LLMRequestConfig) (<-chan StreamingLLMResponse, error)
}

type EmbeddingProvider interface {
    GenerateEmbedding(ctx context.Context, input interface{}, model string) (*EmbeddingResponse, error)
}
```

## Documentation

Full documentation is available on [pkg.go.dev/github.com/shaharia-lab/guti](https://pkg.go.dev/github.com/shaharia-lab/guti#section-documentation)

### 📝 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/shaharia-lab/guti/blob/master/LICENSE) file for details.