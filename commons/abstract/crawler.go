package abstract

import (
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
)

const DefaultManifestFilename string = "MANIFEST.yaml"

type Crawler interface {
	InitConnection(cfg *conf.Config) (Crawler, error)
	WalkWithFilter(root string, filenameFilter string) ([]Asset, error)
}
