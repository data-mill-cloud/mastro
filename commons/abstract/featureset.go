package abstract

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	resterrors "github.com/data-mill-cloud/mastro/commons/utils/errors"

	"github.com/data-mill-cloud/mastro/commons/utils/conf"
)

// FeatureSet ... a versioned set of features
type FeatureSet struct {
	Name        string            `json:"name,omitempty"`
	InsertedAt  time.Time         `json:"inserted_at,omitempty"`
	Version     string            `json:"version,omitempty"`
	Features    []Feature         `json:"features,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// Feature ... a named variable with a data type
type Feature struct {
	Name     string      `json:"name,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	DataType string      `json:"data_type,omitempty"`
}

// Validate ... validate a featureSet
func (fs *FeatureSet) Validate() error {
	// the name should not be empty or we may not be able to retrieve the fset
	if len(strings.TrimSpace(fs.Name)) == 0 {
		return errors.New("FeatureSet Name is undefined")
	}

	if len(strings.TrimSpace(fs.Version)) == 0 {
		return errors.New("FeatureSet Version is undefined")
	}

	for _, f := range fs.Features {
		if err := f.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate ... validate a feature
func (f *Feature) Validate() error {
	if len(strings.TrimSpace(f.Name)) == 0 {
		return errors.New("Feature Name is undefined")
	}

	log.Println(f.Name, f.Value)
	if f.Value == nil {
		return errors.New(fmt.Sprintf("Feature Value for Feature %s is undefined", f.Name))
	}

	if len(strings.TrimSpace(f.DataType)) == 0 {
		return errors.New(fmt.Sprintf("Feature Data Type for Feature %s is undefined", f.Name))
	}

	return nil
}

// FeatureSetDAOProvider ... The interface each dao must implement
type FeatureSetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(fs *FeatureSet) error
	GetById(id string) (*FeatureSet, error)
	GetByName(name string, limit int, page int) (*Paginated[FeatureSet], error)
	SearchFeatureSetsByLabels(labels map[string]string, limit int, page int) (*Paginated[FeatureSet], error)
	Search(query string, limit int, page int) (*Paginated[FeatureSet], error)
	ListAllFeatureSets(limit int, page int) (*Paginated[FeatureSet], error)
	CloseConnection()
}

// Service ... FeatureStoreService Interface listing implemented methods
type FeatureStoreService interface {
	Init(cfg *conf.Config) *resterrors.RestErr
	CreateFeatureSet(fs FeatureSet) (*FeatureSet, *resterrors.RestErr)
	GetFeatureSetByID(fsID string) (*FeatureSet, *resterrors.RestErr)
	GetFeatureSetByName(fsName string, limit int, page int) (*Paginated[FeatureSet], *resterrors.RestErr)
	SearchFeatureSetsByLabels(labels map[string]string, limit int, page int) (*Paginated[FeatureSet], *resterrors.RestErr)
	Search(query string, limit int, page int) (*Paginated[FeatureSet], *resterrors.RestErr)
	ListAllFeatureSets(limit int, page int) (*Paginated[FeatureSet], *resterrors.RestErr)
}
