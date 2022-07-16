package abstract

import (
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	resterrors "github.com/data-mill-cloud/mastro/commons/utils/errors"

	"errors"
	"strings"
	"time"
)

// Embedding ... a named feature vector to be used for similarity search
type Embedding struct {
	Id         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	InsertedAt time.Time `json:"inserted_at,omitempty"`
	Vector     []float32 `json:"vector,omitempty"`
}

// Validate ... validate an embedding
func (em *Embedding) Validate() error {
	// the name should not be empty or we may not be able to retrieve the embedding
	if len(strings.TrimSpace(em.Name)) == 0 {
		return errors.New("Embedding Name is undefined")
	}

	return nil
}

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
