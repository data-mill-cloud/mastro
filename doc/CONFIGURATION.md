# Mastro
## Configuration

The package `conf` defines the structure of the Yaml configuration, to be provided as input.
The config can be used to start one of the three different types: i) crawler, ii) catalogue or iii) featurestore.
This is defined using the `ConfigType`, an alias for those cases.
Additional `Details` are also provided as a map to start the component.
Each component is defined by a `DataSourceDefinition` defining the connection details to a backend persistence service.

```go
// Config ... Defines a model for the input config files
type Config struct {
	ConfigType           ConfigType           `yaml:"type"`
	Details              map[string]string    `yaml:"details,omitempty"`
	DataSourceDefinition DataSourceDefinition `yaml:"backend"`
}

// ConfigType ... config type
type ConfigType string

const (
	// Crawler ... crawler agent config type
	Crawler ConfigType = "crawler"
	// Catalogue ... catalogue config type
	Catalogue = "catalogue"
	// FeatureStore ... featurestore config type
	FeatureStore = "featurestore"
)
```

The `DataSourceDefinition` is defined as a user-selected `name` and a `type`.

```go
// DataSourceDefinition ... connection details for a data source connector
type DataSourceDefinition struct {
	Name              string            `yaml:"name"`
	Type              string            `yaml:"type"`
	CrawlerDefinition CrawlerDefinition `yaml:"crawler,omitempty"`
	Settings          map[string]string `yaml:"settings,omitempty"`
	// optional kerberos section
	KerberosDetails *KerberosDetails `yaml:"kerberos"`
	// optional tls section
	TLSDetails *TLSDetails `yaml:"tls"`
}
```

A `CrawlerDefinition` is optionally provided to the `crawler` component to determine scraping information.

```go
// CrawlerDefinition ... Config for a Crawler service
type CrawlerDefinition struct {
	Root              string `yaml:"root"`
	FilterFilename    string `yaml:"filter-filename"`
	ScheduleEvery     Period `yaml:"schedule-period"`
	ScheduleValue     uint64 `yaml:"schedule-value"`
	StartNow          bool   `yaml:"start-now"`
	CatalogueEndpoint string `yaml:"catalogue-endpoint"`
}
```

### Feature store

An example configuration for a feature store is defined below:

```yaml
type: featurestore
details:
  port: 8085
backend:
  name: test-mongo
  type: mongo
  settings:
    username: mongo
    password: test
    host: "localhost:27017"
    database: mastro
    collection: mastro-featurestore
```

### Catalogue

An example configuration for a mongo-based catalogue service is defined below:

```yaml
type: catalogue
details:
  port: 8085
backend:
  name: test-mongo
  type: mongo
  settings:
    username: mongo
    password: test
    host: "localhost:27017"
    database: mastro
    collection: mastro-catalogue
```

### Crawler

An example configuration for an S3 crawler is defined below:

```yaml
type: crawler
backend:
  name: public-minio-s3
  type: s3
  crawler:
    root: ""
    filter-filename: "MANIFEST.yaml"
    schedule-period: ""
    schedule-value: 1
    catalogue-endpoint: "localhost:8085"
  settings:
    endpoint: "play.min.io"
    access-key-id: "Q3AM3UQ867SPQQA43P2F"
    secret-access-key: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
    use-ssl: "true"
```

```yaml
type: crawler
backend:
  name: local-impala
  type: impala
  crawler:
    root: ""
    schedule-period: "sunday"
    schedule-value: 1
    start-now: true
    catalogue-endpoint: "http://localhost:8085/assets"
  settings:
    host: "localhost"
    port: "21000"
    use-kerberos: false
```