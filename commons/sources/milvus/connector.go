package milvus

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

// NewMilvusConnector ... Factory
func NewMilvusConnector() *Connector {
	return &Connector{
		ConfigurableConnector: abstract.ConfigurableConnector{
			RequiredFields: map[string]string{
				"endpoint":        "endpoint",
				"collection":      "collection",
				"partitions":      "partitions",
				"shardsNum":       "shards-num",
				"denseVectorSize": "dense-vector-size",
			},
			OptionalFields: map[string]string{
				"description":          "description",
				"denseVectorFieldName": "dense-vector-field-name",
			},
		},
	}
}

// Connector ... struct containing info on how to connect to a milvus db
type Connector struct {
	abstract.ConfigurableConnector
	Client               client.Client
	Collection           string
	DenseVectorFieldName string
	Partitions           []string
}

// InitConnection ... Instantiate the connection with the remote DB
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	var err error
	c.Collection = def.Settings[c.RequiredFields["collection"]]
	c.Partitions = def.Settings[c.RequiredFields["partitions"]]
	if c.Client, err = client.NewGrpcClient(
		context.Background(),
		def.Settings[c.RequiredFields["endpoint"]],
	); err != nil {
		log.Fatal("failed to connect to Milvus:", err.Error())
	}
	log.Printf("Connected to Milvus at %s", def.Settings[c.RequiredFields["endpoint"]])

	var exists bool
	if c.DenseVectorFieldName, exists = def.Settings[c.OptionalFields["denseVectorFieldName"]]; !exists {
		c.DenseVectorFieldName = "dense_vector"
	}

	if err := c.ensureCollectionExists(def); err != nil {
		log.Fatalf("failed to ensure collection exists: %v", err.Error())
	}

	if err := c.ensureIndexExists(def); err != nil {
		log.Fatalf("failed to create index: %v", err.Error())
	}

	if err := c.Client.LoadCollection(context.Background(), c.Collection, false); err != nil {
		log.Fatal("failed to load collection:", err.Error())
	}
}

func (c *Connector) ensureCollectionExists(def *conf.DataSourceDefinition) error {
	var exists bool
	var err error
	if exists, err = c.Client.HasCollection(context.Background(), c.Collection); err != nil {
		return fmt.Errorf("failed to check collection existence: %v", err.Error())
	}

	if exists {
		log.Printf("Milvus collection %s already exists", c.Collection)
	} else {
		description, exist := def.Settings[c.OptionalFields["description"]]
		if !exist {
			description = ""
		}
		denseVectorSize := def.Settings[c.RequiredFields["denseVectorSize"]]
		schemaDef := &entity.Schema{
			CollectionName: c.Collection,
			Description:    description,
			Fields: []*entity.Field{
				{
					Name:       "id",
					DataType:   entity.FieldTypeInt64,
					PrimaryKey: true,
					AutoID:     false,
				},
				/*{
					Name:       "name",
					DataType:   entity.FieldTypeString,
					PrimaryKey: false,
					AutoID:     false,
				},*/
				{
					Name:     c.DenseVectorFieldName,
					DataType: entity.FieldTypeFloatVector,
					TypeParams: map[string]string{
						"dim": denseVectorSize,
					},
				},
			},
		}

		shardsNum, err := strconv.ParseInt(def.Settings[c.RequiredFields["shardsNum"]], 10, 32)
		if err != nil {
			return err
		}

		if err := c.Client.CreateCollection(context.Background(), schemaDef, int32(shardsNum)); err != nil {
			return fmt.Errorf("failed to create collection: %v", err.Error())
		}
		log.Printf("Successfully created Milvus collection %s", c.Collection)
	}

	return nil
}

func (c *Connector) ensureIndexExists(def *conf.DataSourceDefinition) error {
	var err error
	var index *entity.IndexIvfFlat
	if index, err = entity.NewIndexIvfFlat(entity.L2, 1024); err != nil {
		return err
	}

	return c.Client.CreateIndex(context.Background(), c.Collection, c.DenseVectorFieldName, index, false)
}

// CloseConnection ... Disconnects and deallocates resources
func (c *Connector) CloseConnection() {
	c.Client.Close()
}
