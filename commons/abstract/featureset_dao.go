package abstract

import "github.com/data-mill-cloud/mastro/commons/utils/conf"

// FeatureSetDAOProvider ... The interface each dao must implement
type FeatureSetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(fs *FeatureSet) error
	GetById(id string) (*FeatureSet, error)
	GetByName(name string) (*[]FeatureSet, error)
	ListAllFeatureSets() (*[]FeatureSet, error)
	CloseConnection()
}
