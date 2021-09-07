package abstract

// DeequMetric ...
type DeequMetric struct {
	ResultKey       ResultKey       `json:"resultKey"`
	AnalyzerContext AnalyzerContext `json:"analyzerContext"`
}

type ResultKey struct {
	DataSetDate int64             `json:"dataSetDate"`
	Tags        map[string]string `json:"tags"`
}

type AnalyzerContext struct {
	MetricMap []MetricInstance `json:"metricMap"`
}

type MetricInstance struct {
	Analyzer Analyzer `json:"analyzer"`
	Metric   Metric   `json:"metric"`
}

type Analyzer struct {
	AnalyzerName string `json:"analyzerName"`
	Column       string `json:"column"`
}

type Metric struct {
	MetricName string  `json:"metricName"`
	Entity     string  `json:"entity"`
	Instance   string  `json:"instance"`
	Name       string  `json:"name"`
	Value      float64 `json:"value"`
}
