package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/riferrei/srclient"
)

// NewKafkaConnector ... Factory
func NewKafkaConnector() *Connector {
	return &Connector{
		ConfigurableConnector: abstract.ConfigurableConnector{
			RequiredFields: map[string]string{
				"bootstrapServers": "bootstrap.servers",
			},
			OptionalFields: map[string]string{
				"schemaRegistryUrl":      "schema.registry.url",
				"schemaRegistryUsername": "schema.registry.username",
				"schemaRegistryPassword": "schema.registry.password",
			},
		},
	}
}

type Connector struct {
	abstract.ConfigurableConnector
	KafkaAdminClient     *kafka.AdminClient
	SchemaRegistryClient *srclient.SchemaRegistryClient
}

// InitConnection ... Instantiate the connection with the remote service
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	var err error

	// create a new kafka admin client (inject directly the connector properties)
	/*
		clientConf := &kafka.ConfigMap{
			"bootstrap.servers": def.Settings[requiredFields["bootstrapServers"]],
		}
	*/
	clientConf := &kafka.ConfigMap{}
	for key, value := range def.Settings {
		clientConf.SetKey(key, value)
	}
	// additional connector properties if any
	c.KafkaAdminClient, err = kafka.NewAdminClient(clientConf)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully instantiated Kafka Admin Client")
	}

	// create a new schema registry client if a url is set
	if schemaRegistryUrl, exist := def.Settings[c.OptionalFields["schemaRegistryUrl"]]; exist {
		log.Printf("Using provided url %s to connect to schema registry", schemaRegistryUrl)
		c.SchemaRegistryClient = srclient.CreateSchemaRegistryClient(schemaRegistryUrl)
		schemaRegistryUsername, existUsername := def.Settings[c.OptionalFields["schemaRegistryUsername"]]
		schemaRegistryPassword, existPassword := def.Settings[c.OptionalFields["schemaRegistryPassword"]]
		if existUsername && existPassword {
			log.Printf("Using provided username %s to connect to schema registry", schemaRegistryUsername)
			c.SchemaRegistryClient.SetCredentials(schemaRegistryUsername, schemaRegistryPassword)
		}
	}

}

// CloseConnection ... Disconnects and deallocates resources
func (c *Connector) CloseConnection() {
	c.KafkaAdminClient.Close()
}
