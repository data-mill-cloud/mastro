package abstract

import "github.com/data-mill-cloud/mastro/commons/utils/conf"

// EmbeddingDAOProvider ... The interface each dao must implement
type EmbeddingDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(e *Embedding) error
	GetById(id string) (*Embedding, error)
	GetByName(name string, limit int, page int) (*PaginatedEmbeddings, error)
	SimilarToThis(id string, limit int, page int) (*PaginatedEmbeddings, error)
	CloseConnection()
}
