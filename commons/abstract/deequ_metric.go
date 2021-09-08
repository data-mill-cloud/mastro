package abstract

// DeequMetric ...
type DeequMetric struct {
	ResultKey       DeequResultKey       `json:"resultKey"`
	AnalyzerContext DeequAnalyzerContext `json:"analyzerContext"`
}

type DeequResultKey struct {
	DataSetDate int64             `json:"dataSetDate"`
	Tags        map[string]string `json:"tags"`
}

type DeequAnalyzerContext struct {
	MetricMap []DeequMetricInstance `json:"metricMap"`
}

type DeequMetricInstance struct {
	Analyzer DeequAnalyzer    `json:"analyzer"`
	Metric   DeequMetricValue `json:"metric"`
}

type DeequAnalyzer struct {
	AnalyzerName string `json:"analyzerName"`
	Column       string `json:"column"`
}

type DeequMetricValue struct {
	MetricName string  `json:"metricName"`
	Entity     string  `json:"entity"`
	Instance   string  `json:"instance"`
	Name       string  `json:"name"`
	Value      float64 `json:"value"`
}
