package main

import (
	"fmt"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/featurestore/daos/elastic"
	"github.com/data-mill-cloud/mastro/featurestore/daos/mongo"
)

// available backends - lazy loaded singleton DAOs
var availableDAOs = map[string]func() abstract.FeatureSetDAOProvider{
	"mongo":   mongo.GetSingleton,
	"elastic": elastic.GetSingleton,
}

func selectDao(cfg *conf.Config) (abstract.FeatureSetDAOProvider, error) {
	if singletonDao, ok := availableDAOs[cfg.DataSourceDefinition.Type]; ok {
		// call singleton constructor on dao
		return singletonDao(), nil
	}
	return nil, fmt.Errorf("Impossible to find specified DAO connector %s", cfg.DataSourceDefinition.Type)
}
