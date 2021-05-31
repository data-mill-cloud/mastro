module github.com/datamillcloud/mastro/crawlers

go 1.14

replace github.com/datamillcloud/mastro/crawlers => ./

replace github.com/datamillcloud/mastro/commons => ../commons

require (
	github.com/alexflint/go-arg v1.4.2
	github.com/datamillcloud/mastro/commons v0.0.0-00010101000000-000000000000
	github.com/go-co-op/gocron v1.6.0
	github.com/go-resty/resty/v2 v2.6.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/minio/minio-go/v7 v7.0.10
	github.com/pilillo/mastro v0.0.0-20210401113305-55cba50ac869
)
