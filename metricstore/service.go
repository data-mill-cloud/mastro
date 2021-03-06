package main

import (
	"log"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/commons/utils/date"
	"github.com/data-mill-cloud/mastro/commons/utils/errors"
)

// metricStoreServiceType ... Service Type
type metricStoreServiceType struct{}

// metricStoreService ... Group all service methods in a kind metricStoreServiceType implementing the metricStoreService
var metricStoreService abstract.MetricStoreService = &metricStoreServiceType{}

// selected dao for the metricStoreService
var dao abstract.MetricSetDAOProvider

func (s *metricStoreServiceType) Init(cfg *conf.Config) *errors.RestErr {
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

// CreateMetricSet ... Create a MetricSet entry
func (s *metricStoreServiceType) CreateMetricSet(ms abstract.MetricSet) (*abstract.MetricSet, *errors.RestErr) {
	if err := ms.Validate(); err != nil {
		return nil, errors.GetBadRequestError(err.Error())
	}
	// set insert time to current date, then insert using selected dao
	ms.InsertedAt = date.GetNow()
	err := dao.Create(&ms)
	if err != nil {
		return nil, errors.GetBadRequestError(err.Error())
	}
	// what should we actually return of the newly inserted object?
	return &ms, nil
}

// GetMetricSetByID ... Retrieves a MetricSet by ID
func (s *metricStoreServiceType) GetMetricSetByID(msID string) (*abstract.MetricSet, *errors.RestErr) {
	mset, err := dao.GetById(msID)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return mset, nil
}

// GetMetricSetByName ... Retrieves a MetricSet by Name
func (s *metricStoreServiceType) GetMetricSetByName(msName string, limit int, page int) (*abstract.Paginated[abstract.MetricSet], *errors.RestErr) {
	mset, err := dao.GetByName(msName, limit, page)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return mset, nil
}

// SearchMetricSetsByLabels ... Retrieve MetricSets by Labels
func (s *metricStoreServiceType) SearchMetricSetsByLabels(labels map[string]string, limit int, page int) (*abstract.Paginated[abstract.MetricSet], *errors.RestErr) {
	ms, err := dao.SearchMetricSetsByLabels(labels, limit, page)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return ms, nil
}

// ListAllMetricSets ... Retrieves all MetricSets
func (s *metricStoreServiceType) ListAllMetricSets(limit int, page int) (*abstract.Paginated[abstract.MetricSet], *errors.RestErr) {
	msets, err := dao.ListAllMetricSets(limit, page)
	if err != nil {
		return nil, errors.GetInternalServerError(err.Error())
	}
	if msets == nil || len(*msets.Data) == 0 {
		return nil, errors.GetNotFoundError("No metricsets in given collection")
	}
	return msets, nil
}

// Search ... Retrieves items by a search query
func (s *metricStoreServiceType) Search(query string, limit int, page int) (*abstract.Paginated[abstract.MetricSet], *errors.RestErr) {
	msets, err := dao.Search(query, limit, page)
	if err != nil {
		return nil, errors.GetInternalServerError(err.Error())
	}
	if msets == nil || len(*msets.Data) == 0 {
		return nil, errors.GetNotFoundError("No assets in given collection")
	}
	return msets, nil
}
