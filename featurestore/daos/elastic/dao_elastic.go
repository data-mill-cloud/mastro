package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/datamillcloud/mastro/commons/abstract"
	"github.com/datamillcloud/mastro/commons/sources/elastic"
	"github.com/datamillcloud/mastro/commons/utils/conf"
	"github.com/elastic/go-elasticsearch/esapi"
)

// both init and sync.Once are thread-safe
// but only sync.Once is lazy
var once sync.Once
var instance *dao

// GetSingleton ... get an instance of the dao backend
func GetSingleton() abstract.FeatureSetDAOProvider {
	// once.do is lazy, we use it to return an instance of the DAO
	once.Do(func() {
		instance = &dao{}
	})
	return instance
}

// dao ... The struct for the ElasticSearch DAO for the FeatureStore service
type dao struct {
	Connector *elastic.Connector
}

/*
{
  "took" : 1,
  "timed_out" : false,
  "_shards" : {
    "total" : 1,
    "successful" : 1,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : {
      "value" : 0,
      "relation" : "eq"
    },
    "max_score" : null,
    "hits" : [ ]
  }
}
*/
type SearchResponse struct {
	Took     float64 `json:"took,omitempty"`
	TimedOut bool    `json:"timed_out,omitempty"`
	Shards   Shards  `json:"_shards,omitempty"`
	Hits     Hits    `json:"hits,omitempty"`
}

type Shards struct {
	Total       float64 `json:"total,omitempty"`
	Successfull float64 `json:"successfull,omitempty"`
	Skipped     float64 `json:"skipped,omitempty"`
	Failed      float64 `json:"failed,omitempty"`
}
type Hits struct {
	Total    Total         `json:"total,omitempty"`
	MaxScore float64       `json:"max_score,omitempty"`
	Hits     []ResponseDoc `json:"hits,omitempty"`
}

type Total struct {
	Value    float64 `json:"value,omitempty"`
	Relation string  `json:"relation,omitempty"`
}

/*
},
      {
        "_index" : "test",
        "_type" : "_doc",
        "_id" : "q5QpE3YBC-0pshNrv60c",
        "_score" : 1.0,
        "_source" : {
          "settings" : {
            "number_of_shards" : 1,
            "number_of_replicas" : 0
          }
        }
      }
*/
type ResponseDoc struct {
	Index  string     `json:"_index,omitempty"`
	Type   string     `json:"_type,omitempty"`
	ID     string     `json:"_id,omitempty"`
	Score  float64    `json:"_score,omitempty"`
	Source FeatureSet `json:"_source,omitempty"`
}

// FeatureSet ... a versioned set of features
type FeatureSet struct {
	InsertedAt  time.Time         `json:"inserted_at,omitempty"`
	Version     string            `json:"version,omitempty"`
	Features    []Feature         `json:"features,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// Version ... definition of version for a feature set
type Version struct{}

// Feature ... a named variable with a data type
type Feature struct {
	Name     string `json:"name,omitempty"`
	Value    string `json:"value,omitempty"`
	DataType string `json:"data-type,omitempty"`
}

func (dao *dao) checkIndex(def *conf.DataSourceDefinition) error {
	// if file location has no specified folder, check at the same position of the config
	indexDefFilePath := def.Settings["index-def"]

	// if the def is indicated without an actual path (parent is itself and dir is ., or directly as ./file)
	if filepath.Dir(indexDefFilePath) == "." || filepath.Dir(indexDefFilePath) == "." {
		// look for the index def in the same location of the application config
		indexDefFilePath = filepath.Join(filepath.Dir(conf.Args.Config), indexDefFilePath)
	}
	log.Println("Attempting loading index def file from folder", indexDefFilePath)

	// read definition from file
	defFile, err := ioutil.ReadFile(indexDefFilePath)
	if err != nil {
		return err
	}

	// Instantiate an indexRequest object
	req := esapi.IndexRequest{
		Index:   dao.Connector.IndexName,
		Body:    strings.NewReader(string(defFile)),
		Refresh: "true",
	}

	// Return an API response object from request
	ctx := context.Background()
	res, err := req.Do(ctx, dao.Connector.Client)
	if err != nil {
		return fmt.Errorf("IndexRequest ERROR: %s", err)
	}
	defer res.Body.Close()
	return nil
}

// Init ... Initialize connection to elastic search and target index
func (dao *dao) Init(def *conf.DataSourceDefinition) {
	// create connector
	dao.Connector = elastic.NewElasticConnector()
	// validate data source definition
	if err := dao.Connector.ValidateDataSourceDefinition(def); err != nil {
		panic(err)
	}
	// init connector
	dao.Connector.InitConnection(def)
	// make sure the target index exists
	//err := dao.checkIndex(def)
	//if err != nil {
	//	log.Fatalln(err)
	//}
}

// Create ... Create featureset on ES
func (dao *dao) Create(fs *abstract.FeatureSet) error {

	jsonVal, err := json.Marshal(fs)
	if err != nil {
		return err
	}

	body := string(jsonVal)
	// Instantiate an indexRequest object
	req := esapi.IndexRequest{
		Index: dao.Connector.IndexName,
		// https://www.elastic.co/guide/en/elasticsearch/reference/6.8/mapping-id-field.html
		//DocumentID: strconv.Itoa(i),
		Body:    strings.NewReader(body),
		Refresh: "true",
	}

	// Return an API response object from request
	ctx := context.Background()
	res, err := req.Do(ctx, dao.Connector.Client)
	if err != nil {
		return fmt.Errorf("IndexRequest ERROR: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Println(res.String())
		return fmt.Errorf("%s ERROR indexing document ", res.Status())
	}

	return nil
}

func (dao *dao) search(buf *bytes.Buffer) (*SearchResponse, error) {
	// Perform a search request.
	res, err := dao.Connector.Client.Search(
		dao.Connector.Client.Search.WithContext(context.Background()),
		dao.Connector.Client.Search.WithIndex(dao.Connector.IndexName),
		dao.Connector.Client.Search.WithBody(buf),
		dao.Connector.Client.Search.WithTrackTotalHits(true),
		dao.Connector.Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("Error parsing the response body: %s", err)
		}
		// Print the response status and error information.
		return nil, fmt.Errorf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	/*
		Example reply is of type:
		{
		"took" : 2,
		"timed_out" : false,
		"_shards" : {
			"total" : 1,
			"successful" : 1,
			"skipped" : 0,
			"failed" : 0
		},
		"hits" : {
			"total" : {
			"value" : 0,
			"relation" : "eq"
			},
			"max_score" : null,
			"hits" : [ ]
		}
		}
	*/

	bodyBuf := new(bytes.Buffer)
	bodyBuf.ReadFrom(res.Body)
	//newStr := bodyBuf.String()

	searchResponse := &SearchResponse{}
	if err := json.Unmarshal(bodyBuf.Bytes(), &searchResponse); err != nil {
		return nil, err
	}

	return searchResponse, nil
}

// GetById ... Retrieve document by given id
func (dao *dao) GetById(id string) (*abstract.FeatureSet, error) {
	// prepare search query
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-ids-query.html
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"ids": map[string]interface{}{
				"values": []string{
					id,
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("Error encoding query: %s", err)
	}

	searchResponse, err := dao.search(&buf)
	if err != nil {
		return nil, err
	}

	log.Println("GetById :: Retrieved", searchResponse.Hits.Total.Value, "documents")
	if searchResponse.Hits.Total.Value > 0 {
		hitDocs, err := convertDocumentsToFeatureSetCollection(searchResponse.Hits.Hits)
		if err != nil {
			return nil, err
		}
		return &((*hitDocs)[0]), nil
	}
	// else return an empty feature set
	return nil, fmt.Errorf("No document found for id %s", id)
}

// GetByName ... Retrieve document by given name
func (dao *dao) GetByName(name string) (*[]abstract.FeatureSet, error) {

	var buf bytes.Buffer
	// use a term query to do an exact match of the name
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-term-query.html
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"name": name,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("Error encoding query: %s", err)
	}

	searchResponse, err := dao.search(&buf)
	if err != nil {
		return nil, err
	}

	log.Println("GetByName :: Retrieved", searchResponse.Hits.Total.Value, "documents")
	if searchResponse.Hits.Total.Value > 0 {
		hitDocs, err := convertDocumentsToFeatureSetCollection(searchResponse.Hits.Hits)
		if err != nil {
			return nil, err
		}
		//return &((*hitDocs)[0]), nil
		return hitDocs, nil
	}
	// else return an empty feature set
	return nil, fmt.Errorf("No document found for name %s", name)
}

// ListAllFeatureSets ... Return all featuresets in index
func (dao *dao) ListAllFeatureSets() (*[]abstract.FeatureSet, error) {

	// in ES to return all documents in an index we need the following query:
	/*
		{
		    "query": {
		        "match_all": {}
		    }
		}
	*/

	// prepare search query
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("Error encoding query: %s", err)
	}

	searchResponse, err := dao.search(&buf)
	if err != nil {
		return nil, err
	}

	log.Println("ListAllFeatureSets :: Retrieved", searchResponse.Hits.Total.Value, "documents")
	if searchResponse.Hits.Total.Value > 0 {
		fsColl, err := convertDocumentsToFeatureSetCollection(searchResponse.Hits.Hits)
		if err != nil {
			return nil, err
		}
		return fsColl, nil
	}
	// else return an empty feature set
	return nil, fmt.Errorf("No document found in index %s", dao.Connector.IndexName)
}

func convertDocumentsToFeatureSetCollection(documents []ResponseDoc) (*[]abstract.FeatureSet, error) {
	featureSetCollection := []abstract.FeatureSet{}
	for _, d := range documents {
		fs, err := convertDocumentToFeatureSet(d)
		if err != nil {
			return nil, err
		}
		featureSetCollection = append(featureSetCollection, *fs)
	}
	return &featureSetCollection, nil
}

func convertDocumentToFeatureSet(document ResponseDoc) (*abstract.FeatureSet, error) {
	fs := abstract.FeatureSet{}
	fs.InsertedAt = document.Source.InsertedAt
	fs.Version = document.Source.Version
	features, err := convertDaoFeaturesToFeatures(document.Source.Features)
	if err != nil {
		return nil, err
	}
	fs.Features = *features
	fs.Description = document.Source.Description
	fs.Labels = document.Source.Labels
	return &fs, nil
}

func convertDaoFeaturesToFeatures(features []Feature) (*[]abstract.Feature, error) {
	result := []abstract.Feature{}
	for _, f := range features {
		af := abstract.Feature{}
		af.Name = f.Name
		af.DataType = f.DataType
		switch af.DataType {
		case "bool":
			b, err := strconv.ParseBool(f.Value)
			if err != nil {
				return nil, err
			}
			af.Value = b
		case "int":
			n, err := strconv.ParseInt(f.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			af.Value = n
		case "float":
			f, err := strconv.ParseFloat(f.Value, 64)
			if err != nil {
				return nil, err
			}
			af.Value = f
		case "string":
			af.Value = f.Value
		}
		result = append(result, af)
	}
	return &result, nil
}

// CloseConnection ... Terminates the connection to ES for the DAO
func (dao *dao) CloseConnection() {
	dao.Connector.CloseConnection()
}
