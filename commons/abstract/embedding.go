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
	Id         int64     `json:"id,omitempty"`
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
	Create(e *Embedding) error
	GetById(id string, partitions []string) (*Embedding, error)
	GetByName(name string, limit int, page int) (*Paginated[Embedding], error)
	SimilarToThisId(id string, limit int, page int) (*Paginated[Embedding], error)
	CloseConnection()
}

// EmbeddingStoreService ... EmbeddingStoreService Interface listing service methods
type EmbeddingStoreService interface {
	Init(cfg *conf.Config) *resterrors.RestErr
	CreateEmbedding(em Embedding) (*Embedding, *resterrors.RestErr)
	GetEmbeddingByID(emID string, partitions []string) (*Embedding, *resterrors.RestErr)
	GetEmbeddingByName(emName string, limit int, page int) (*Paginated[Embedding], *resterrors.RestErr)
	SimilarToThisId(emId string, limit int, page int) (*Paginated[Embedding], *resterrors.RestErr)
}
