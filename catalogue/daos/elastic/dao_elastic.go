package elastic

import (
	"sync"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/sources/elastic"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
)

var once sync.Once
var instance *dao

type dao struct {
	Connector *elastic.Connector
}

// GetSingleton ... get an instance of the dao backend
func GetSingleton() abstract.AssetDAOProvider {
	// once.do is lazy, we use it to return an instance of the DAO
	once.Do(func() {
		instance = &dao{}
	})
	return instance
}

// Init ... Initialize connection to elastic search and target index
func (dao *dao) Init(def *conf.DataSourceDefinition) {
	return
}

// Create ... Create asset on ES
func (dao *dao) Upsert(fs *abstract.Asset) error {
	return nil
}

// SearchAssetsByTags ... search for the provided tags
func (dao *dao) SearchAssetsByTags(tags []string) (*[]abstract.Asset, error) {
	return nil, nil
}

// ListAllFeatureSets ... Return all assets in index
func (dao *dao) ListAllAssets() (*[]abstract.Asset, error) {
	return nil, nil
}

// GetById ... Retrieve document by given id
func (dao *dao) GetById(id string) (*abstract.Asset, error) {
	return nil, nil
}

// GetByName ... Retrieve document by given id
func (dao *dao) GetByName(id string) (*abstract.Asset, error) {
	return nil, nil
}

// CloseConnection ... Terminates the connection to ES for the DAO
func (dao *dao) CloseConnection() {
	dao.Connector.CloseConnection()
}
