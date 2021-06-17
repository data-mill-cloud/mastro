package hdfs

import (
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	cfg     *conf.Config
	connector *Connector
	connErr error
)

func TestMain(m *testing.M) {
	cfg = conf.Load("example_hdfs.cfg")
	connector = NewHDFSConnector()
	connector.InitConnection(&cfg.DataSourceDefinition)

	// run module tests
	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestConnection(t *testing.T) {
	t.Log("TestConnection running")
	assert := assert.New(t)
	assert.Equal(connErr, nil)
	assert.NotEqual(connector.GetClient(), nil)
}

