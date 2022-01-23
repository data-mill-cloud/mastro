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

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"gopkg.in/mgo.v2/bson"
)

// featureSetMongoDao ... DAO for the FeatureSet in Mongo
type featureSetMongoDao struct {
	//ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string            `bson:"name,omitempty"`
	InsertedAt  time.Time         `bson:"inserted-at,omitempty"`
	Version     string            `bson:"version,omitempty"`
	Features    []featureMongoDao `bson:"features,omitempty"`
	Description string            `bson:"description,omitempty"`
	Labels      map[string]string `bson:"labels,omitempty"`
}

// featureMongoDao ... a named variable with a data type
type featureMongoDao struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Value    interface{}        `bson:"value,omitempty"`
	DataType string             `bson:"data-type,omitempty"`
}

type dao struct {
	Connector *mongo.Connector
}

var timeout = 5 * time.Second

func convertFeatureDTOtoDAO(f *abstract.Feature) *featureMongoDao {
	fmd := &featureMongoDao{}

	//fmd.ID = f.ID // not set at time of insert
	fmd.Name = f.Name
	fmd.Value = f.Value
	fmd.DataType = f.DataType

	return fmd
}

func convertFeatureSetDTOtoDAO(fs *abstract.FeatureSet) *featureSetMongoDao {
	fsmd := &featureSetMongoDao{}

	//fsmd.ID = fs.ID // not set at time of insert
	fsmd.Name = fs.Name
	fsmd.InsertedAt = fs.InsertedAt
	fsmd.Version = fs.Version

	var feats []featureMongoDao
	for _, element := range fs.Features {
		feats = append(feats, *convertFeatureDTOtoDAO(&element))
	}
	fsmd.Features = feats

	fsmd.Description = fs.Description
	fsmd.Labels = fs.Labels

	return fsmd
}

func convertFeatureDAOToDTO(fmd *featureMongoDao) *abstract.Feature {
	f := &abstract.Feature{}

	//f.ID = fmd.ID.String() // set it in DAO, propagate to DTO?
	f.Name = fmd.Name
	f.Value = fmd.Value
	f.DataType = fmd.DataType

	return f
}

func convertFeatureSetDAOToDTO(fsmd *featureSetMongoDao) *abstract.FeatureSet {
	fs := &abstract.FeatureSet{}

	//fs.ID = fsmd.ID.String() // set it in DAO, propagate to DTO?
	fs.Name = fsmd.Name
	fs.InsertedAt = fsmd.InsertedAt
	fs.Version = fsmd.Version

	fs.Features = convertAllFeatures(&fsmd.Features)
	fs.Description = fsmd.Description
	fs.Labels = fsmd.Labels

	return fs
}

func convertAllFeatureSets(inFeats *[]featureSetMongoDao) []abstract.FeatureSet {
	var feats []abstract.FeatureSet
	for _, element := range *inFeats {
		feats = append(feats, *convertFeatureSetDAOToDTO(&element))
	}
	return feats
}

func convertAllFeatures(inFeats *[]featureMongoDao) []abstract.Feature {
	var feats []abstract.Feature
	for _, element := range *inFeats {
		feats = append(feats, *convertFeatureDAOToDTO(&element))
	}
	return feats
}

// both init and sync.Once are thread-safe
// but only sync.Once is lazy
var once sync.Once
var instance *dao

// GetSingleton ... lazy singleton on DAO
func GetSingleton() abstract.FeatureSetDAOProvider {
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

func (dao *dao) CloseConnection() {
	dao.Connector.CloseConnection()
}

func (dao *dao) Create(fs *abstract.FeatureSet) error {
	// convert DTO to DAO
	//bsonVal := bson.M{"name": "pi", "value": 3.14159}
	fsmd := convertFeatureSetDTOtoDAO(fs)

	bsonVal, err := bson.Marshal(fsmd)
	if err != nil {
		return err
	}

	// insert
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := dao.Connector.Collection.InsertOne(ctx, bsonVal)
	if err != nil {
		return fmt.Errorf("Error while creating feature set :: %v", err)
	}
	id := res.InsertedID
	log.Printf("Inserted FeatureSet %d", id)
	return nil
}

func (dao *dao) getOneDocumentUsingFilter(filter interface{}) (*abstract.FeatureSet, error) {
	var result featureSetMongoDao
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := dao.Connector.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving feature set :: %v", err)
	}

	// convert DAO to DTO
	return convertFeatureSetDAOToDTO(&result), nil
}

type sorter struct {
	sortField string
	sortValue interface{}
}

func (dao *dao) getAnyDocumentUsingFilter(filter interface{}, sorter *sorter, limit int, page int) (*abstract.PaginatedFeatureSets, error) {
	var features []featureSetMongoDao

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	paginator := paginate.New(dao.Connector.Collection).Context(ctx).Limit(int64(limit)).Page(int64(page)).Filter(filter)
	if sorter != nil {
		paginator = paginator.Sort(sorter.sortField, sorter.sortValue)
	}

	paginatedData, err := paginator.Decode(&features).Find()
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving asset :: %v", err)
	}

	if features == nil {
		return nil, fmt.Errorf("Error while retrieving featuresets using filter :: empty result set")
	}

	var resultFeats []abstract.FeatureSet = convertAllFeatureSets(&features)
	return &abstract.PaginatedFeatureSets{
		Data:       &resultFeats,
		Pagination: abstract.FromMongoPaginationData(paginatedData.Pagination),
	}, nil
}

// GetById ... Retrieve document by given id
func (dao *dao) GetById(id string) (*abstract.FeatureSet, error) {
	filter := bson.M{"_id": id}
	return dao.getOneDocumentUsingFilter(filter)
}

// GetByName ... Retrieve document by given name
func (dao *dao) GetByName(name string, limit int, page int) (*abstract.PaginatedFeatureSets, error) {
	filter := bson.M{"name": name}
	sorter := &sorter{"inserted-at", -1}
	return dao.getAnyDocumentUsingFilter(filter, sorter, limit, page)
}

// ListAllFeatureSets ... Return all feature sets available in collection
func (dao *dao) ListAllFeatureSets(limit int, page int) (*abstract.PaginatedFeatureSets, error) {
	filter := bson.M{}
	var sorter *sorter = nil
	return dao.getAnyDocumentUsingFilter(filter, sorter, limit, page)
}

// Search ... Return all featuresets matching the text search query
func (dao *dao) Search(query string, limit int, page int) (*abstract.PaginatedFeatureSets, error) {
	filter := bson.M{
		"$text": bson.M{"$search": query},
	}
	sorter := &sorter{sortField: "score", sortValue: bson.M{"$meta": "textScore"}}
	return dao.getAnyDocumentUsingFilter(filter, sorter, limit, page)
}

// SearchFeatureSetsByLabels ... Return all featuresets matching the search labels
func (dao *dao) SearchFeatureSetsByLabels(labels map[string]string, limit int, page int) (*abstract.PaginatedFeatureSets, error) {
	filter := bson.M{}
	for k, v := range labels {
		filter["labels."+k] = v
	}
	sorter := &sorter{"inserted-at", -1}
	return dao.getAnyDocumentUsingFilter(filter, sorter, limit, page)
}
