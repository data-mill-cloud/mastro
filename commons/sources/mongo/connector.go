package mongo

import (
	"context"
	"fmt"
	"log"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NewMongoConnector ... Factory
func NewMongoConnector() *Connector {
	return &Connector{
		ConfigurableConnector: abstract.ConfigurableConnector{
			RequiredFields: map[string]string{
				// surely needed the DB and the target collection
				"database":   "database",
				"collection": "collection",
			},
			OptionalFields: map[string]string{
				// connect either by providing the credentials separately
				"username": "username",
				"password": "password",
				"host":     "host",
				// or else by specifying the connection string
				"connectionString": "connection-string",
			},
		},
	}
}

// Connector ... struct containing info on how to connect to a mongo db
type Connector struct {
	abstract.ConfigurableConnector
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

// InitConnection ... Instantiate the connection with the remote DB
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	var connectionString string
	var exist bool

	// if connectionString is provided then use it
	if connectionString, exist = def.Settings[c.OptionalFields["connectionString"]]; exist {
		log.Println("Using provided connection string")
	} else {
		log.Println("No connection string, building from mandatory fields")
		// todo: mongo connection string varies a lot, maybe just pass the whole string from a secret rather than composing it here??
		connectionString = fmt.Sprintf(
			"mongodb://%s:%s@%s",
			def.Settings[c.OptionalFields["username"]],
			def.Settings[c.OptionalFields["password"]],
			def.Settings[c.OptionalFields["host"]],
		)
	}

	var err error
	ctx := context.Background()
	//c.DBCLient, err = mongo.NewClient(options.Client().ApplyURI(connectionString))
	//err = c.DBCLient.Connect(context.Background())
	c.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	if err != nil {
		log.Fatal(err)
	} else {
		if err = c.Client.Ping(ctx, readpref.Primary()); err != nil {
			log.Fatal(err)
		} else {
			log.Println("Successfully connected to db")
		}
	}

	// set target db and connections
	c.Database = c.Client.Database(def.Settings[c.RequiredFields["database"]])
	c.Collection = c.Database.Collection(def.Settings[c.RequiredFields["collection"]])
}

// CloseConnection ... Disconnects and deallocates resources
func (c *Connector) CloseConnection() {
	ctx := context.Background()
	c.Client.Disconnect(ctx)
}
