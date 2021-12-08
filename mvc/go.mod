module github.com/data-mill-cloud/mastro/mvc

go 1.15

replace github.com/data-mill-cloud/mastro/commons => ../commons

require (
	github.com/alexflint/go-arg v1.4.2
	github.com/cheggaaa/pb v1.0.29
	github.com/data-mill-cloud/mastro/commons v0.0.0
	github.com/fatih/color v1.10.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/minio/minio-go/v7 v7.0.11-0.20210517200026-f0518ca447d6
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sger/go-hashdir v0.0.1
	gopkg.in/yaml.v2 v2.4.0
)
