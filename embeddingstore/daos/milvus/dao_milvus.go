package milvus

import (
	"context"
	"log"
	"sync"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/sources/milvus"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type dao struct {
	Connector *milvus.Connector
}

// both init and sync.Once are thread-safe
// but only sync.Once is lazy
var once sync.Once
var instance *dao

// GetSingleton ... lazy singleton on DAO
func GetSingleton() abstract.EmbeddingDAOProvider {
	// once.do is lazy, we use it to return an instance of the DAO
	once.Do(func() {
		instance = &dao{}
	})
	return instance
}

func (dao *dao) Init(def *conf.DataSourceDefinition) {
	// create milvus connector
	dao.Connector = milvus.NewMilvusConnector()
	// validate data source definition
	if err := dao.Connector.ValidateDataSourceDefinition(def); err != nil {
		panic(err)
	}
	// init milvus connector
	dao.Connector.InitConnection(def)
}

func (dao *dao) CloseConnection() {
	dao.Connector.CloseConnection()
}

func (dao *dao) convertDtoToDao(eb *abstract.Embedding) []entity.Column {
	return []entity.Column{
		entity.NewColumnString("name", []string{
			eb.Name,
		}),
		entity.NewColumnFloatVector(dao.Connector.DenseVectorFieldName, 2, [][]float32{
			eb.Vector,
		}),
	}
}

func (dao *dao) Create(embedding *abstract.Embedding) error {

	columns := dao.convertDtoToDao(embedding)
	if _, err := dao.Connector.Client.Insert(context.Background(), dao.Connector.Collection,
		"", // partitionName // TODO: fixme
		columns...,
	); err != nil {
		log.Fatal("failed to insert data:", err.Error())
	}
	return nil
}

func (dao *dao) GetById(id string) (*abstract.Embedding, error) {
	return nil, nil
}

func (dao *dao) GetByName(name string, limit int, page int) (*abstract.PaginatedEmbeddings, error) {
	return nil, nil
}

func (dao *dao) SimilarToThis(id string, limit int, page int) (*abstract.PaginatedEmbeddings, error) {
	return nil, nil
}
