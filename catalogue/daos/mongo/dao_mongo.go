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
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"gopkg.in/mgo.v2/bson"
)

type assetMongoDao struct {
	//ID primitive.ObjectID `bson:"_id,omitempty"`
	// asset discovery datetime
	LastDiscoveredAt time.Time `bson:"last-discovered-at"`
	// asset publication datetime
	PublishedOn time.Time `bson:"published-on"`
	// name of the asset
	Name string `bson:"_id"`
	// description of the asset
	Description string `bson:"description"`
	// the list of assets this depends on
	DependsOn []string `bson:"depends-on"`
	// asset type
	Type abstract.AssetType `bson:"type"`
	// asset labels
	Labels map[string]interface{} `bson:"labels"`
	// tags are flags used to simplify asset search
	Tags []string `bson:"tags"`
	// versions specify available variants of the same asset
	Versions map[string]interface{} `bson:"versions"`
}

func convertAssetDTOtoDAO(as *abstract.Asset) *assetMongoDao {
	asmd := &assetMongoDao{}

	// id not set at the time of insert (DTO->DAO)
	// however id is set if we are updating an existing asset
	asmd.LastDiscoveredAt = as.LastDiscoveredAt
	asmd.PublishedOn = as.PublishedOn
	asmd.Name = as.Name
	asmd.Description = as.Description
	asmd.DependsOn = as.DependsOn

	asmd.Type = as.Type

	asmd.Labels = as.Labels
	asmd.Tags = as.Tags
	asmd.Versions = as.Versions

	return asmd
}

func convertAssetDAOtoDTO(asmd *assetMongoDao) *abstract.Asset {
	as := &abstract.Asset{}

	as.LastDiscoveredAt = asmd.LastDiscoveredAt
	as.PublishedOn = asmd.PublishedOn
	as.Name = asmd.Name
	as.Description = asmd.Description
	as.DependsOn = asmd.DependsOn

	as.Type = asmd.Type

	as.Labels = asmd.Labels
	as.Tags = asmd.Tags
	as.Versions = asmd.Versions

	return as
}

func convertAllAssets(inAssets *[]assetMongoDao) []abstract.Asset {
	var assets []abstract.Asset
	for _, element := range *inAssets {
		assets = append(assets, *convertAssetDAOtoDTO(&element))
	}
	return assets
}

var once sync.Once
var instance *dao

type dao struct {
	Connector *mongo.Connector
}

var timeout = 5 * time.Second

// GetSingleton ... get an instance of the dao backend
func GetSingleton() abstract.AssetDAOProvider {
	// once.do is lazy, we use it to return an instance of the DAO
	once.Do(func() {
		instance = &dao{}
	})
	return instance
}

// Init ... Initialize connection to db and target index
func (dao *dao) Init(def *conf.DataSourceDefinition) {
	dao.Connector = mongo.NewMongoConnector()
	if err := dao.Connector.ValidateDataSourceDefinition(def); err != nil {
		panic(err)
	}
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

// Upsert ... Upsert asset
func (dao *dao) Upsert(as *abstract.Asset) error {
	asmd := convertAssetDTOtoDAO(as)

	bsonVal, err := bson.Marshal(asmd)
	if err != nil {
		return err
	}

	// https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-update-documents
	// When working with the ReplaceOne function, update operators such as $set cannot be used since
	// it is a complete replace rather than an update of particular fields.
	// replace all fields in a document while maintaining the id of that document
	// https://stackoverflow.com/questions/59311020/upsert-not-working-when-using-updateone-with-the-mongodb-golang-driver

	opts := options.Replace().SetUpsert(true) //.Update().SetUpsert(true)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	session, err := dao.Connector.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionContext mongodriver.SessionContext) (interface{}, error) {
		//filter := bson.M{"name": as.Name}
		filter := bson.M{"_id": as.Name}
		result, err := dao.Connector.Collection.ReplaceOne(ctx, filter, bsonVal, opts)
		if err != nil {
			return nil, err
		}
		//id := result.UpsertedID.(primitive.ObjectID).Hex()

		if result.MatchedCount > 0 {
			// we are updating an existing asset
			log.Printf("Matched %d Asset(s) :: Modified %d documents", result.MatchedCount, result.ModifiedCount)
		} else {
			// upsert/insert
			log.Printf("Upserted %d Asset(s) :: id = '%v'", result.UpsertedCount, result.UpsertedID)
		}
		return result, nil
	}

	if _, err = session.WithTransaction(context.Background(), callback); err != nil {
		return fmt.Errorf("error while upserting asset :: %v", err)
	}

	return nil
}

func (dao *dao) getOneDocumentUsingFilter(filter interface{}) (*abstract.Asset, error) {
	var result assetMongoDao
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := dao.Connector.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving asset :: %v", err)
	}
	return convertAssetDAOtoDTO(&result), nil
}

type sorter struct {
	sortField string
	sortValue interface{}
}

func (dao *dao) getAnyDocumentUsingFilter(filter interface{}, sorter *sorter, limit int, page int) (*abstract.Paginated[abstract.Asset], error) {
	var assets []assetMongoDao

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	paginator := paginate.New(dao.Connector.Collection).Context(ctx).Limit(int64(limit)).Page(int64(page))
	if sorter != nil {
		paginator = paginator.Sort(sorter.sortField, sorter.sortValue)
	}
	paginator = paginator.Filter(filter)

	paginatedData, err := paginator.Decode(&assets).Find()
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving asset :: %v", err)
	}

	if assets == nil {
		return nil, fmt.Errorf("Error while retrieving assets using filter :: empty result set")
	}

	var resultAssets []abstract.Asset = convertAllAssets(&assets)

	return &abstract.Paginated[abstract.Asset]{
		Data:       &resultAssets,
		Pagination: abstract.FromMongoPaginationData(paginatedData.Pagination),
	}, nil
}

// GetById ... Retrieve document by given id
func (dao *dao) GetById(id string) (*abstract.Asset, error) {
	//primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": id}
	return dao.getOneDocumentUsingFilter(filter)
}

// GetByName ... Retrieve document by given name
func (dao *dao) GetByName(name string) (*abstract.Asset, error) {
	//filter := bson.M{"name": name}
	filter := bson.M{"_id": name}
	return dao.getOneDocumentUsingFilter(filter)
}

// SearchAssetsByTags ... Retrieve assets by given tags
func (dao *dao) SearchAssetsByTags(tags []string, limit int, page int) (*abstract.Paginated[abstract.Asset], error) {
	// https://www.mongodb.com/blog/post/quick-start-golang--mongodb--data-aggregation-pipeline
	// https://docs.mongodb.com/manual/tutorial/query-arrays/#match-an-array
	// find all docs whose tags field contains all the elements provided as tags []string in input
	// without regard of the order
	filter := bson.M{"tags": bson.M{"$all": tags}}
	var sorter *sorter = nil
	return dao.getAnyDocumentUsingFilter(filter, sorter, limit, page)
}

// ListAllAssets ... Return all assets in index
func (dao *dao) ListAllAssets(limit int, page int) (*abstract.Paginated[abstract.Asset], error) {
	filter := bson.M{}
	var sorter *sorter = nil
	return dao.getAnyDocumentUsingFilter(filter, sorter, limit, page)
}

// Search ... Return all assets matching the text search query
func (dao *dao) Search(query string, limit int, page int) (*abstract.Paginated[abstract.Asset], error) {
	filter := bson.M{
		"$text": bson.M{"$search": query},
	}
	sorter := &sorter{sortField: "score", sortValue: bson.M{"$meta": "textScore"}}
	return dao.getAnyDocumentUsingFilter(filter, sorter, limit, page)
}

// CloseConnection ... Terminates the connection to ES for the DAO
func (dao *dao) CloseConnection() {
	dao.Connector.CloseConnection()
}
