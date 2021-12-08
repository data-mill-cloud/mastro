module github.com/data-mill-cloud/mastro/metricstore

go 1.15

//replace github.com/data-mill-cloud/mastro/commons => ../commons

require (
	github.com/alexflint/go-arg v1.4.2
	github.com/data-mill-cloud/mastro/commons v0.0.0-20211208161933-320a88c6ba52
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.2
	github.com/kelseyhightower/envconfig v1.4.0
	go.mongodb.org/mongo-driver v1.5.2
)
