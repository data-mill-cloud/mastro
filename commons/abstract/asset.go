package abstract

import (
	"errors"
	"fmt"
	"strings"
	"time"

	resterrors "github.com/data-mill-cloud/mastro/commons/utils/errors"

	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"gopkg.in/yaml.v2"
)

// Asset ... managed resource
type Asset struct {
	// asset last found by crawler at - only added by service (not crawler/manifest itself, i.e. no yaml)
	LastDiscoveredAt time.Time `json:"last-discovered-at" yaml:"last-discovered-at,omitempty"`
	// asset publication datetime
	PublishedOn time.Time `yaml:"published-on" json:"published-on"`
	// name of the asset
	Name string `yaml:"name" json:"name"`
	// description of the asset
	Description string `yaml:"description" json:"description"`
	// the list of assets this depends on
	DependsOn []string `yaml:"depends-on" json:"depends-on"`
	// asset type
	Type AssetType `yaml:"type" json:"type"`
	// labels for the specific asset
	Labels map[string]interface{} `yaml:"labels" json:"labels"`
	// tags are flags used to simplify asset search
	Tags []string `yaml:"tags" json:"tags"`
	// versions specify available variants of the same asset
	Versions map[string]interface{} `yaml:"versions" json:"versions"`
}

// AssetType ... Asset type information
type AssetType string

const (
	_Database   AssetType = "database"
	_Dataset              = "dataset"
	_FeatureSet           = "featureset"
	_Embedding            = "embedding"
	_Model                = "model"
	_Notebook             = "notebook"
	_Pipeline             = "pipeline"
	_Report               = "report"
	_Service              = "service"
	_Stream               = "stream"
	_Table                = "table"
	_User                 = "user"
	_Workflow             = "workflow"
)

func NewDatasetAsset() *Asset {
	asset := &Asset{}
	asset.PublishedOn = time.Now()
	asset.Type = _Dataset

	return asset
}

var assetTypes = []AssetType{
	_Database,
	_Dataset,
	_FeatureSet,
	_Embedding,
	_Model,
	_Notebook,
	_Pipeline,
	_Report,
	_Service,
	_Stream,
	_Table,
	_User,
	_Workflow,
}

func isValidType(t AssetType) bool {
	for _, b := range assetTypes {
		if t == b {
			return true
		}
	}
	return false
}

// ParseAsset ... Parse an asset specification file
func ParseAsset(data []byte) (*Asset, error) {
	asset := Asset{}

	err := yaml.Unmarshal(data, &asset)

	return &asset, err
}

func (assetType *AssetType) Validate() error {
	inputStr := strings.TrimSpace(string(*assetType))
	if len(inputStr) == 0 {
		return errors.New("AssetType is empty")
	}

	// Validate the valid enum values
	if !isValidType(*assetType) {
		return errors.New(fmt.Sprintf("invalid value %s for AssetType", inputStr))
	}
	return nil
}

// Validate ... Validate asset specification file
func (asset *Asset) Validate() error {

	// validate required fields for an asset
	// - name
	// - assetType
	if len(strings.TrimSpace(asset.Name)) == 0 {
		return errors.New("Name is undefined")
	}

	if err := asset.Type.Validate(); err != nil {
		return err
	}

	// validate optional fields if any available

	return nil
}

// Label types

const (
	L_SCHEMA = "schema"
)

// AssetDAOProvider ... The interface each dao must implement
type AssetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Upsert(asset *Asset) error
	GetById(id string) (*Asset, error)
	GetByName(id string) (*Asset, error)
	SearchAssetsByTags(tags []string, limit int, page int) (*Paginated[Asset], error)
	ListAllAssets(limit int, page int) (*Paginated[Asset], error)
	Search(query string, limit int, page int) (*Paginated[Asset], error)
	CloseConnection()
}

// CatalogueService ... CatalogueService Interface listing service methods
type CatalogueService interface {
	Init(cfg *conf.Config) *resterrors.RestErr
	UpsertAssets(assets *[]Asset) (*[]Asset, *resterrors.RestErr)
	GetAssetByID(assetID string) (*Asset, *resterrors.RestErr)
	GetAssetByName(name string) (*Asset, *resterrors.RestErr)
	SearchAssetsByTags(tags []string, limit int, page int) (*Paginated[Asset], *resterrors.RestErr)
	Search(query string, limit int, page int) (*Paginated[Asset], *resterrors.RestErr)
	ListAllAssets(limit int, page int) (*Paginated[Asset], *resterrors.RestErr)
}
