# Vector Storage

The vector storage system provides a flexible and efficient way to store, manage, and search vector embeddings with associated metadata. It is particularly useful for applications involving semantic search, recommendation systems, and similarity matching.

## Key Features

- Collection-based organization of vector data
- Multiple index types for performance optimization
- Configurable distance metrics
- Support for document metadata
- Batch operations for efficient data management
- Rich search capabilities with filtering options

## Installation

First, ensure you have the required dependencies:

```go
go get github.com/shaharia-lab/guti
```

If using PostgreSQL storage (recommended), you'll need:
- PostgreSQL 11 or later
- pgvector extension installed

## Quick Start

Here's a basic example of using vector storage:

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/shaharia-lab/guti/ai"
)

func main() {
    // Create PostgreSQL provider
    config := ai.PostgresStorageConfig{
        ConnectionString: "postgres://user:pass@localhost:5432/dbname",
        MaxDimension:    384,
        SchemaName:      "vectors", // Optional, defaults to "public"
    }
    
    provider, err := ai.NewPostgresProvider(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Initialize vector storage
    storage, err := ai.NewVectorStorage(context.Background(), provider)
    if err != nil {
        log.Fatal(err)
    }
    defer storage.Close()
    
    // Create a collection
    err = storage.CreateCollection(context.Background(), &ai.VectorCollectionConfig{
        Name:         "documents",
        Dimension:    384,
        IndexType:    ai.IndexTypeHNSW,
        DistanceType: ai.DistanceTypeCosine,
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

## Working with Documents

### Storing Documents

```go
// Create a document with vector embedding and metadata
doc := &ai.VectorDocument{
    ID:      "doc1",
    Vector:  []float32{0.1, 0.2, 0.3, ...}, // 384-dimensional vector
    Content: "Document text content",
    Metadata: map[string]interface{}{
        "category": "technology",
        "author":   "John Doe",
        "tags":     []string{"ai", "machine-learning"},
    },
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}

// Store the document
err = storage.UpsertDocument(context.Background(), "documents", doc)
if err != nil {
    log.Fatal(err)
}

// Store multiple documents in batch
docs := []*ai.VectorDocument{doc1, doc2, doc3}
err = storage.UpsertDocuments(context.Background(), "documents", docs)
if err != nil {
    log.Fatal(err)
}
```

### Searching Documents

```go
// Search by vector similarity
searchVector := []float32{0.15, 0.25, 0.35, ...}
results, err := storage.SearchByVector(context.Background(), "documents", searchVector, &ai.VectorSearchOptions{
    Limit:           10,
    Offset:          0,
    IncludeMetadata: true,
    Filter: map[string]interface{}{
        "category": "technology",
    },
})
if err != nil {
    log.Fatal(err)
}

// Process results
for _, result := range results {
    fmt.Printf("Document ID: %s, Score: %f\n", result.Document.ID, result.Score)
    fmt.Printf("Content: %s\n", result.Document.Content)
    fmt.Printf("Category: %s\n", result.Document.Metadata["category"])
}

// Search by existing document ID
similarDocs, err := storage.SearchByID(context.Background(), "documents", "doc1", &ai.VectorSearchOptions{
    Limit: 5,
})
```

## Collection Configuration

### Index Types

The system supports different index types for optimizing search performance:

- `IndexTypeFlat`: Brute force search, best for small collections
- `IndexTypeIVFFlat`: IVF-based index, good balance of speed and accuracy
- `IndexTypeHNSW`: Hierarchical Navigable Small World graph, fastest search but more memory intensive

### Distance Metrics

Available distance metrics for similarity calculation:

- `DistanceTypeCosine`: Cosine similarity, best for normalized vectors
- `DistanceTypeEuclidean`: Euclidean distance, good for absolute distances
- `DistanceTypeDotProduct`: Dot product similarity, useful for specialized cases

Example configuration:

```go
config := &ai.VectorCollectionConfig{
    Name:      "semantic_search",
    Dimension: 384,
    IndexType: ai.IndexTypeHNSW,
    DistanceType: ai.DistanceTypeCosine,
    CustomFields: map[string]ai.VectorFieldConfig{
        "category": {
            Type:     "string",
            Required: true,
            Indexed:  true,
        },
        "rating": {
            Type:     "float",
            Required: false,
            Indexed:  true,
        },
    },
}
```

## Common Use Cases

1. **Semantic Search**
    - Store document embeddings generated from text
    - Search for semantically similar documents
    - Use metadata filtering for refined results

2. **Recommendation Systems**
    - Store user/item embeddings
    - Find similar items or users
    - Combine with metadata for personalized recommendations

3. **Image Similarity**
    - Store image feature vectors
    - Find visually similar images
    - Use metadata for filtering by attributes

4. **Duplicate Detection**
    - Store document fingerprints as vectors
    - Search for near-duplicates
    - Use distance thresholds for similarity detection

## Best Practices

1. **Vector Dimensionality**
    - Choose dimension based on your embedding model
    - Maintain consistent dimensionality within collections
    - Consider memory/storage implications

2. **Index Selection**
    - Use `IndexTypeHNSW` for large-scale, high-performance search
    - Use `IndexTypeFlat` for small collections or when accuracy is critical
    - Use `IndexTypeIVFFlat` for a balance of performance and resource usage

3. **Batch Operations**
    - Use `UpsertDocuments` for bulk insertions
    - Consider batch size based on memory constraints
    - Use transactions for data consistency

4. **Search Optimization**
    - Use metadata filters to reduce search space
    - Set appropriate limit/offset for pagination
    - Choose suitable distance metrics for your use case

## Error Handling

The system provides detailed error types for different scenarios:

```go
if err != nil {
    switch {
    case errors.Is(err, ai.ErrDocumentNotFound):
        // Handle missing document
    case errors.Is(err, ai.ErrCollectionNotFound):
        // Handle missing collection
    case errors.Is(err, ai.ErrInvalidDimension):
        // Handle dimension mismatch
    default:
        // Handle other errors
    }
}
```