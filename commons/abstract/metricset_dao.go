package abstract

import "github.com/data-mill-cloud/mastro/commons/utils/conf"

// MetricSetDAOProvider ... The interface each dao must implement
type MetricSetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(m *MetricSet) error
	GetById(id string) (*MetricSet, error)
	GetByName(name string, limit int, page int) (*PaginatedMetricSets, error)
	SearchMetricSetsByLabels(labels map[string]string, limit int, page int) (*PaginatedMetricSets, error)
	ListAllMetricSets(limit int, page int) (*PaginatedMetricSets, error)
	Search(query string, limit int, page int) (*PaginatedMetricSets, error)
	CloseConnection()
}
