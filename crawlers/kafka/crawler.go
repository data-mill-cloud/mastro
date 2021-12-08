package kafka

import (
	"fmt"
	"log"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/sources/kafka"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
)

const timeoutMs = 5 * 1000

type kafkaCrawler struct {
	connector *kafka.Connector
}

func NewCrawler() abstract.Crawler {
	return &kafkaCrawler{}
}

func (crawler *kafkaCrawler) InitConnection(cfg *conf.Config) (abstract.Crawler, error) {
	crawler.connector = kafka.NewKafkaConnector()
	if err := crawler.connector.ValidateDataSourceDefinition(&cfg.DataSourceDefinition); err != nil {
		log.Panicln(err)
	}
	crawler.connector.InitConnection(&cfg.DataSourceDefinition)
	return crawler, nil
}

func (crawler *kafkaCrawler) WalkWithFilter(root string, filter string) ([]abstract.Asset, error) {
	metadata, err := crawler.connector.KafkaAdminClient.GetMetadata(nil, true, timeoutMs)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %s", err)
	}

	var assets []abstract.Asset
	for _, m := range metadata.Topics {
		streamInfo, err := abstract.GetStreamInfoByName(m.Topic)
		if err != nil {
			return nil, err
		}

		schemaVersions, err := crawler.connector.SchemaRegistryClient.GetSchemaVersions(m.Topic)
		if err != nil {
			// skip without failing the entire process
			log.Printf("couldn't get schema versions for topic %s: %v", m.Topic, err)
		} else {
			for _, schemaVersion := range schemaVersions {
				schema, err := crawler.connector.SchemaRegistryClient.GetSchemaByVersion(m.Topic, schemaVersion)
				if err != nil {
					return nil, fmt.Errorf("couldn't get schema for topic %s: %v", m.Topic, err)
				}
				streamInfo.Schema[schema.Version()] = schema.Schema()
			}
		}

		a, err := streamInfo.BuildAsset()
		if err != nil {
			return nil, err
		}
		assets = append(assets, *a)
	}

	return assets, nil
}
