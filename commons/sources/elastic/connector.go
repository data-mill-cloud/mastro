package elastic

import (
	"fmt"
	"io/ioutil"
	"log"

	"strings"

	"github.com/datamillcloud/mastro/commons/utils/conf"
	stringutils "github.com/datamillcloud/mastro/commons/utils/strings"
	es7 "github.com/elastic/go-elasticsearch/v7"
)

var requiredFields = map[string]string{
	"esUser":  "username",
	"esPwd":   "password",
	"esHosts": "hosts",
	"esIndex": "index",
}

var optionalFields = map[string]string{
	"cert": "cert",
}

// NewElasticConnector factory
func NewElasticConnector() *Connector {
	return &Connector{}
}

/*
https://www.elastic.co/blog/the-go-client-for-elasticsearch-working-with-data
*/

// todo: find a way not to export this

// Connector ... Connector type
type Connector struct {
	Client    *es7.Client
	IndexName string
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
		// https://stackoverflow.com/questions/28799110/how-to-join-a-slice-of-strings-into-a-single-string
		return fmt.Errorf("The following fields are missing from the data source configuration: %s", strings.Join(missingFields, ","))
	}

	log.Println("Successfully validated data source definition")
	return nil
}

// InitConnection ... Starts a connection with Elastic Search
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	var err error
	//c.client, err = es7.NewDefaultClient()
	elasticHostnames := stringutils.SplitAndTrim(def.Settings[requiredFields["esHosts"]], ",")

	esConfig := es7.Config{
		Addresses: elasticHostnames,
		Username:  def.Settings[requiredFields["esUser"]],
		Password:  def.Settings[requiredFields["esPwd"]],
	}
	// if encryption is enabled then set the server certificate
	if certFile, exist := def.Settings[optionalFields["cert"]]; exist {
		cert, err := ioutil.ReadFile(certFile)
		if err != nil {
			log.Fatal("Error while reading certificate", err)
		}
		esConfig.CACert = cert
	}

	c.Client, err = es7.NewClient(esConfig)
	// set the index for the client
	c.IndexName = def.Settings[requiredFields["esIndex"]]

	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	log.Println("Successfully connected to ES")
	log.Println(res)
}

func (c *Connector) CloseConnection() {
	c.CloseConnection()
}
