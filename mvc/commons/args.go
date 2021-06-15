package commons

type InitCmd struct {
	DestinationPath   string  `arg:"-d,required"`
	LocalManifestPath *string `arg:"-f"`
}

type NewCmd struct {
	DestinationPath string `arg:"-d,required"`
}

type AddCmd struct {
	DestinationPath string `arg:"-d,required"`
	LocalPath       string `arg:"-l,required"`
}

type VersionsCmd struct {
	DestinationPath string `arg:"-d,required"`
}

type LatestCmd struct {
	DestinationPath string `arg:"-d,required"`
}

type DeleteCmd struct {
	DestinationPath string `arg:"-d,required"`
	Version         string `arg:"-v,required"`
}

type OverwriteCmd struct {
	DestinationPath string `arg:"-d,required"`
	Version         string `arg:"-v,required"`
	LocalPath       string `arg:"-l,required"`
}

type CheckCmd struct {
	LocalPath string `arg:"-l,required"`
}
