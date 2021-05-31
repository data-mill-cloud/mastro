package conf

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Args ... Arguments provided either as env vars or string args
var Args struct {
	Config string `required:"true" arg:"-c,required"`
}

// Config ... Defines a model for the input config files
type Config struct {
	ConfigType           ConfigType           `yaml:"type"`
	Details              map[string]string    `yaml:"details,omitempty"`
	DataSourceDefinition DataSourceDefinition `yaml:"backend"`
}

// ConfigType ... config type
type ConfigType string

const (
	// Crawler ... crawler agent config type
	Crawler ConfigType = "crawler"
	// Catalogue ... catalogue config type
	Catalogue = "catalogue"
	// FeatureStore ... featurestore config type
	FeatureStore = "featurestore"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func parseCfg(data []byte) (*Config, error) {
	cfg := &Config{}

	err := yaml.Unmarshal(data, &cfg)
	log.Println("Successfully loaded config", cfg.ConfigType, cfg.DataSourceDefinition.Name)

	return cfg, err
}

func validateCfg(cfg *Config) (*Config, error) {
	// todo add validation of input config
	return cfg, nil
}

// Load ... load configuration from file path
func Load(filename string) *Config {
	if !fileExists(filename) {
		log.Fatalf("Configuration file %s does not exist (or is a directory)", filename)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error - %v", err)
	}

	config, err := parseCfg(data)
	if err != nil {
		log.Fatalf("Error - %v", err)
	}
	config, err = validateCfg(config)
	if err != nil {
		log.Fatalf("Error - %v", err)
	}

	return config
}
