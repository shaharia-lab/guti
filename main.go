package main

import (
	"context"
	"log"
	"net/http"

	"github.com/shaharia-lab/guti/ai"
)

func main() {
	// Create provider
	pgConfig := ai.PostgresStorageConfig{
		ConnectionString: "postgres://app:pass@localhost:5432/app?sslmode=disable",
		MaxDimension:     1536,
		SchemaName:       "vectors",
	}
	provider, err := ai.NewPostgresProvider(pgConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create and initialize storage
	ctx := context.Background()
	storage, err := ai.NewVectorStorage(ctx, provider)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	// Create a collection
	/*err = storage.CreateCollection(ctx, &ai.VectorCollectionConfig{
		Name:         "documents",
		Dimension:    384,
		IndexType:    ai.IndexTypeHNSW,
		DistanceType: ai.DistanceTypeCosine,
	})
	if err != nil {
		log.Fatal(err)
	}*/

	em := ai.NewEmbeddingService("http://localhost:8000", &http.Client{})
	embedd, _ := em.GenerateEmbedding(ctx, "Hello", ai.EmbeddingModelAllMiniLML6V2)

	/*err = storage.UpsertDocument(ctx, "documents", &ai.VectorDocument{
		ID:        "sdfadsf",
		Vector:    embedd.Data[0].Embedding,
		Content:   "asfdasdf",
		Metadata:  nil,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	})
	if err != nil {
		log.Fatal(err)
	}*/

	v, err := storage.SearchByVector(ctx, "documents", embedd.Data[0].Embedding, &ai.VectorSearchOptions{
		Limit:           1,
		Offset:          0,
		Filter:          nil,
		IncludeMetadata: true,
		IncludeVectors:  false,
	})
	if err != nil {
		return
	}

	log.Fatal(v[0].Document)
}
