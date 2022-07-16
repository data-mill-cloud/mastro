package hdfs

import (
	"os"
	"os/user"

	gohdfs "github.com/colinmarc/hdfs/v2"
	"github.com/colinmarc/hdfs/v2/hadoopconf"
	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/commons/utils/kerberos"
)

// NewHDFSConnector factory
func NewHDFSConnector() *Connector {
	return &Connector{
		ConfigurableConnector: abstract.ConfigurableConnector{
			RequiredFields: map[string]string{},
			OptionalFields: map[string]string{},
		},
	}
}

// Connector ... HDFS connector
type Connector struct {
	abstract.ConfigurableConnector
	client *gohdfs.Client
}

// GetClient ... Returns the client from the connector
func (c *Connector) GetClient() *gohdfs.Client {
	return c.client
}

// InitConnection ... inits connection
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {

	// "HADOOP_CONF_DIR" should be set for this to work
	_, present := os.LookupEnv("HADOOP_CONF_DIR")
	if !present {
		panic("HADOOP_CONF_DIR not set!")
	}

	/*
		LoadFromEnvironment tries to locate the Hadoop configuration files based on the environment,
		and returns a HadoopConf object representing the parsed configuration.
		If the HADOOP_CONF_DIR environment variable is specified, it uses that, or if HADOOP_HOME is specified, it uses $HADOOP_HOME/conf.
	*/
	hadoopConf, err := hadoopconf.LoadFromEnvironment()
	if err != nil {
		panic(err)
	}

	// https://godoc.org/github.com/colinmarc/hdfs#ClientOptionsFromConf
	clientOptions := gohdfs.ClientOptionsFromConf(hadoopConf)

	username, err := user.Current()
	if err != nil {
		panic(err)
	}
	clientOptions.User = username.Username

	if clientOptions.KerberosClient != nil {
		clientOptions.KerberosClient = kerberos.GetKerberosClient(def.KerberosDetails)
	}

	if c.client, err = gohdfs.NewClient(clientOptions); err != nil {
		panic(err)
	}
}

// CloseConnection ... terminates the connection
func (c *Connector) CloseConnection() {
	c.client.Close()
}
