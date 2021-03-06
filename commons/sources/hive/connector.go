package hive

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/beltran/gohive"
	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
)

type authType string

const (
	kerberos authType = "kerberos"
	plain             = "plain"
	none              = "none"
)

// NewHiveConnector ... connector constructor
func NewHiveConnector() *Connector {
	return &Connector{
		ConfigurableConnector: abstract.ConfigurableConnector{
			RequiredFields: map[string]string{
				"host":      "host",
				"port":      "port",
				"auth-type": "auth-type",
			},
			OptionalFields: map[string]string{
				"kerberos-service-name": "kerberos-service-name",
			},
		},
	}
}

// Connector ... hive connector
type Connector struct {
	abstract.ConfigurableConnector
	connection *gohive.Connection
}

// InitConnection ... init connection
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	var err error

	host := def.Settings[c.RequiredFields["host"]]
	port, err := strconv.Atoi(def.Settings[c.RequiredFields["port"]])
	if err != nil {
		panic(err)
	}

	configuration := gohive.NewConnectConfiguration()

	// todo: convert all settings to map[string]interface{}
	authType := authType(def.Settings[c.RequiredFields["auth-type"]])
	switch authType {
	case kerberos:
		configuration.Service = def.Settings[c.OptionalFields["krb-service-name"]]
		c.connection, err = gohive.Connect(host, port, "KERBEROS", configuration)
	case plain:
		configuration.Username = def.Settings[c.OptionalFields["username"]]
		configuration.Password = def.Settings[c.OptionalFields["password"]]
		c.connection, err = gohive.Connect(host, port, "NONE", configuration)
	case none:
		c.connection, err = gohive.Connect(host, port, "NOSASL", configuration)
	default:
		log.Panicf("Auth type %s not available!", authType)
	}

	if err != nil {
		panic(err)
	}
}

// CloseConnection ... close connection
func (c *Connector) CloseConnection() {
	c.connection.Close()
}

// ListDatabases ...
func (c *Connector) ListDatabases() ([]abstract.DBInfo, error) {
	result := []abstract.DBInfo{}

	cursor := c.connection.Cursor()
	defer cursor.Close()
	ctx := context.Background()
	cursor.Exec(ctx, "show databases")
	if cursor.Err != nil {
		return nil, cursor.Err
	}
	for cursor.HasMore(ctx) {
		db := abstract.DBInfo{}
		cursor.FetchOne(ctx, &db.Name, &db.Comment)
		result = append(result, db)
	}
	return result, nil
}

// ListTables ...
func (c *Connector) ListTables(dbName string) ([]abstract.TableInfo, error) {
	//result := []string{}
	result := []abstract.TableInfo{}

	cursor := c.connection.Cursor()
	defer cursor.Close()
	ctx := context.Background()
	cursor.Exec(ctx, fmt.Sprintf("show tables in %s", dbName))
	if cursor.Err != nil {
		return nil, cursor.Err
	}

	for cursor.HasMore(ctx) {
		var tableName string
		// show tables only return table names!
		cursor.FetchOne(ctx, &tableName)
		if cursor.Err != nil {
			return nil, cursor.Err
		}
		//result = append(result, tableName)
		tableInfo, err := abstract.GetTableInfoByName(tableName)
		if err != nil {
			return nil, err
		}
		result = append(result, tableInfo)
	}

	return result, nil
}

// DescribeTable ...
func (c *Connector) DescribeTable(dbName string, tableName string) (map[string]abstract.ColumnInfo, error) {
	var result = make(map[string]abstract.ColumnInfo)

	cursor := c.connection.Cursor()
	defer cursor.Close()
	ctx := context.Background()
	cursor.Exec(ctx, fmt.Sprintf("describe %s.%s", dbName, tableName))
	if cursor.Err != nil {
		return nil, cursor.Err
	}

	for cursor.HasMore(ctx) {
		var cName string
		cInfo := abstract.ColumnInfo{}

		cursor.FetchOne(ctx, &cName, &(cInfo.Type), &(cInfo.Comment))
		result[cName] = cInfo
	}

	return result, nil
}
