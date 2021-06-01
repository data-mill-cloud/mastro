package main

import (
	"log"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/commons/utils/date"
	"github.com/data-mill-cloud/mastro/commons/utils/errors"
)

// Service ... Service Interface listing implemented methods
type Service interface {
	Init(cfg *conf.Config) *errors.RestErr
	UpsertAssets(assets *[]abstract.Asset) (*[]abstract.Asset, *errors.RestErr)
	GetAssetByID(assetID string) (*abstract.Asset, *errors.RestErr)
	GetAssetByName(name string) (*abstract.Asset, *errors.RestErr)
	SearchAssetsByTags(tags []string) (*[]abstract.Asset, *errors.RestErr)
	ListAllAssets() (*[]abstract.Asset, *errors.RestErr)
}

// assetServiceType ... Service Type
type assetServiceType struct{}

// assetService ... Group all service methods in a kind FeatureSetServiceType implementing the FeatureSetService
var assetService Service = &assetServiceType{}

// selected dao for the featureSetService
var dao abstract.AssetDAOProvider

// Init ... initializes the service
func (s *assetServiceType) Init(cfg *conf.Config) *errors.RestErr {
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
func (s *assetServiceType) UpsertAssets(assets *[]abstract.Asset) (*[]abstract.Asset, *errors.RestErr) {
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
func (s *assetServiceType) GetAssetByID(assetID string) (*abstract.Asset, *errors.RestErr) {
	asset, err := dao.GetById(assetID)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return asset, nil
}

// GetAssetByName ... Retrieves an asset by its unique name
func (s *assetServiceType) GetAssetByName(name string) (*abstract.Asset, *errors.RestErr) {
	asset, err := dao.GetByName(name)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return asset, nil
}

func (s *assetServiceType) SearchAssetsByTags(tags []string) (*[]abstract.Asset, *errors.RestErr) {
	assets, err := dao.SearchAssetsByTags(tags)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return assets, nil
}

// ListAllAssets ... Retrieves all stored assets
func (s *assetServiceType) ListAllAssets() (*[]abstract.Asset, *errors.RestErr) {
	assets, err := dao.ListAllAssets()
	if err != nil {
		return nil, errors.GetInternalServerError(err.Error())
	}
	// n.b. - assets empty if collection is empty
	// better to return an error or an empty list?
	if assets == nil || len(*assets) == 0 {
		return nil, errors.GetNotFoundError("No assets in given collection")
	}
	return assets, nil
}
