package abstract

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeequMetricParsing(t *testing.T) {

	inputJson := `
	[
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
	]`

	metrics := []DeequMetric{}
	err := json.Unmarshal([]byte(inputJson), &metrics)

	assert := assert.New(t)
	assert.Equal(err, nil)
	assert.NotEqual(metrics, nil)
	assert.Equal(len(metrics), 1)
	assert.Equal(metrics[0].ResultKey.DataSetDate, int64(1630876393300))

	t.Log(metrics)
}
