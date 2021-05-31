package abstract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsing(t *testing.T) {

	inputYaml := `
    published-on: "2015-08-06T17:52:48Z"
    name: "testAsset"
    description: "this is an example asset"
    depends-on: ["asset1", "someother"]
    type: dataset
    labels:
      key1: "val1"
      key2: "val2"
    tags:
      - testtag`

	asset, err := ParseAsset([]byte(inputYaml))

	assert := assert.New(t)
	assert.Equal(err, nil)
	assert.NotEqual(asset, nil)

  t.Log(asset)
}
