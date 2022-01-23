package mongo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/sources/mongo"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	paginate "github.com/gobeam/mongo-go-pagination"
	mongodriver "go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/x/bsonx"
	"gopkg.in/mgo.v2/bson"
)

// metricSetMongoDao ... DAO for the MetricSet in Mongo
type metricSetMongoDao struct {
	//ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string            `bson:"name,omitempty"`
	InsertedAt  time.Time         `bson:"inserted-at,omitempty"`
	Version     string            `bson:"version,omitempty"`
	Description string            `bson:"description,omitempty"`
	Labels      map[string]string `bson:"labels,omitempty"`
	Metrics     []metricMongoDao  `bson:"metrics,omitempty"`
}

// metricMongoDao ... a named variable with a data type
type metricMongoDao struct {
	DeequMetric *deequMetricMongoDao `bson:"deequ,omitempty"`
}

type deequMetricMongoDao struct {
	ResultKey       deequResultKeyMongoDao       `bson:"resultKey,omitempty"`
	AnalyzerContext deequAnalyzerContextMongoDao `bson:"analyzerContext,omitempty"`
}

type deequResultKeyMongoDao struct {
	DataSetDate int64             `bson:"dataSetDate,omitempty"`
	Tags        map[string]string `bson:"tags,omitempty"`
}

type deequAnalyzerContextMongoDao struct {
	MetricMap []deequMetricInstanceMongoDao `bson:"metricMap,omitempty"`
}

type deequMetricInstanceMongoDao struct {
	Analyzer deequAnalyzerMongoDao    `bson:"analyzer,omitempty"`
	Metric   deequMetricValueMongoDao `bson:"metric,omitempty"`
}

type deequAnalyzerMongoDao struct {
	AnalyzerName string `bson:"analyzerName,omitempty"`
	Column       string `bson:"column,omitempty"`
}

type deequMetricValueMongoDao struct {
	MetricName string  `bson:"metricName,omitempty"`
	Entity     string  `bson:"entity,omitempty"`
	Instance   string  `bson:"instance,omitempty"`
	Name       string  `bson:"name,omitempty"`
	Value      float64 `bson:"value,omitempty"`
}

// ------------------------------
func convertMetricInstanceDTOtoDAO(mi *abstract.DeequMetricInstance) deequMetricInstanceMongoDao {
	return deequMetricInstanceMongoDao{
		Analyzer: deequAnalyzerMongoDao{
			AnalyzerName: mi.Analyzer.AnalyzerName,
			Column:       mi.Analyzer.Column,
		},
		Metric: deequMetricValueMongoDao{
			MetricName: mi.Metric.MetricName,
			Entity:     mi.Metric.Entity,
			Instance:   mi.Metric.Instance,
			Name:       mi.Metric.Name,
			Value:      mi.Metric.Value,
		},
	}
}

func convertAllMetricInstancesDTOtoDAO(inputInstances []abstract.DeequMetricInstance) []deequMetricInstanceMongoDao {
	var instances []deequMetricInstanceMongoDao
	for _, element := range inputInstances {
		instances = append(instances, convertMetricInstanceDTOtoDAO(&element))
	}
	return instances
}

func convertMetricDTOtoDAO(m *abstract.Metric) *metricMongoDao {
	result := &metricMongoDao{}
	if m.DeequMetric != nil {
		result.DeequMetric = &deequMetricMongoDao{
			ResultKey: deequResultKeyMongoDao{
				DataSetDate: m.DeequMetric.ResultKey.DataSetDate,
				Tags:        m.DeequMetric.ResultKey.Tags,
			},
			AnalyzerContext: deequAnalyzerContextMongoDao{
				MetricMap: convertAllMetricInstancesDTOtoDAO(m.DeequMetric.AnalyzerContext.MetricMap),
			},
		}
	}
	return result
}

func convertAllMetricsDTOtoDAO(mm []abstract.Metric) []metricMongoDao {
	var metrics []metricMongoDao
	for _, element := range mm {
		metrics = append(metrics, *convertMetricDTOtoDAO(&element))
	}
	return metrics
}

func convertMetricSetDTOtoDAO(ms *abstract.MetricSet) *metricSetMongoDao {
	return &metricSetMongoDao{
		Name:        ms.Name,
		InsertedAt:  ms.InsertedAt,
		Version:     ms.Version,
		Description: ms.Description,
		Labels:      ms.Labels,
		Metrics:     convertAllMetricsDTOtoDAO(ms.Metrics),
	}
}

func convertMetricInstanceDAOToDTO(mi deequMetricInstanceMongoDao) *abstract.DeequMetricInstance {
	return &abstract.DeequMetricInstance{
		Analyzer: abstract.DeequAnalyzer{
			AnalyzerName: mi.Analyzer.AnalyzerName,
			Column:       mi.Analyzer.Column,
		},
		Metric: abstract.DeequMetricValue{
			MetricName: mi.Metric.MetricName,
			Entity:     mi.Metric.Entity,
			Instance:   mi.Metric.Instance,
			Name:       mi.Metric.Name,
			Value:      mi.Metric.Value,
		},
	}
}

func convertAllMetricInstancesDAOToDTO(inMi []deequMetricInstanceMongoDao) []abstract.DeequMetricInstance {
	var mi []abstract.DeequMetricInstance
	for _, element := range inMi {
		mi = append(mi, *convertMetricInstanceDAOToDTO(element))
	}
	return mi
}

func convertMetricDAOToDTO(mmd *metricMongoDao) *abstract.Metric {
	result := &abstract.Metric{}
	if mmd.DeequMetric != nil {
		result.DeequMetric = &abstract.DeequMetric{
			ResultKey: &abstract.DeequResultKey{
				DataSetDate: mmd.DeequMetric.ResultKey.DataSetDate,
				Tags:        mmd.DeequMetric.ResultKey.Tags,
			},
			AnalyzerContext: &abstract.DeequAnalyzerContext{
				MetricMap: convertAllMetricInstancesDAOToDTO(mmd.DeequMetric.AnalyzerContext.MetricMap),
			},
		}
	}
	return result
}

func convertMetricSetDAOToDTO(msmd *metricSetMongoDao) *abstract.MetricSet {
	return &abstract.MetricSet{
		Name:        msmd.Name,
		InsertedAt:  msmd.InsertedAt,
		Version:     msmd.Version,
		Description: msmd.Description,
		Labels:      msmd.Labels,
		Metrics:     convertAllMetricsDAOToDTO(&msmd.Metrics),
	}
}

func convertAllMetricSetsDAOToDTO(inFeats *[]metricSetMongoDao) []abstract.MetricSet {
	var metricSets []abstract.MetricSet
	for _, element := range *inFeats {
		metricSets = append(metricSets, *convertMetricSetDAOToDTO(&element))
	}
	return metricSets
}

func convertAllMetricsDAOToDTO(inMetrics *[]metricMongoDao) []abstract.Metric {
	var metrics []abstract.Metric
	for _, element := range *inMetrics {
		metrics = append(metrics, *convertMetricDAOToDTO(&element))
	}
	return metrics
}

// ------------------------------

var timeout = 5 * time.Second

type dao struct {
	Connector *mongo.Connector
}

// both init and sync.Once are thread-safe
// but only sync.Once is lazy
var once sync.Once
var instance *dao

// GetSingleton ... lazy singleton on DAO
func GetSingleton() abstract.MetricSetDAOProvider {
	// once.do is lazy, we use it to return an instance of the DAO
	once.Do(func() {
		instance = &dao{}
	})
	return instance
}

func (dao *dao) Init(def *conf.DataSourceDefinition) {
	// create mongo connector
	dao.Connector = mongo.NewMongoConnector()
	// validate data source definition
	if err := dao.Connector.ValidateDataSourceDefinition(def); err != nil {
		panic(err)
	}
	// init mongo connector
	dao.Connector.InitConnection(def)

	if err := dao.EnsureIndexesExist(); err != nil {
		panic(err)
	}
}

func (dao *dao) EnsureIndexesExist() error {
	ctx := context.Background()
	// make sure a full text index exists on the description
	indexModel := mongodriver.IndexModel{
		Keys: bsonx.Doc{{Key: "description", Value: bsonx.String("text")}},
	}
	if _, err := dao.Connector.Collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		return err
	}
	return nil
}

func (dao *dao) Create(ms *abstract.MetricSet) error {
	msmd := convertMetricSetDTOtoDAO(ms)

	bsonVal, err := bson.Marshal(msmd)
	if err != nil {
		return err
	}

	// insert
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := dao.Connector.Collection.InsertOne(ctx, bsonVal)
	if err != nil {
		return fmt.Errorf("Error while creating metric set :: %v", err)
	}
	id := res.InsertedID
	log.Printf("Inserted MetricSet %d", id)
	return nil
}

func (dao *dao) CloseConnection() {
	dao.Connector.CloseConnection()
}

func (dao *dao) getOneDocumentUsingFilter(filter interface{}) (*abstract.MetricSet, error) {
	var result metricSetMongoDao
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := dao.Connector.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving metric set :: %v", err)
	}

	// convert DAO to DTO
	return convertMetricSetDAOToDTO(&result), nil
}

type sorter struct {
	sortField string
	sortValue interface{}
}

func (dao *dao) getAnyDocumentUsingFilter(filter interface{}, sorter *sorter, limit int, page int) (*abstract.PaginatedMetricSets, error) {
	var metrics []metricSetMongoDao

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	paginator := paginate.New(dao.Connector.Collection).Context(ctx).Limit(int64(limit)).Page(int64(page)).Filter(filter)
	if sorter != nil {
		paginator = paginator.Sort(sorter.sortField, sorter.sortValue)
	}
	paginatedData, err := paginator.Decode(&metrics).Find()
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving asset :: %v", err)
	}
	if metrics == nil {
		return nil, fmt.Errorf("Error while retrieving metrics using filter :: empty result set")
	}

	var resultMetrics []abstract.MetricSet = convertAllMetricSetsDAOToDTO(&metrics)
	//return &resultMetrics, nil
	return &abstract.PaginatedMetricSets{
		Data:       &resultMetrics,
		Pagination: abstract.FromMongoPaginationData(paginatedData.Pagination),
	}, nil
}

// GetById ... Retrieve document by given id
func (dao *dao) GetById(id string) (*abstract.MetricSet, error) {
	filter := bson.M{"_id": id}
	return dao.getOneDocumentUsingFilter(filter)
}

// GetByName ... Retrieve document by given name
func (dao *dao) GetByName(name string, limit int, page int) (*abstract.PaginatedMetricSets, error) {
	filter := bson.M{"name": name}
	sorter := &sorter{"inserted-at", -1}
	return dao.getAnyDocumentUsingFilter(filter, sorter, limit, page)
}

// SearchMetricSetsByLabels ... Retrieve assets by given labels
func (dao *dao) SearchMetricSetsByLabels(labels map[string]string, limit int, page int) (*abstract.PaginatedMetricSets, error) {
	// https://docs.mongodb.com/manual/reference/operator/query/
	// we can not simply use filter := bson.M{"labels": bson.M{"$eq": labels}} since the order of the keys would matter
	// using this the result would be non-deterministic (empty, and not empty)
	// so the solution is to use dot notation instead
	// https://stackoverflow.com/questions/37303989/exact-match-on-the-embedded-document-when-field-order-is-not-known
	// https://docs.mongodb.com/manual/tutorial/query-embedded-documents/#query-on-nested-fields
	filter := bson.M{}
	for k, v := range labels {
		filter["labels."+k] = v
	}
	sorter := &sorter{"inserted-at", -1}
	return dao.getAnyDocumentUsingFilter(filter, sorter, limit, page)
}

// ListAllMetricSets ... Return all MetricSets in index
func (dao *dao) ListAllMetricSets(limit int, page int) (*abstract.PaginatedMetricSets, error) {
	filter := bson.M{}
	var sorter *sorter = nil
	return dao.getAnyDocumentUsingFilter(filter, sorter, limit, page)
}

// Search ... Return all metric sets matching the text search query
func (dao *dao) Search(query string, limit int, page int) (*abstract.PaginatedMetricSets, error) {
	filter := bson.M{
		"$text": bson.M{"$search": query},
	}
	sorter := &sorter{sortField: "score", sortValue: bson.M{"$meta": "textScore"}}
	return dao.getAnyDocumentUsingFilter(filter, sorter, limit, page)
}
