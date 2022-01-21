package abstract

import "github.com/data-mill-cloud/mastro/commons/utils/conf"

// FeatureSetDAOProvider ... The interface each dao must implement
type FeatureSetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(fs *FeatureSet) error
	GetById(id string) (*FeatureSet, error)
	GetByName(name string, limit int, page int) (*PaginatedFeatureSets, error)
	ListAllFeatureSets(limit int, page int) (*PaginatedFeatureSets, error)
	Search(query string, limit int, page int) (*PaginatedFeatureSets, error)
	CloseConnection()
}
