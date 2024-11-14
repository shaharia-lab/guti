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
go get go get github.com/shaharia-lab/guti
```

And start using like -

```golang
import (
    "go get github.com/shaharia-lab/guti"
)

guti.ContainsAll()
```

### AI Operations

The `ai` package provides a flexible interface for interacting with various Language Learning Models (LLMs). Currently supports OpenAI's GPT models with an extensible interface for other providers.

#### Basic Usage

```go
import (
    "github.com/shaharia-lab/guti/ai"
)

// Create an OpenAI provider
provider := ai.NewOpenAILLMProvider(ai.OpenAIProviderConfig{
    APIKey: "your-api-key",
    Model:  "gpt-3.5-turbo", // Optional, defaults to gpt-3.5-turbo
})

// Create a request with default configuration
request := ai.NewLLMRequest(ai.NewRequestConfig())

// Generate a response
response, err := request.Generate("What is the capital of France?", provider)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Response: %s\n", response.Text)
fmt.Printf("Input tokens: %d\n", response.TotalInputToken)
fmt.Printf("Output tokens: %d\n", response.TotalOutputToken)
fmt.Printf("Completion time: %.2f seconds\n", response.CompletionTime)
```

#### Custom Configuration

You can customize the LLM request configuration using the functional options pattern:

```go
// Use specific configuration options
config := ai.NewRequestConfig(
    ai.WithMaxToken(2000),
    ai.WithTemperature(0.8),
    ai.WithTopP(0.95),
    ai.WithTopK(100),
)

request := ai.NewLLMRequest(config)
```

#### Using Templates

The package also supports templated prompts:

```go
template := &ai.LLMPromptTemplate{
    Template: "Hello {{.Name}}! Please tell me about {{.Topic}}.",
    Data: map[string]interface{}{
        "Name":  "Alice",
        "Topic": "artificial intelligence",
    },
}

prompt, err := template.Parse()
if err != nil {
    log.Fatal(err)
}

response, err := request.Generate(prompt, provider)
```

#### Configuration Options

| Option      | Default | Description                          |
|-------------|---------|--------------------------------------|
| MaxToken    | 1000    | Maximum number of tokens to generate |
| TopP        | 0.9     | Nucleus sampling parameter (0-1)     |
| Temperature | 0.7     | Randomness in output (0-2)           |
| TopK        | 50      | Top-k sampling parameter             |

#### Error Handling

The package provides structured error handling:

```go
response, err := request.Generate(prompt, provider)
if err != nil {
    if llmErr, ok := err.(*ai.LLMError); ok {
        fmt.Printf("LLM Error %d: %s\n", llmErr.Code, llmErr.Message)
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

#### Custom Providers

You can implement the `LLMProvider` interface to add support for additional LLM providers:

```go
type LLMProvider interface {
    GetResponse(question string, config LLMRequestConfig) (LLMResponse, error)
}
```

## Documentation

Full documentation is available on [pkg.go.dev/github.com/shaharia-lab/guti](https://pkg.go.dev/github.com/shaharia-lab/guti#section-documentation)

### 📝 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/shaharia-lab/guti/blob/master/LICENSE) file for details.