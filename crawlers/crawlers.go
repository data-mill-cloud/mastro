package main

import (
	"fmt"
	"time"

	"log"

	"github.com/go-co-op/gocron"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/crawlers/hdfs"
	"github.com/data-mill-cloud/mastro/crawlers/hive"
	"github.com/data-mill-cloud/mastro/crawlers/impala"
	"github.com/data-mill-cloud/mastro/crawlers/local"
	"github.com/data-mill-cloud/mastro/crawlers/s3"

	"github.com/go-resty/resty/v2"
)

var factories = map[string]func() abstract.Crawler{
	"local":  local.NewCrawler,
	"hdfs":   hdfs.NewCrawler,
	"s3":     s3.NewCrawler,
	"impala": impala.NewCrawler,
	"hive":   hive.NewCrawler,
}

var client = resty.New()

// Start ... Starts the crawler defined in the provided config
func Start(cfg *conf.Config) (abstract.Crawler, error) {
	// start crawler defined in Config
	var crawler abstract.Crawler

	if crawlerFactory, ok := factories[cfg.DataSourceDefinition.Type]; ok {
		// call factory for selected crawler
		crawler = crawlerFactory()
		// init connection on the selected crawler
		crawler.InitConnection(cfg)
		log.Println("Successfully initialized connection", cfg.DataSourceDefinition.Name)

		if cfg.DataSourceDefinition.CrawlerDefinition.Schedule != nil {
			// schedule crawler if a schedule is set
			scheduler := gocron.NewScheduler(time.UTC)

			_, err := scheduler.
				Cron(*cfg.DataSourceDefinition.CrawlerDefinition.Schedule).
				Do(Reconcile, crawler, cfg)

			// if err get out
			if err != nil {
				return nil, err
			}

			log.Printf("Scheduled crawler with %s", *cfg.DataSourceDefinition.CrawlerDefinition.Schedule)

			// start a run right now if necessary
			if cfg.DataSourceDefinition.CrawlerDefinition.StartNow != nil && *cfg.DataSourceDefinition.CrawlerDefinition.StartNow {
				log.Println("Starting first run")
				go Reconcile(crawler, cfg)
			}

			// start gocron - move outside if we decide to start multiple crawlers within the same agent
			//<-gocron.Start()
			scheduler.StartAsync() // start and continue
			//s.StartBlocking() // start scheduler and wait
		} else {
			// start now!
			go Reconcile(crawler, cfg)
		}

		return crawler, nil
	}
	return nil, fmt.Errorf("Impossible to find specified Crawler %s", cfg.DataSourceDefinition.Type)
}

// Reconcile ... call to walkWithFilter to traverse the FS tree and post all found assets to the catalogue endpoint
func Reconcile(crawler abstract.Crawler, cfg *conf.Config) {
	log.Println("Running crawler", cfg.DataSourceDefinition.Name)
	assets, err := crawler.WalkWithFilter(cfg.DataSourceDefinition.CrawlerDefinition.Root, cfg.DataSourceDefinition.CrawlerDefinition.FilterFilename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("Found %d assets to merge in catalogue", len(assets))
	// call a remote catalogue endpoint to add those assets that were just found
	// https://github.com/go-resty/resty/blob/master/example_test.go
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(assets).
		//Post(cfg.DataSourceDefinition.CrawlerDefinition.CatalogueEndpoint)
		Put(cfg.DataSourceDefinition.CrawlerDefinition.CatalogueEndpoint)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// print response info
	ti := resp.Request.TraceInfo()

	//log.Printf("Catalogue response - status:%s statusCode:%d time:%v body:%s", resp.Status(), resp.StatusCode(), ti.ResponseTime, string(resp.Body()))
	log.Printf("Catalogue response - status:%s statusCode:%d time:%v", resp.Status(), resp.StatusCode(), ti.ResponseTime)
}
