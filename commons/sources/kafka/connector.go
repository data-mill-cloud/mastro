package kafka

import (
	"fmt"
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/riferrei/srclient"
)

var requiredFields = map[string]string{
	"bootstrapServers": "bootstrap-servers",
}

var optionalFields = map[string]string{
	"schemaRegistryUrl":      "schema-registry-url",
	"schemaRegistryUsername": "schema-registry-username",
	"schemaRegistryPassword": "schema-registry-password",
}

// NewKafkaConnector ... Factory
func NewKafkaConnector() *Connector {
	return &Connector{}
}

type Connector struct {
	KafkaAdminClient     *kafka.AdminClient
	SchemaRegistryClient *srclient.SchemaRegistryClient
}

// ValidateDataSourceDefinition ... validates the provided data source definition
func (c *Connector) ValidateDataSourceDefinition(def *conf.DataSourceDefinition) error {
	// check all required fields are available
	var missingFields []string
	for _, reqvalue := range requiredFields {
		if _, exist := def.Settings[reqvalue]; !exist {
			missingFields = append(missingFields, reqvalue)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("The following %d fields are missing from the data source configuration: %s", len(missingFields), strings.Join(missingFields[:], ","))
	}

	log.Println("Successfully validated data source definition")
	return nil
}

// InitConnection ... Instantiate the connection with the remote service
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	var err error

	// create a new kafka admin client
	c.KafkaAdminClient, err = kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": def.Settings[requiredFields["bootstrapServers"]],
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully instantiated Kafka Admin Client")
	}

	// create a new schema registry client if a url is set
	if schemaRegistryUrl, exist := def.Settings[optionalFields["schemaRegistryUrl"]]; exist {
		log.Printf("Using provided url %s to connect to schema registry", schemaRegistryUrl)
		c.SchemaRegistryClient = srclient.CreateSchemaRegistryClient(schemaRegistryUrl)
		schemaRegistryUsername, existUsername := def.Settings[optionalFields["schemaRegistryUsername"]]
		schemaRegistryPassword, existPassword := def.Settings[optionalFields["schemaRegistryPassword"]]
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
