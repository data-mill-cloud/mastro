package main

import (
	"log"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/commons/utils/date"
	"github.com/data-mill-cloud/mastro/commons/utils/errors"
)

// catalogueServiceType ... Service Type
type catalogueServiceType struct{}

var catalogueService abstract.CatalogueService = &catalogueServiceType{}
var dao abstract.AssetDAOProvider

// Init ... initializes the service
func (s *catalogueServiceType) Init(cfg *conf.Config) *errors.RestErr {
	// select target DAO based on used connector
	// set a connector to the selected backend here
	var err error
	// select dao using mapping function in same package
	dao, err = selectDao(cfg)
	if err != nil {
		log.Panicln(err)
	}
	dao.Init(&cfg.DataSourceDefinition)
	return nil
}

// UpsertAsset ... Adds and asset description
func (s *catalogueServiceType) UpsertAssets(assets *[]abstract.Asset) (*[]abstract.Asset, *errors.RestErr) {
	for _, a := range *assets {
		if err := a.Validate(); err != nil {
			return nil, errors.GetBadRequestError(err.Error())
		}
		// add last discovered date
		a.LastDiscoveredAt = date.GetNow()
		err := dao.Upsert(&a)

		if err != nil {
			return nil, errors.GetBadRequestError(err.Error())
		}
	}

	// what should we actually return of the newly inserted object?
	return assets, nil
}

// GetAssetById ... Retrieves an asset by its unique id
func (s *catalogueServiceType) GetAssetByID(assetID string) (*abstract.Asset, *errors.RestErr) {
	asset, err := dao.GetById(assetID)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return asset, nil
}

// GetAssetByName ... Retrieves an asset by its unique name
func (s *catalogueServiceType) GetAssetByName(name string) (*abstract.Asset, *errors.RestErr) {
	asset, err := dao.GetByName(name)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return asset, nil
}

func (s *catalogueServiceType) SearchAssetsByTags(tags []string, limit int, page int) (*abstract.Paginated[abstract.Asset], *errors.RestErr) {
	assets, err := dao.SearchAssetsByTags(tags, limit, page)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return assets, nil
}

// ListAllAssets ... Retrieves all stored assets
func (s *catalogueServiceType) ListAllAssets(limit int, page int) (*abstract.Paginated[abstract.Asset], *errors.RestErr) {
	assets, err := dao.ListAllAssets(limit, page)
	if err != nil {
		return nil, errors.GetInternalServerError(err.Error())
	}
	// n.b. - assets empty if collection is empty
	// better to return an error or an empty list?
	if assets == nil || len(*assets.Data) == 0 {
		return nil, errors.GetNotFoundError("No assets in given collection")
	}
	return assets, nil
}

// Search ... Retrieves items by a search query
func (s *catalogueServiceType) Search(query string, limit int, page int) (*abstract.Paginated[abstract.Asset], *errors.RestErr) {
	assets, err := dao.Search(query, limit, page)
	if err != nil {
		return nil, errors.GetInternalServerError(err.Error())
	}
	// n.b. - assets empty if collection is empty
	// better to return an error or an empty list?
	if assets == nil || len(*assets.Data) == 0 {
		return nil, errors.GetNotFoundError("No assets in given collection")
	}
	return assets, nil
}
