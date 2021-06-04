package impala

import (
	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/sources/impala"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"

	"fmt"
	"log"

	"github.com/data-mill-cloud/mastro/commons/utils/strings"
)

type impalaCrawler struct {
	connector *impala.Connector
}

// NewCrawler ... returns an instance of the crawler
func NewCrawler() abstract.Crawler {
	return &impalaCrawler{}
}

func (crawler *impalaCrawler) InitConnection(cfg *conf.Config) (abstract.Crawler, error) {
	crawler.connector = impala.NewImpalaConnector()
	if err := crawler.connector.ValidateDataSourceDefinition(&cfg.DataSourceDefinition); err != nil {
		return nil, err
	}
	crawler.connector.InitConnection(&cfg.DataSourceDefinition)
	return crawler, nil
}

func (crawler *impalaCrawler) WalkWithFilter(root string, filter string) ([]abstract.Asset, error) {
	var assets []abstract.Asset

	levels := strings.SplitAndTrim(root, "/")

	// create empty map of kind, DBInfo -> []TableInfo
	dbTables := map[*abstract.DBInfo][]abstract.TableInfo{}

	// check if a specific database and table was defined
	// N.B. golang split returns a slice with one element, the empty string so len is 1 and we gotta check it
	// https://stackoverflow.com/questions/28330908/how-to-string-split-an-empty-string-in-go
	if levels != nil && len(levels) > 0 && levels[0] != "" {
		log.Printf("Provided specific db levels to locate: '%s'", root)

		dbInfo, err := abstract.GetDBInfoByName(levels[0])
		if err != nil {
			return nil, err
		}

		// a table is defined, use that
		if len(levels) > 1 {
			// construct TableInfo using the provided table name
			tableInfo, err := abstract.GetTableInfoByName(levels[1])
			if err != nil {
				return nil, err
			}
			// put dbInfo -> [tableInfo]
			dbTables[&dbInfo] = []abstract.TableInfo{tableInfo}
		} else {
			// only db is provided, list all tables, construct table info with sole name
			tables, err := crawler.connector.ListTables(dbInfo.Name)
			if err != nil {
				// error while accessing the sole DB we desired to access
				return nil, err
			}
			log.Printf("Found %d tables in requested database %s: %v", len(tables), dbInfo.Name, tables)
			dbTables[&dbInfo] = tables
		}
	} else {
		// list all databases, skip those we can't access, as may be a right issue
		dbs, err := crawler.connector.ListDatabases()

		if err != nil {
			return nil, err
		}

		// list all tables in available dbs
		for _, dbInfo := range dbs {
			tables, err := crawler.connector.ListTables(dbInfo.Name)
			if err != nil {
				// skipping DB
				log.Println(fmt.Sprintf("Error while accessing DB %s! Skipping..", dbInfo.Name))
			} else {
				// add all found tables to map for given db name
				log.Printf("Found %d tables in database %s: %v", len(tables), dbInfo.Name, tables)
				dbTables[&dbInfo] = tables
			}
		}
	}

	// visit each found db->[]tables
	for dbInfo, tableNames := range dbTables {

		// create an asset for the database
		a, err := dbInfo.BuildAsset()
		if err != nil {
			return nil, err
		}
		assets = append(assets, *a)

		// describe each table in the db - create an asset for each
		for _, tableInfo := range tableNames {
			// map[string]abstract.ColumnInfo
			tableSchema, err := crawler.connector.DescribeTable(dbInfo.Name, tableInfo.Name)
			if err != nil {
				log.Print(fmt.Sprintf("Error while accessing %s.%s! Skipping..", dbInfo.Name, tableInfo.Name))
			} else {
				log.Printf("Retrieved schema for table %s.%s", dbInfo.Name, tableInfo.Name)
				// add table schema
				tableInfo.Schema = tableSchema
				// convert to actual Asset definition
				a, err := tableInfo.BuildAsset()
				if err != nil {
					return nil, err
				}
				assets = append(assets, *a)
			}
		}
	}

	return assets, nil
}
