package s3

import (
	"log"
	"strconv"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// NewS3Connector factory
func NewS3Connector() *Connector {
	return &Connector{
		ConfigurableConnector: abstract.ConfigurableConnector{
			RequiredFields: map[string]string{
				"endpoint":  "endpoint",
				"accesskey": "access-key-id",
				"secretkey": "secret-access-key",
				"usessl":    "use-ssl",
			},
			OptionalFields: map[string]string{
				"region": "region",
				"bucket": "bucket",
			},
		},
	}
}

// Connector ... S3 connector
type Connector struct {
	abstract.ConfigurableConnector
	client *minio.Client
	Bucket string
	Prefix string
	Region string
}

// GetClient ... Returns the client from the connector
func (c *Connector) GetClient() *minio.Client {
	return c.client
}

// InitConnection ... inits connection
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	var err error
	endpoint := def.Settings[c.RequiredFields["endpoint"]]
	accessKeyID := def.Settings[c.RequiredFields["accesskey"]]
	secretKey := def.Settings[c.RequiredFields["secretkey"]]

	var useSSL bool
	if useSSL, err = strconv.ParseBool(def.Settings[c.RequiredFields["usessl"]]); err != nil {
		log.Fatalf("Impossible to convert usessl to boolean")
	}

	// optional
	var exist bool
	if c.Region, exist = def.Settings[c.OptionalFields["region"]]; exist {
		log.Printf("Conf setting region '%s'", c.Region)
	}
	// optional - in mvc can be provided by cli
	if c.Bucket, exist = def.Settings[c.OptionalFields["bucket"]]; exist {
		log.Printf("Conf setting bucket '%s'", c.Bucket)
	}

	c.client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Panicln(err)
	}
}

// CloseConnection ... terminates the connection
func (c *Connector) CloseConnection() {
	// c.client.Close() // close not available
}
