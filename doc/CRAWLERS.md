# Mastro
## Crawlers
A crawler is an agent traversing file systems to seek for asset definition files.
Crawlers implement the Crawler interface:

```go
type Crawler interface {
	InitConnection(cfg *conf.Config) (Crawler, error)
	WalkWithFilter(root string, filenameFilter string) ([]Asset, error)
}
```

Specifically, the crawler inits the connection to a volume (e.g., hdfs, s3) whereas in the WalkWithFilter it traverses the file system starting from the provided root path.
A filter is provided to only select specific metadata files, whose naming follows a reserved global setting such as `MANIFEST.yml`. Selected files are then marshalled and returned using the `abstract.Asset` definition:

```go
type Asset struct {
	// asset last found by crawler at - only added by service (not crawler/manifest itself, i.e. no yaml)
	LastDiscoveredAt time.Time `json:"last-discovered-at"`
	// asset publication datetime
	PublishedOn time.Time `yaml:"published-on" json:"published-on"`
	// name of the asset
	Name string `yaml:"name" json:"name"`
	// description of the asset
	Description string `yaml:"description" json:"description"`
	// the list of assets this depends on
	DependsOn []string `yaml:"depends-on" json:"depends-on"`
	// asset type
	Type AssetType `yaml:"type,omitempty" json:"type,omitempty"`
	// labels for the specific asset
	Labels map[string]interface{} `yaml:"labels,omitempty" json:"labels,omitempty"`
	// tags are flags used to simplify asset search
	Tags []string `yaml:"tags,omitempty" json:"tags,omitempty"`
}
```

The package also provide means to parse and validate assets:
```go
func ParseAsset(data []byte) (*Asset, error) {}
func (asset *Asset) Validate() error {}
```