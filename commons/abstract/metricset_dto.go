package abstract

import (
	"errors"
	"strings"
	"time"
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
	DeequMetric
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

	for _, f := range ms.Metrics {
		if err := f.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate ... validate a metric
func (m *Metric) Validate() error {
	return nil
}
