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
type MetricSetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(m *MetricSet) error
	GetById(id string) (*MetricSet, error)
	GetByName(name string, limit int, page int) (*Paginated[MetricSet], error)
	SearchMetricSetsByLabels(labels map[string]string, limit int, page int) (*Paginated[MetricSet], error)
	ListAllMetricSets(limit int, page int) (*Paginated[MetricSet], error)
	Search(query string, limit int, page int) (*Paginated[MetricSet], error)
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

A basic interface is defined to retrieve metricSets:

```go
type MetricStoreService interface {
	Init(cfg *conf.Config) *resterrors.RestErr
	CreateMetricSet(ms MetricSet) (*MetricSet, *resterrors.RestErr)
	GetMetricSetByID(msID string) (*MetricSet, *resterrors.RestErr)
	GetMetricSetByName(msName string, limit int, page int) (*Paginated[MetricSet], *resterrors.RestErr)
	SearchMetricSetsByLabels(labels map[string]string, limit int, page int) (*Paginated[MetricSet], *resterrors.RestErr)
	Search(query string, limit int, page int) (*Paginated[MetricSet], *resterrors.RestErr)
	ListAllMetricSets(limit int, page int) (*Paginated[MetricSet], *resterrors.RestErr)
}
```

This is translated to the following endpoint:


| Verb        | Endpoint                           | Maps to                                                                     |
|-------------|------------------------------------|-----------------------------------------------------------------------------|
| **GET**     | /healthcheck/metricstore           | github.com/data-mill-cloud/mastro/metricstore.Ping                          |
| ~~**GET**~~ | ~~/metricstore/id/:metricset_id~~  | ~~github.com/data-mill-cloud/mastro/metricstore.GetMetricSetByID~~          |
| **GET**     | /metricstore/name/:metricset_name  | github.com/data-mill-cloud/mastro/metricstore.GetMetricSetByName            |
| **PUT**     | /metricstore/                      | github.com/data-mill-cloud/mastro/metricstore.CreateMetricSet               |
| **POST**    | /metricstore/labels                | github.com/data-mill-cloud/mastro/metricstore.SearchMetricSetsByLabels      |
| **GET**     | /metricstore/labels                | github.com/data-mill-cloud/mastro/metricStore.SearchMetricSetsByQueryLabels |
| **POST**    | /metricstore/search                | github.com/data-mill-cloud/mastro/metricstore.Search                        |
| ~~**GET**~~ | ~~/metricstore/~~                  | ~~github.com/data-mill-cloud/mastro/metricstore.ListAllMetricSets~~         | 

### Examples

To push a metric set a PUT to `/metricstore/` is used, along with a JSON body of kind:
```bash
{
    "name" : "gilberto",
	"version" : "test-v1.0",
	"description" : "example metricset for testing purposes",
	"labels" : {
	    "refers-to" : "project-gilberto",
	    "environment" : "test"
	},
	"metrics" : [
		...
	]
}
```

where metrics can be of various types. For instance, Deequ returns an AnalysisResult type such as the following:
```bash
[
  {
    "resultKey": {
      "dataSetDate": 1630876393300,
      "tags": {}
    },
    "analyzerContext": {
      "metricMap": [
        {
          "analyzer": {
            "analyzerName": "Size"
          },
          "metric": {
            "metricName": "DoubleMetric",
            "entity": "Dataset",
            "instance": "*",
            "name": "Size",
            "value": 5.0
          }
        },
        {
          "analyzer": {
            "analyzerName": "Minimum",
            "column": "numViews"
          },
          "metric": {
            "metricName": "DoubleMetric",
            "entity": "Column",
            "instance": "numViews",
            "name": "Minimum",
            "value": 0.0
          }
        }
      ]
    }
  }
]
```

To retrieve metric sets by labels we can either use a POST to `metricstore/labels` and provide a JSON dict of label	values, for instance:

```bash
{
   "labels": {
         "environment": "test",
         "refers-to": "project-gilberto"
    },
    "limit": 4,
    "page": 1
}
```

or use a GET to the same URL with a query string of the form `?label1=value1&label2=value2`, such as `localhost:8087/metricstore/labels?environment=test&refers-to=project-gilberto&limit=4&page=1`.

