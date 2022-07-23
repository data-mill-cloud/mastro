package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	stringutils "github.com/data-mill-cloud/mastro/commons/utils/strings"

	//es7 "github.com/elastic/go-elasticsearch/v7"

	"github.com/elastic/go-elasticsearch/esapi"
	es "github.com/elastic/go-elasticsearch/v8"
)

// NewElasticConnector factory
func NewElasticConnector[T any]() *Connector[T] {
	return &Connector[T]{
		ConfigurableConnector: abstract.ConfigurableConnector{
			RequiredFields: map[string]string{
				"esUser":  "username",
				"esPwd":   "password",
				"esHosts": "hosts",
				"esIndex": "index",
			},
			OptionalFields: map[string]string{
				"cert": "cert",
			},
		},
	}
}

/*
https://www.elastic.co/blog/the-go-client-for-elasticsearch-working-with-data
*/

// todo: find a way not to export this

// Connector ... Connector type
type Connector[T any] struct {
	abstract.ConfigurableConnector
	Client    *es.Client
	IndexName string
}

// InitConnection ... Starts a connection with Elastic Search
func (c *Connector[T]) InitConnection(def *conf.DataSourceDefinition) {
	var err error
	//c.client, err = es7.NewDefaultClient()
	elasticHostnames := stringutils.SplitAndTrim(def.Settings[c.RequiredFields["esHosts"]], ",")

	esConfig := es.Config{
		Addresses: elasticHostnames,
		Username:  def.Settings[c.RequiredFields["esUser"]],
		Password:  def.Settings[c.RequiredFields["esPwd"]],
	}
	// if encryption is enabled then set the server certificate
	if certFile, exist := def.Settings[c.OptionalFields["cert"]]; exist {
		cert, err := ioutil.ReadFile(certFile)
		if err != nil {
			log.Fatal("Error while reading certificate", err)
		}
		esConfig.CACert = cert
	}

	c.Client, err = es.NewClient(esConfig)
	// set the index for the client
	c.IndexName = def.Settings[c.RequiredFields["esIndex"]]

	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	log.Println("Successfully connected to ES")

	// make sure the target index exists
	err = c.CheckIndex(def)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Connector[T]) CheckIndex(def *conf.DataSourceDefinition) error {
	response, err := c.Client.Indices.Exists([]string{c.IndexName})
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusNotFound {
		log.Println("Index already exists, skipping creation")
		return nil
	}
	/*
		response, err = c.Client.Indices.Create(c.IndexName)
		if err != nil || response.IsError() {
			return err
		}
	*/

	// if file location has no specified folder, check at the same position of the config
	indexDefFilePath := def.Settings["index-def"]

	// if the def is indicated without an actual path (parent is itself and dir is ., or directly as ./file)
	if filepath.Dir(indexDefFilePath) == "." {
		// look for the index def in the same location of the application config
		indexDefFilePath = filepath.Join(filepath.Dir(conf.Args.Config), indexDefFilePath)
	}
	log.Println("Attempting loading index def file from folder", indexDefFilePath)

	// read definition from file
	defFile, err := ioutil.ReadFile(indexDefFilePath)
	if err != nil {
		return err
	}

	indexReq := esapi.IndicesCreateRequest{
		Index: c.IndexName,
		Body:  strings.NewReader(string(defFile)),
	}
	ctx := context.Background()
	res, err := indexReq.Do(ctx, c.Client)
	if err != nil {
		return fmt.Errorf("IndexRequest ERROR: %s", err)
	}
	defer res.Body.Close()

	return nil
}

type SearchResponse[T any] struct {
	Took     float64 `json:"took,omitempty"`
	TimedOut bool    `json:"timed_out,omitempty"`
	Shards   Shards  `json:"_shards,omitempty"`
	Hits     Hits[T] `json:"hits,omitempty"`
}

type Shards struct {
	Total       float64 `json:"total,omitempty"`
	Successfull float64 `json:"successfull,omitempty"`
	Skipped     float64 `json:"skipped,omitempty"`
	Failed      float64 `json:"failed,omitempty"`
}
type Hits[T any] struct {
	Total    Total            `json:"total,omitempty"`
	MaxScore float64          `json:"max_score,omitempty"`
	Hits     []ResponseDoc[T] `json:"hits,omitempty"`
}

type Total struct {
	Value    float64 `json:"value,omitempty"`
	Relation string  `json:"relation,omitempty"`
}

type ResponseDoc[T any] struct {
	Index  string  `json:"_index,omitempty"`
	Type   string  `json:"_type,omitempty"`
	ID     string  `json:"_id,omitempty"`
	Score  float64 `json:"_score,omitempty"`
	Source T       `json:"_source,omitempty"`
}

func (c *Connector[T]) Delete(id string) error {
	ctx := context.Background()
	res, err := c.Client.Delete(
		c.IndexName,
		id,
		c.Client.Delete.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}
	if res.IsError() {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		return fmt.Errorf("error getting response: %s", buf.String())
	}
	defer res.Body.Close()
	return nil
}

func (c *Connector[T]) DeleteByQuery(body *bytes.Buffer) error {
	ctx := context.Background()
	res, err := c.Client.DeleteByQuery(
		[]string{c.IndexName},
		body,
		c.Client.DeleteByQuery.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()
	return err
}

func (c *Connector[T]) Search(buf *bytes.Buffer) (*SearchResponse[T], error) {
	// Perform a search request.
	res, err := c.Client.Search(
		c.Client.Search.WithContext(context.Background()),
		c.Client.Search.WithIndex(c.IndexName),
		c.Client.Search.WithBody(buf),
		c.Client.Search.WithTrackTotalHits(true),
		c.Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %s", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		}
		// Print the response status and error information.
		return nil, fmt.Errorf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	bodyBuf := new(bytes.Buffer)
	bodyBuf.ReadFrom(res.Body)
	//newStr := bodyBuf.String()

	searchResponse := &SearchResponse[T]{}
	if err := json.Unmarshal(bodyBuf.Bytes(), &searchResponse); err != nil {
		return nil, err
	}

	return searchResponse, nil
}

func (c *Connector[T]) SimilarToThis(vectorFieldName string, vector []float32, k int, numCandidates int, projectionFields []string, filter *map[string]interface{}) (*SearchResponse[T], error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"knn": map[string]interface{}{
			"field":          vectorFieldName,
			"query_vector":   vector,
			"k":              k,
			"num_candidates": numCandidates,
		},
	}

	if projectionFields != nil {
		query["fields"] = projectionFields
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %s", err)
	}

	res, err := c.Client.KnnSearch(
		[]string{c.IndexName},
		c.Client.KnnSearch.WithContext(context.Background()),
		c.Client.KnnSearch.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %s", err)
	}

	if err != nil {
		return nil, fmt.Errorf("error getting response: %s", err)
	}
	if res.IsError() {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		return nil, fmt.Errorf("error getting response: %s", buf.String())
	}
	defer res.Body.Close()

	bodyBuf := new(bytes.Buffer)
	bodyBuf.ReadFrom(res.Body)

	searchResponse := &SearchResponse[T]{}
	if err := json.Unmarshal(bodyBuf.Bytes(), &searchResponse); err != nil {
		return nil, err
	}

	return searchResponse, nil
}

func (c *Connector[T]) CloseConnection() {

}
