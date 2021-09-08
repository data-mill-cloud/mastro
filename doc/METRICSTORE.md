# Mastro

## Metric Store

A metric store is a service to store and version metrics.

Similarly to features, metrics can be computed on any data asset using either a batch or a stream processing pipeline.
For instance, data quality metrics may be computed on a data source regarding the completeness of certain columns or fields and their distribution; This can be used to detect anomalies in the data by observing the evolution of those metrics over time.

```go
// MetricSet ... a timestamped set of Metrics
type MetricSet struct {
	// unique metric set name
	Name string `json:"name,omitempty"`
	// insertion time
	InsertedAt time.Time `json:"inserted_at,omitempty"`
	// version relates to the environment and the pipeline version
	Version string `json:"version,omitempty"`
	// description is related to the metrics and their extraction process and not the datasource they were calculated on
	Description string `json:"description,omitempty"`
	// labels used for query purposes
	Labels map[string]string `json:"labels,omitempty"`

	// actual metrics
	Metrics []Metric `json:"metrics,omitempty"`
}

type Metric struct {
	DeequMetric
}
```

A data access object (DAO) for a metricSet is defined as follows:

```go
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
```

The interface is then implemented for specific targets in the `metricstore/daos/*` packages.

Each DAO also implements a lazy singleton using `sync.once` (see [blog post](https://medium.com/@ishagirdhar/singleton-pattern-in-golang-9f60d7fdab23)).
This way, all DAO implementations can be efficiently linked from a `dao_mappings.go` file, for instance:

```go
var availableDAOs = map[string]func() abstract.MetricSetDAOProvider{
	"mongo":   mongo.GetSingleton,
	"elastic": elastic.GetSingleton,
}
```

## Service

As for the exposed service, the `metricstore/service.go` defines a basic interface to retrieve metricSets:

```go
type Service interface {
	Init(cfg *conf.Config) *errors.RestErr
	CreateMetricSet(ms abstract.MetricSet) (*abstract.MetricSet, *errors.RestErr)
	GetMetricSetByID(msID string) (*abstract.MetricSet, *errors.RestErr)
	GetMetricSetByName(msName string) (*[]abstract.MetricSet, *errors.RestErr)
	SearchMetricSetsByLabels(labels map[string]string) (*[]abstract.MetricSet, *errors.RestErr)
	ListAllMetricSets() (*[]abstract.MetricSet, *errors.RestErr)
}
```

This is translated to the following endpoint:


| Verb        | Endpoint                           | Maps to                                                                |
|-------------|------------------------------------|------------------------------------------------------------------------|
| **GET**     | /healthcheck/metricstore           | github.com/data-mill-cloud/mastro/metricstore.Ping                     |
| ~~**GET**~~ | ~~/metricstore/id/:featureset_id~~ | ~~github.com/data-mill-cloud/mastro/metricstore.GetMetricSetByID~~     |
| **GET**     | /metricstore/name/:featureset_name | github.com/data-mill-cloud/mastro/metricstore.GetMetricSetByName       |
| **PUT**     | /metricstore/                      | github.com/data-mill-cloud/mastro/metricstore.CreateMetricSet          |
| **POST**    | /metricstore/labels                | github.com/data-mill-cloud/mastro/metricstore.SearchMetricSetsByLabels |
| ~~**GET**~~ | ~~/metricstore/~~                  | ~~github.com/data-mill-cloud/mastro/metricstore.ListAllMetricSets~~    | 

### Examples
