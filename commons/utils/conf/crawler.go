package conf

// CrawlerDefinition ... Config for a Crawler service
type CrawlerDefinition struct {
	CatalogueEndpoint string  `yaml:"catalogue-endpoint"`
	Root              string  `yaml:"root"`
	FilterFilename    string  `yaml:"filter-filename"`
	Schedule          *string `yaml:"schedule,omitempty"`
	StartNow          *bool   `yaml:"start-now,omitempty"`
}
