package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/commons/utils/ux"
	"github.com/kelseyhightower/envconfig"
)

func waitForCtrlC() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		endWaiter.Done()
	}()
	endWaiter.Wait()
}

func loadCfg() *conf.Config {
	err := envconfig.Process("mastro", &conf.Args)
	if err != nil {
		log.Printf("Impossible to parse from env vars - %v", err.Error())
		log.Printf("Attempting parsing string arguments")
		arg.MustParse(&conf.Args)
	}
	// load config from file
	return conf.Load(conf.Args.Config)
}

func start() {
	switch Cfg.ConfigType {
	case "catalogue":
		StartEndpoint(Cfg)
	default:
		log.Println("Invalid config type", Cfg.ConfigType)
	}
}

var (
	// Cfg ... global Config
	Cfg *conf.Config
)

func main() {
	log.Println("Starting")
	log.Println(ux.Header)
	log.Println(ux.Description)

	// load configuration
	Cfg = loadCfg()

	// start selected service
	start()

	log.Println("Waiting for Ctrl+C...")
	waitForCtrlC()
}
