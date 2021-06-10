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
}

var optionalFields = map[string]string{
	"region": "region",
	"bucket": "bucket",
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
	secretKey := def.Settings[requiredFields["secretkey"]]
	useSSL, _ := strconv.ParseBool(def.Settings[requiredFields["usessl"]])

	// optional
	var exist bool
	if c.Region, exist = def.Settings[optionalFields["region"]]; exist {
		log.Println(fmt.Sprintf("Using specified region %s", c.Region))
	}
	// optional - in mvc can be provided by cli
	if c.Bucket, exist = def.Settings[optionalFields["bucket"]]; exist {
		log.Println(fmt.Sprintf("Using specified bucket %s", c.Bucket))
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
