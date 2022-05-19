package milvus

import (
	"context"
	"fmt"
	"strconv"
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
		return fmt.Errorf("failed to insert data: %v", err.Error())
	}
	return nil
}

func (dao *dao) GetById(id string, partitions []string) (*abstract.Embedding, error) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s of type %T", id, id)
	}

	queryResult, err := dao.Connector.Client.QueryByPks(
		context.Background(),
		dao.Connector.Collection,
		partitions,
		entity.NewColumnInt64("id", []int64{intId}),
		[]string{"id", "name"},
	)
	if err != nil {
		return nil, fmt.Errorf("fail to query collection: %v", err.Error())
	}

	if len(queryResult) == 0 {
		return nil, fmt.Errorf("query result for %s is empty", id)
	}

	return &abstract.Embedding{
		Id:     queryResult[0].FieldData().GetFieldId(),
		Name:   queryResult[1].Name(),
		Vector: queryResult[2].FieldData().GetScalars().GetFloatData().GetData(),
	}, nil
}

func (dao *dao) GetByName(name string, limit int, page int) (*abstract.Paginated[abstract.Embedding], error) {
	return nil, nil
}

func (dao *dao) SimilarToThisId(id string, limit int, page int) (*abstract.Paginated[abstract.Embedding], error) {
	return nil, nil
}
