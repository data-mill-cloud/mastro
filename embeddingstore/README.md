# Mastro

## Embedding Store

An Embedding store is a service to store embeddings and compute similarity between them.
Embeddings are numerical vectors that can be extracted from unstructured data, such as text, images, audio and video files.

```go
// Embedding ... a named feature vector to be used for similarity search
type Embedding struct {
	Id         int64     `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	InsertedAt time.Time `json:"inserted_at,omitempty"`
	Vector     []float32 `json:"vector,omitempty"`
}
```

A data access object (DAO) for an embedding is defined as follows:

```go
// EmbeddingDAOProvider ... The interface each dao must implement
type EmbeddingDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(e *Embedding) error
	GetById(id string, partitions []string) (*Embedding, error)
	GetByName(name string, limit int, page int) (*PaginatedEmbeddings, error)
	SimilarToThis(id string, limit int, page int) (*PaginatedEmbeddings, error)
	CloseConnection()
}
```

The interface is then implemented for specific targets in the `embeddingstore/daos/*` packages.

## Service

A basic interface is dedined to retrieve embeddings:

```go
// Service
type EmbeddingStoreService interface {
	Init(cfg *conf.Config) *errors.RestErr
	CreateEmbedding(em abstract.Embedding) (*abstract.Embedding, *errors.RestErr)
	GetEmbeddingByID(emID string, partitions []string) (*abstract.Embedding, *errors.RestErr)
	GetEmbeddingByName(emName string, limit int, page int) (*abstract.PaginatedEmbeddings, *errors.RestErr)
	SimilarToThis(emId string, limit int, page int) (*abstract.PaginatedEmbeddings, *errors.RestErr)
}
```