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
	CreateFeatureSet(fs abstract.FeatureSet) (*abstract.FeatureSet, *errors.RestErr)
	GetFeatureSetByID(fsID string) (*abstract.FeatureSet, *errors.RestErr)
	GetFeatureSetByName(fsName string) (*[]abstract.FeatureSet, *errors.RestErr)
	ListAllFeatureSets() (*[]abstract.FeatureSet, *errors.RestErr)
}

// featureSetServiceType ... Service Type
type featureSetServiceType struct{}

// FeatureSetService ... Group all service methods in a kind FeatureSetServiceType implementing the FeatureSetService
var featureSetService Service = &featureSetServiceType{}

// selected dao for the featureSetService
var dao abstract.FeatureSetDAOProvider

// Init ... Initializes the connector by validating the config and initializing the connection
func (s *featureSetServiceType) Init(cfg *conf.Config) *errors.RestErr {
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

// CreateFeatureSet ... Create a FeatureSet entry
func (s *featureSetServiceType) CreateFeatureSet(fs abstract.FeatureSet) (*abstract.FeatureSet, *errors.RestErr) {
	if err := fs.Validate(); err != nil {
		return nil, errors.GetBadRequestError(err.Error())
	}
	// set insert time to current date, then insert using selected dao
	fs.InsertedAt = date.GetNow()
	err := dao.Create(&fs)
	if err != nil {
		return nil, errors.GetBadRequestError(err.Error())
	}
	// what should we actually return of the newly inserted object?
	return &fs, nil
}

// GetFeatureSetByID ... Retrieves a FeatureSet
func (s *featureSetServiceType) GetFeatureSetByID(fsID string) (*abstract.FeatureSet, *errors.RestErr) {
	fset, err := dao.GetById(fsID)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return fset, nil
}

// GetFeatureSetByName ... Retrieves a FeatureSet
func (s *featureSetServiceType) GetFeatureSetByName(fsName string) (*[]abstract.FeatureSet, *errors.RestErr) {
	fset, err := dao.GetByName(fsName)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return fset, nil
}

// ListAllFeatureSets ... Retrieves all FeatureSets
func (s *featureSetServiceType) ListAllFeatureSets() (*[]abstract.FeatureSet, *errors.RestErr) {
	fsets, err := dao.ListAllFeatureSets()
	if err != nil {
		return nil, errors.GetInternalServerError(err.Error())
	}
	// n.b. - fsets empty if collection is empty
	// better to return an error or an empty list?
	if fsets == nil || len(*fsets) == 0 {
		return nil, errors.GetNotFoundError("No feature sets in given collection")
	}
	return fsets, nil
}
