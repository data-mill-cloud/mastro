module github.com/data-mill-cloud/mastro/catalogue

go 1.15

replace github.com/data-mill-cloud/mastro/commons => ../commons

require (
	github.com/alexflint/go-arg v1.4.2
	github.com/data-mill-cloud/mastro/commons v0.0.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.2
	github.com/gobeam/mongo-go-pagination v0.0.8
	github.com/kelseyhightower/envconfig v1.4.0
	go.mongodb.org/mongo-driver v1.7.4
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)
