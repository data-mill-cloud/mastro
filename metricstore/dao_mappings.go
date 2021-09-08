package main

import (
	"fmt"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/metricstore/daos/mongo"
)

// available backends - lazy loaded singleton DAOs
var availableDAOs = map[string]func() abstract.MetricSetDAOProvider{
	"mongo": mongo.GetSingleton,
}

func selectDao(cfg *conf.Config) (abstract.MetricSetDAOProvider, error) {
	if singletonDao, ok := availableDAOs[cfg.DataSourceDefinition.Type]; ok {
		// call singleton constructor on dao
		return singletonDao(), nil
	}
	return nil, fmt.Errorf("Impossible to find specified DAO connector %s", cfg.DataSourceDefinition.Type)
}
