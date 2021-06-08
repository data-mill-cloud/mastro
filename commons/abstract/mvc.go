package abstract

type MvcProvider interface {
	InitConnection(cfg *conf.Config) (MvcProvider, error)
	InitDataset(cmd *commons.InitCmd, manifestFilename string)
	NewVersion(cmd *commons.NewCmd, manifestFilename string)
	Add(cmd *commons.AddCmd, manifestFilename string)
	AllVersions(cmd *commons.VersionsCmd, manifestFilename string)
	LatestVersion(cmd *commons.LatestCmd, manifestFilename string)
	OverwriteVersion(cmd *commons.OverwriteCmd)
	DeleteVersion(cmd *commons.DeleteCmd, manifestFilename string)
}