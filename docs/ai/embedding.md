# Embedding Generation

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