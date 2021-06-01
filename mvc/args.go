package main

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

var args struct {
	Init      *InitCmd      `arg:"subcommand:init"`
	New       *NewCmd       `arg:"subcommand:new"`
	Add       *AddCmd       `arg:"subcommand:add"`
	Latest    *LatestCmd    `arg:"subcommand:latest"`
	Versions  *VersionsCmd  `arg:"subcommand:versions"`
	Delete    *DeleteCmd    `arg:"subcommand:delete"`
	Overwrite *OverwriteCmd `arg:"subcommand:overwrite"`
}
