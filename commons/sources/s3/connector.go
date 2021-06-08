package s3

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// NewS3Connector factory
func NewS3Connector() *Connector {
	return &Connector{}
}

// Connector ... S3 connector
type Connector struct {
	client *minio.Client
	Bucket string
	Prefix string
	Region string
}

var requiredFields = map[string]string{
	"endpoint":  "endpoint",
	"accesskey": "access-key-id",
	"secretkey": "secret-access-key",
	"usessl":    "use-ssl",
	"bucket":    "bucket",
}

var optionalFields = map[string]string{
	"region" : "region",
}

// GetClient ... Returns the client from the connector
func (c *Connector) GetClient() *minio.Client {
	return c.client
}

// ValidateDataSourceDefinition ... Validates the input data source definition
func (c *Connector) ValidateDataSourceDefinition(def *conf.DataSourceDefinition) error {
	// check all required fields are available
	var missingFields []string
	for _, reqvalue := range requiredFields {
		if _, exist := def.Settings[reqvalue]; !exist {
			missingFields = append(missingFields, reqvalue)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("The following fields are missing from the data source configuration: %s", strings.Join(missingFields, ","))
	}

	_, err := strconv.ParseBool(def.Settings[requiredFields["usessl"]])
	if err != nil {
		return fmt.Errorf("Impossible to convert usessl to boolean")
	}

	log.Println("Successfully validated data source definition")
	return nil
}

// InitConnection ... inits connection
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {

	endpoint := def.Settings[requiredFields["endpoint"]]
	accessKeyID := def.Settings[requiredFields["accesskey"]]
	secretKey := def.Settings[requiredFields["secretKey"]]
	useSSL, _ := strconv.ParseBool(def.Settings[requiredFields["usessl"]])
	bucket := def.Settings[requiredFields["bucket"]]

	// optional, e.g. when using minio this is not necessary
	var exist bool
	if c.Region, exist = def.Settings[optionalFields["region"]]; exist {
		log.Println(fmt.Sprintf("Using provided region %s", c.Region))
	} 

	var err error
	c.client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretKey, ""),
		Secure: useSSL,
	})
	c.Bucket = bucket

	if err != nil {
		log.Panicln(err)
	}
}

// CloseConnection ... terminates the connection
func (c *Connector) CloseConnection() {
	// c.client.Close() // close not available
}
