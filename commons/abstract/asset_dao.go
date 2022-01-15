package abstract

import "github.com/data-mill-cloud/mastro/commons/utils/conf"

// AssetDAOProvider ... The interface each dao must implement
type AssetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Upsert(asset *Asset) error
	GetById(id string) (*Asset, error)
	GetByName(id string) (*Asset, error)
	SearchAssetsByTags(tags []string, limit int, page int) (*PaginatedAssets, error)
	ListAllAssets(limit int, page int) (*PaginatedAssets, error)
	CloseConnection()
}
