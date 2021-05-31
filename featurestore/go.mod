module github.com/datamillcloud/mastro/featurestore

go 1.14

replace github.com/datamillcloud/mastro/featurestore => ./

replace github.com/datamillcloud/mastro/commons => ../commons

require (
	github.com/alexflint/go-arg v1.4.2
	github.com/datamillcloud/mastro/commons v0.0.0-00010101000000-000000000000
	github.com/elastic/go-elasticsearch v0.0.0
	github.com/elastic/go-elasticsearch/v7 v7.13.0 // indirect
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.2
	github.com/kelseyhightower/envconfig v1.4.0
	go.mongodb.org/mongo-driver v1.5.2
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
