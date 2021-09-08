package abstract

import "github.com/data-mill-cloud/mastro/commons/utils/conf"

// MetricSetDAOProvider ... The interface each dao must implement
type MetricSetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(m *MetricSet) error
	GetById(id string) (*MetricSet, error)
	GetByName(name string) (*[]MetricSet, error)
	SearchMetricSetsByLabels(labels map[string]string) (*[]MetricSet, error)
	ListAllMetricSets() (*[]MetricSet, error)
	CloseConnection()
}
