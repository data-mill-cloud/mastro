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
	CreateEmbedding(em abstract.Embedding) (*abstract.Embedding, *errors.RestErr)
	GetEmbeddingByID(emID string) (*abstract.Embedding, *errors.RestErr)
	GetEmbeddingByName(emName string, limit int, page int) (*abstract.PaginatedEmbeddings, *errors.RestErr)
	SimilarToThis(emId string, limit int, page int) (*abstract.PaginatedEmbeddings, *errors.RestErr)
}

// embeddingServiceType ... Service Type
type embeddingServiceType struct{}

// embeddingService ... Group all service methods in a kind embeddingServiceType implementing the embeddingService
var embeddingService Service = &embeddingServiceType{}

// selected dao for the embeddingService
var dao abstract.EmbeddingDAOProvider

// Init ... Initializes the connector by validating the config and initializing the connection
func (s *embeddingServiceType) Init(cfg *conf.Config) *errors.RestErr {
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

// CreateEmbedding ... Create an embedding entry
func (s *embeddingServiceType) CreateEmbedding(em abstract.Embedding) (*abstract.Embedding, *errors.RestErr) {
	if err := em.Validate(); err != nil {
		return nil, errors.GetBadRequestError(err.Error())
	}
	// set insert time to current date, then insert using selected dao
	em.InsertedAt = date.GetNow()
	err := dao.Create(&em)
	if err != nil {
		return nil, errors.GetBadRequestError(err.Error())
	}
	// what should we actually return of the newly inserted object?
	return &em, nil
}

// GetEmbeddingByID ... Retrieves an embedding
func (s *embeddingServiceType) GetEmbeddingByID(emID string) (*abstract.Embedding, *errors.RestErr) {
	em, err := dao.GetById(emID)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return em, nil
}

// GetEmbeddingByName ... Retrieves an embedding
func (s *embeddingServiceType) GetEmbeddingByName(emName string, limit int, page int) (*abstract.PaginatedEmbeddings, *errors.RestErr) {
	em, err := dao.GetByName(emName, limit, page)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return em, nil
}

// GetEmbeddingByName ... Retrieves an embedding
func (s *embeddingServiceType) SimilarToThis(emId string, limit int, page int) (*abstract.PaginatedEmbeddings, *errors.RestErr) {
	em, err := dao.SimilarToThis(emId, limit, page)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return em, nil
}
