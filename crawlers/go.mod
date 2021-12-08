module github.com/data-mill-cloud/mastro/crawlers

go 1.15

//replace github.com/data-mill-cloud/mastro/commons => ../commons

require (
	github.com/alexflint/go-arg v1.4.2
	github.com/data-mill-cloud/mastro/commons v0.0.0-20211208161933-320a88c6ba52
	github.com/go-co-op/gocron v1.11.0
	github.com/go-resty/resty/v2 v2.7.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/minio/minio-go/v7 v7.0.17
	github.com/stretchr/testify v1.7.0
)
