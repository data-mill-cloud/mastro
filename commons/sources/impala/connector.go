package impala

import (
	"fmt"
	"strconv"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/koblas/impalathing"
)

func NewImpalaConnector() *Connector {
	return &Connector{
		ConfigurableConnector: abstract.ConfigurableConnector{
			RequiredFields: map[string]string{
				"host":         "host",
				"port":         "port",
				"use-kerberos": "use-kerberos",
			},
			OptionalFields: map[string]string{},
		},
	}
}

type Connector struct {
	abstract.ConfigurableConnector
	connection *impalathing.Connection
}

func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	var err error

	host := def.Settings[c.RequiredFields["host"]]
	port, err := strconv.Atoi(def.Settings[c.RequiredFields["port"]])
	if err != nil {
		panic(err)
	}

	// todo: convert all settings to map[string]interface{}
	if def.Settings[c.RequiredFields["use-kerberos"]] == "true" {
		options := impalathing.WithGSSAPISaslTransport()
		c.connection, err = impalathing.Connect(host, port, options)
	} else {
		c.connection, err = impalathing.Connect(host, port)
	}

	if err != nil {
		panic(err)
	}
}

func (c *Connector) CloseConnection() {
	c.connection.Close()
}

// Impala specific methods and structs

func (c *Connector) ListDatabases() ([]abstract.DBInfo, error) {
	result := []abstract.DBInfo{}

	query, err := c.connection.Query("show databases")
	if err != nil {
		return nil, err
	}
	for query.Next() {
		db := abstract.DBInfo{}
		query.Scan(&db.Name, &db.Comment)

		result = append(result, db)
	}
	return result, nil
}

func (c *Connector) ListTables(dbName string) ([]abstract.TableInfo, error) {
	//result := []string{}
	result := []abstract.TableInfo{}

	query, err := c.connection.Query(fmt.Sprintf("show tables in %s", dbName))
	if err != nil {
		return nil, err
	}

	for query.Next() {
		var tableName string
		query.Scan(&tableName)
		//result = append(result, tableName)
		tableInfo, err := abstract.GetTableInfoByName(tableName)
		if err != nil {
			return nil, err
		}
		result = append(result, tableInfo)
	}

	return result, nil
}

func (c *Connector) DescribeTable(dbName string, tableName string) (map[string]abstract.ColumnInfo, error) {
	var result = make(map[string]abstract.ColumnInfo)

	query, err := c.connection.Query(fmt.Sprintf("describe %s.%s", dbName, tableName))
	if err != nil {
		return nil, err
	}

	for query.Next() {
		var cName string
		cInfo := abstract.ColumnInfo{}

		query.Scan(&cName, &(cInfo.Type), &(cInfo.Comment))
		result[cName] = cInfo
	}

	return result, nil
}
