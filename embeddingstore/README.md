# Mastro

## Embedding Store

An Embedding store is a service to store embeddings and compute similarity between them.
Embeddings are numerical vectors that can be extracted from unstructured data, such as text, images, audio and video files.

```go
// Embedding ... a named feature vector to be used for similarity search
type Embedding struct {
	Id         string    `json:"id,omitempty"`
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
	Upsert(e []Embedding) error
	GetById(id string) (*Embedding, error)
	GetByName(name string) ([]Embedding, error)
	SimilarToThis(vector []float32, k int) ([]Embedding, error)
	DeleteByName(name string) error
	DeleteByIds(ids ...string) error
	CloseConnection()
}
```

The interface is then implemented for specific targets in the `embeddingstore/daos/*` packages.

## Service

A basic interface is dedined to retrieve embeddings:

```go
// EmbeddingStoreService ... EmbeddingStoreService Interface listing service methods
type EmbeddingStoreService interface {
	Init(cfg *conf.Config) *resterrors.RestErr
	UpsertEmbeddings(embeddings []Embedding) *resterrors.RestErr
	GetEmbeddingByID(id string) (*Embedding, *resterrors.RestErr)
	GetEmbeddingByName(name string) ([]Embedding, *resterrors.RestErr)
	SimilarToThis(vector []float32, k int) ([]Embedding, *resterrors.RestErr)
	DeleteEmbeddingByName(name string) *resterrors.RestErr
	DeleteEmbeddingByIds(ids ...string) *resterrors.RestErr
}
```