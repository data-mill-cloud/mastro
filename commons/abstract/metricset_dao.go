package abstract

import "github.com/data-mill-cloud/mastro/commons/utils/conf"

// MetricSetDAOProvider ... The interface each dao must implement
type MetricSetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(m *MetricSet) error
	GetByName(id string) (*MetricSet, error)
	SearchMetricSetsByLabels(labels map[string]string) (*[]MetricSet, error)
	CloseConnection()
}
