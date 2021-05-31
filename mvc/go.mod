module github.com/datamillcloud/mvc

go 1.14

replace github.com/datamillcloud/mastro/commons => ../commons

require (
	github.com/alexflint/go-arg v1.4.2
	github.com/cheggaaa/pb v1.0.29
	github.com/cheggaaa/pb/v3 v3.0.8
	github.com/datamillcloud/mastro/commons v0.0.0-00010101000000-000000000000
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/minio/mc v0.0.0-20210519052947-8d65ee71302e // indirect
	github.com/minio/minio-go/v7 v7.0.11-0.20210517200026-f0518ca447d6
	gopkg.in/yaml.v2 v2.3.0
)
