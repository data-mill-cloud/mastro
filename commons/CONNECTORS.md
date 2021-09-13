# Mastro
## Connectors

A Data source is defined in the `abstract` package as follows:

```go
type ConnectorProvider interface {
	ValidateDataSourceDefinition(*conf.DataSourceDefinition) error
	InitConnection(*conf.DataSourceDefinition)
	CloseConnection()
}
```

Have a look at the `sources/*` packages for specific implementations of the interface.

A factory is generally used to instantiate the connector with default settings:

```go
func NewElasticConnector() *Connector {
	return &Connector{}
}
```
The connector can be then used for any of the implemented DAOs to be started.
