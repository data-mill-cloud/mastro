package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/mvc/commons"
	"github.com/data-mill-cloud/mastro/mvc/connectors/s3"

	"github.com/kelseyhightower/envconfig"
)

func loadCfg() *conf.Config {
	err := envconfig.Process("mvc", &conf.Args)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// load config from file
	return conf.Load(conf.Args.Config)
}

var (
	// Cfg ... global Config
	Cfg *conf.Config
)

const DefaultManifestFilename string = "MANIFEST.yaml"

var manifestFilename = DefaultManifestFilename

var args struct {
	Init      *commons.InitCmd      `arg:"subcommand:init"`
	New       *commons.NewCmd       `arg:"subcommand:new"`
	Add       *commons.AddCmd       `arg:"subcommand:add"`
	Latest    *commons.LatestCmd    `arg:"subcommand:latest"`
	Versions  *commons.VersionsCmd  `arg:"subcommand:versions"`
	Delete    *commons.DeleteCmd    `arg:"subcommand:delete"`
	Overwrite *commons.OverwriteCmd `arg:"subcommand:overwrite"`
}

const header string = `
╔╦╗╦  ╦╔═╗
║║║╚╗╔╝║
╩ ╩ ╚╝ ╚═╝`

// factories for available connectors
var factories = map[string]func(string) commons.MvcProvider{
	"s3": s3.NewMvc,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No subcommand provided. -h for help")
		return
	}

	// load standard config
	Cfg = loadCfg()
	if Cfg == nil {
		return
	}

	// parse command arguments
	arg.MustParse(&args)

	// instantiate mvc provider
	var mvc commons.MvcProvider
	if mvcFactory, exists := factories[Cfg.DataSourceDefinition.Type]; exists {
		mvc = mvcFactory(manifestFilename)
	} else {
		fmt.Println(fmt.Sprintf("Specified type %s does not exist!", Cfg.DataSourceDefinition.Type))
		return
	}

	mvc.InitConnection(Cfg)

	// call specific subcommand handler
	switch {
	case args.Init != nil:
		mvc.InitDataset(args.Init)
	case args.New != nil:
		mvc.NewVersion(args.New)
	case args.Add != nil:
		mvc.Add(args.Add)
	case args.Versions != nil:
		mvc.AllVersions(args.Versions)
	case args.Latest != nil:
		mvc.LatestVersion(args.Latest)
	case args.Overwrite != nil:
		mvc.OverwriteVersion(args.Overwrite)
	case args.Delete != nil:
		mvc.DeleteVersion(args.Delete)
	default:
		fmt.Println(fmt.Sprintf("unknown command %q", os.Args[0]))
		return
	}
}
