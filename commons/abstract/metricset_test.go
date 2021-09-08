package abstract

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricSetParsing(t *testing.T) {

	// actual metrics
	inputJson := `
	{
		"name" : "test",
		"inserted_at" : "2020-11-29T17:24:01.747543Z",	
		"version" : "dq_pipeline_test_env",
		"description" : "my first metricset",
		"labels" : {
			"year" : "2021",
			"month" : "09",
			"day" : "13"
		},
		"metrics" : [
			{
				"resultKey": {
					"dataSetDate": 1630876393300,
					"tags": {}
				},
				"analyzerContext": {
					"metricMap": [
						{
						"analyzer": {
							"analyzerName": "Size"
						},
						"metric": {
							"metricName": "DoubleMetric",
							"entity": "Dataset",
							"instance": "*",
							"name": "Size",
							"value": 5.0
						}
						},
						{
						"analyzer": {
							"analyzerName": "Minimum",
							"column": "numViews"
						},
						"metric": {
							"metricName": "DoubleMetric",
							"entity": "Column",
							"instance": "numViews",
							"name": "Minimum",
							"value": 0.0
						}
						}
					]
				}
			}
		]
	}
	`
	metricset := MetricSet{}
	err := json.Unmarshal([]byte(inputJson), &metricset)

	assert := assert.New(t)
	assert.Equal(err, nil)
	assert.NotEqual(metricset, nil)
	assert.Equal(1, len(metricset.Metrics))
	assert.Equal(metricset.Metrics[0].ResultKey.DataSetDate, int64(1630876393300))
}
