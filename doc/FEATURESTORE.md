# Mastro

## Feature Store

A feature store is a service to store and version features.

A Feature can either be computed on a dataset or a data stream, respectively using a batch or a stream processing pipeline.
This is due to the different life cycle and performance requirements for collecting and serving those data to end applications.

```go
// FeatureSet ... a versioned set of features
type FeatureSet struct {
	Name        string            `json:"name,omitempty"`
	InsertedAt  time.Time         `json:"inserted_at,omitempty"`
	Version     string            `json:"version,omitempty"`
	Features    []Feature         `json:"features,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// Feature ... a named variable with a data type
type Feature struct {
	Name     string      `json:"name,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	DataType string      `json:"data-type,omitempty"`
}
```

A data access object (DAO) for a featureSet is defined as follows:

```go
type FeatureSetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(fs *FeatureSet) error
	GetById(id string) (*FeatureSet, error)
	GetByName(name string) (*[]FeatureSet, error)
	ListAllFeatureSets() (*[]FeatureSet, error)
	CloseConnection()
}
```

The interface is then implemented for specific targets in the `featurestore/daos/*` packages.

Each DAO also implements a lazy singleton using `sync.once` (see [blog post](https://medium.com/@ishagirdhar/singleton-pattern-in-golang-9f60d7fdab23)).
This way, all DAO implementations can be efficiently linked from a `dao_mappings.go` file, for instance:

```go
var availableDAOs = map[string]func() abstract.FeatureSetDAOProvider{
	"mongo":   mongo.GetSingleton,
	"elastic": elastic.GetSingleton,
}
```

## Service

As for the exposed service, the `featurestore/service.go` defines a basic interface to retrieve featureSets:

```go
type Service interface {
	Init(cfg *conf.Config) *errors.RestErr
	CreateFeatureSet(fs abstract.FeatureSet) (*abstract.FeatureSet, *errors.RestErr)
	GetFeatureSetByID(fsID string) (*abstract.FeatureSet, *errors.RestErr)
	GetFeatureSetByName(fsName string) (*[]abstract.FeatureSet, *errors.RestErr)
	ListAllFeatureSets() (*[]abstract.FeatureSet, *errors.RestErr)
}
```

This is translated to the following endpoint:


| Verb        | Endpoint                          | Maps to                                                       |
|-------------|-----------------------------------|---------------------------------------------------------------|
| **GET**     | /healthcheck/featureset           | github.com/pilillo/mastro/featurestore.Ping                   |
| ~~**GET**~~ | ~~/featureset/id/:featureset_id~~ | ~~github.com/pilillo/mastro/featurestore.GetFeatureSetByID~~  |
| **GET**     | /featureset/name/:featureset_name | github.com/pilillo/mastro/featurestore.GetFeatureSetByName    |
| **PUT**     | /featureset/                      | github.com/pilillo/mastro/featurestore.CreateFeatureSet       |
| ~~**GET**~~ | ~~/featureset/~~                  | ~~github.com/pilillo/mastro/featurestore.ListAllFeatureSets~~ | 

### Examples

This is for instance how to add a new featureSet calculated in the test environment of a fictional project.


*PUT* on `localhost:8085/featureset` with body:
```json
{
	"name" : "mypipelinegeneratedfeatureset",
	"version" : "test-v1.0",
	"description" : "example feature set for testing purposes",
	"labels" : {
	    "refers-to" : "project-gilberto",
	    "environment" : "test"
	},
	"features" : [
		{
			"name":"feature1",
			"value":10,
			"data-type":"int"
		},
		{
			"name":"feature2",
			"value":true,
			"data-type":"bool"
		}
	]
}
```

with the service adding a date time for additional versioning and finally replying with:
```json
{
	"name" : "mypipelinegeneratedfeatureset",
    "inserted_at": "2020-11-29T17:24:01.747543Z",
    "version": "test-v1.0",
    "features": [
        {
            "name": "feature1",
            "value": 10,
            "data-type": "int"
        },
        {
            "name": "feature2",
            "value": true,
            "data-type": "bool"
        }
    ],
    "description": "example feature set for testing purposes",
    "labels": {
        "environment": "test",
        "refers-to": "project-gilberto"
    }
}
```

Mind that the `data-type` is provided as additional information, while go(lang) can correctly deserialize primitive values from Json.
Moreover, the name here is used to group featuresets computed by the same process and it is therefore not to be considered as unique.