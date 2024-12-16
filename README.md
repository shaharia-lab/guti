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

The `ai` package provides a flexible interface for interacting with various Language Learning Models (LLMs), supporting both regular and streaming responses.

#### Basic Usage

```go
import "github.com/shaharia-lab/guti/ai"

// Create provider and request
provider := ai.NewOpenAILLMProvider(ai.OpenAIProviderConfig{
    APIKey: "your-api-key",
    Model:  "gpt-3.5-turbo", // Optional, defaults to gpt-3.5-turbo
})
request := ai.NewLLMRequest(ai.NewRequestConfig(), provider)

// Generate response
response, err := request.Generate([]LLMMessage{
    {Role: "user", Text: "What is the capital of France?"},
})
```

#### Streaming Responses

```go
stream, err := request.GenerateStream(context.Background(), []LLMMessage{
    {Role: "user", Text: "Tell me a story"},
})
if err != nil {
    log.Fatal(err)
}

for response := range stream {
    if response.Error != nil {
        break
    }
    fmt.Print(response.Text)
}
```

#### Custom Configuration

```go
config := ai.NewRequestConfig(
    ai.WithMaxToken(2000),
    ai.WithTemperature(0.8),
    ai.WithTopP(0.95),
)
request := ai.NewLLMRequest(config)
```

#### Using Templates

```go
template := &ai.LLMPromptTemplate{
    Template: "Hello {{.Name}}! Tell me about {{.Topic}}.",
    Data: map[string]interface{}{
        "Name":  "Alice",
        "Topic": "AI",
    },
}

prompt, err := template.Parse()
response, err := request.Generate(prompt, provider)
```

#### Configuration Options

| Option      | Default | Description                          |
|-------------|---------|--------------------------------------|
| MaxToken    | 1000    | Maximum tokens to generate           |
| TopP        | 0.9     | Nucleus sampling parameter (0-1)     |
| Temperature | 0.7     | Randomness in output (0-2)           |
| TopK        | 50      | Top-k sampling parameter             |

#### Generate Embeddings

```go
provider := ai.NewEmbeddingService("http://localhost:8000", nil)
embedding, err := provider.GenerateEmbedding(
    context.Background(),
    "Hello world",
    ai.EmbeddingModelAllMiniLML6V2,
)
```

#### Custom Providers

You can implement the `LLMProvider` interface to add support for additional LLM providers:

```go
type LLMProvider interface {
    GetResponse(messages []LLMMessage, config LLMRequestConfig) (LLMResponse, error)
}
```

#### Generate Embedding Vector

You can generate embeddings using the provider-based approach:

```go
import (
    "github.com/shaharia-lab/guti/ai"
)

// Create an embedding provider
provider := ai.NewLocalEmbeddingProvider(ai.LocalProviderConfig{
    BaseURL: "http://localhost:8000",
    Client:  &http.Client{},
})

// Generate embedding
embedding, err := provider.GenerateEmbedding(context.Background(), "Hello world", ai.EmbeddingModelAllMiniLML6V2)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Embedding vector: %+v\n", embedding)
```

The library supports multiple embedding providers. You can implement the `EmbeddingProvider` interface to add support for additional providers:

```go
type EmbeddingProvider interface {
    GenerateEmbedding(ctx context.Context, text string, model EmbeddingModel) ([]float32, error)
}
```

## Documentation

Full documentation is available on [pkg.go.dev/github.com/shaharia-lab/guti](https://pkg.go.dev/github.com/shaharia-lab/guti#section-documentation)

### 📝 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/shaharia-lab/guti/blob/master/LICENSE) file for details.