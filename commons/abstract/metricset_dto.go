package abstract

import "time"

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
