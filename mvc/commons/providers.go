package commons

import (
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
)

type MvcProvider interface {
	InitConnection(cfg *conf.Config) (MvcProvider, error)
	InitDataset(cmd *InitCmd)
	NewVersion(cmd *NewCmd)
	Add(cmd *AddCmd)
	AllVersions(cmd *VersionsCmd)
	LatestVersion(cmd *LatestCmd)
	OverwriteVersion(cmd *OverwriteCmd)
	DeleteVersion(cmd *DeleteCmd)
}
