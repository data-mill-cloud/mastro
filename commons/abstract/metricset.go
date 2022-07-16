package abstract

import (
	"errors"
	"fmt"
	"strings"
	"time"

	resterrors "github.com/data-mill-cloud/mastro/commons/utils/errors"

	"github.com/data-mill-cloud/mastro/commons/utils/conf"
)

// MetricSet ... a timestamped set of Metrics
type MetricSet struct {
	// unique metric set name
	Name string `json:"name,omitempty"`
	// insertion time
	InsertedAt time.Time `json:"inserted_at,omitempty"`
	// version relates to the environment and the pipeline version
	Version string `json:"version,omitempty"`
	// description is related to the metrics and their extraction process and not the datasource they were calculated on
	Description string `json:"description,omitempty"`
	// labels used for query purposes
	Labels map[string]string `json:"labels,omitempty"`

	// actual metrics
	Metrics []Metric `json:"metrics,omitempty"`
}

type Metric struct {
	*DeequMetric
}

// Validate ... validate a metricSet
func (ms *MetricSet) Validate() error {
	// the name should not be empty or we may not be able to retrieve the mset
	if len(strings.TrimSpace(ms.Name)) == 0 {
		return errors.New("MetricSet Name is undefined")
	}

	if len(strings.TrimSpace(ms.Version)) == 0 {
		return errors.New("MetricSet Version is undefined")
	}

	for i, f := range ms.Metrics {
		if err := f.Validate(); err != nil {
			return errors.New(fmt.Sprintf("%s at position %d", err.Error(), i))
		}
	}

	return nil
}

// Validate ... validate a metric
func (m *Metric) Validate() error {
	// make sure at least one metric type is defined
	if m.DeequMetric == nil {
		return errors.New("Incorrect metric definition")
	}

	return nil
}

// MetricSetDAOProvider ... The interface each dao must implement
type MetricSetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(m *MetricSet) error
	GetById(id string) (*MetricSet, error)
	GetByName(name string, limit int, page int) (*Paginated[MetricSet], error)
	SearchMetricSetsByLabels(labels map[string]string, limit int, page int) (*Paginated[MetricSet], error)
	ListAllMetricSets(limit int, page int) (*Paginated[MetricSet], error)
	Search(query string, limit int, page int) (*Paginated[MetricSet], error)
	CloseConnection()
}

// MetricStoreService ... MetricStoreService Interface listing service methods
type MetricStoreService interface {
	Init(cfg *conf.Config) *resterrors.RestErr
	CreateMetricSet(ms MetricSet) (*MetricSet, *resterrors.RestErr)
	GetMetricSetByID(msID string) (*MetricSet, *resterrors.RestErr)
	GetMetricSetByName(msName string, limit int, page int) (*Paginated[MetricSet], *resterrors.RestErr)
	SearchMetricSetsByLabels(labels map[string]string, limit int, page int) (*Paginated[MetricSet], *resterrors.RestErr)
	Search(query string, limit int, page int) (*Paginated[MetricSet], *resterrors.RestErr)
	ListAllMetricSets(limit int, page int) (*Paginated[MetricSet], *resterrors.RestErr)
}
