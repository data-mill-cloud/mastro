package elastic

import (
	"io/ioutil"
	"log"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	stringutils "github.com/data-mill-cloud/mastro/commons/utils/strings"
	es7 "github.com/elastic/go-elasticsearch/v7"
)

// NewElasticConnector factory
func NewElasticConnector() *Connector {
	return &Connector{
		ConfigurableConnector: abstract.ConfigurableConnector{
			RequiredFields: map[string]string{
				"esUser":  "username",
				"esPwd":   "password",
				"esHosts": "hosts",
				"esIndex": "index",
			},
			OptionalFields: map[string]string{
				"cert": "cert",
			},
		},
	}
}

/*
https://www.elastic.co/blog/the-go-client-for-elasticsearch-working-with-data
*/

// todo: find a way not to export this

// Connector ... Connector type
type Connector struct {
	abstract.ConfigurableConnector
	Client    *es7.Client
	IndexName string
}

// InitConnection ... Starts a connection with Elastic Search
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	var err error
	//c.client, err = es7.NewDefaultClient()
	elasticHostnames := stringutils.SplitAndTrim(def.Settings[c.RequiredFields["esHosts"]], ",")

	esConfig := es7.Config{
		Addresses: elasticHostnames,
		Username:  def.Settings[c.RequiredFields["esUser"]],
		Password:  def.Settings[c.RequiredFields["esPwd"]],
	}
	// if encryption is enabled then set the server certificate
	if certFile, exist := def.Settings[c.OptionalFields["cert"]]; exist {
		cert, err := ioutil.ReadFile(certFile)
		if err != nil {
			log.Fatal("Error while reading certificate", err)
		}
		esConfig.CACert = cert
	}

	c.Client, err = es7.NewClient(esConfig)
	// set the index for the client
	c.IndexName = def.Settings[c.RequiredFields["esIndex"]]

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

}
