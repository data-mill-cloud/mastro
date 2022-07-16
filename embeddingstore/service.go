package main

import (
	"log"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/commons/utils/date"
	"github.com/data-mill-cloud/mastro/commons/utils/errors"
)

// embeddingServiceType ... Service Type
type embeddingServiceType struct{}

// embeddingService ... Group all service methods in a kind embeddingServiceType implementing the embeddingService
var embeddingService abstract.EmbeddingStoreService = &embeddingServiceType{}

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

// UpsertEmbeddings ... Create embeddings
func (s *embeddingServiceType) UpsertEmbeddings(embeddings []abstract.Embedding) *errors.RestErr {
	now := date.GetNow()
	for i, em := range embeddings {
		if err := em.Validate(); err != nil {
			return errors.GetBadRequestError(err.Error())
		}
		// set insert time to current date, then insert using selected dao
		em.InsertedAt = now
		embeddings[i] = em
	}

	if err := dao.Upsert(embeddings); err != nil {
		return errors.GetBadRequestError(err.Error())
	}
	return nil
}

// GetEmbeddingByID ... Retrieves an embedding
func (s *embeddingServiceType) GetEmbeddingByID(id string) (*abstract.Embedding, *errors.RestErr) {
	em, err := dao.GetById(id)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return em, nil
}

// GetEmbeddingByName ... Retrieves an embedding
func (s *embeddingServiceType) GetEmbeddingByName(emName string) ([]abstract.Embedding, *errors.RestErr) {
	em, err := dao.GetByName(emName)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return em, nil
}

// SimilarToThis ... Retrieves embeddings similar to the one provided
func (s *embeddingServiceType) SimilarToThis(vector []float32, k int) ([]abstract.Embedding, *errors.RestErr) {
	em, err := dao.SimilarToThis(vector, k)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return em, nil
}

func (s *embeddingServiceType) DeleteEmbeddingByName(name string) *errors.RestErr {
	if err := dao.DeleteByName(name); err != nil {
		return errors.GetBadRequestError(err.Error())
	}
	return nil
}

func (s *embeddingServiceType) DeleteEmbeddingByIds(ids ...string) *errors.RestErr {
	if err := dao.DeleteByIds(ids...); err != nil {
		return errors.GetBadRequestError(err.Error())
	}
	return nil
}
