package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/sources/elastic"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
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
	Connector *elastic.Connector[FeatureSet]
}

// FeatureSet ... a versioned set of features
type FeatureSet struct {
	Name        string            `json:"name,omitempty"`
	InsertedAt  time.Time         `json:"inserted_at,omitempty"`
	Version     string            `json:"version,omitempty"`
	Features    []Feature         `json:"features,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// Feature ... a named variable with a data type
type Feature struct {
	Name     string `json:"name,omitempty"`
	Value    string `json:"value,omitempty"`
	DataType string `json:"data-type,omitempty"`
}

// Init ... Initialize connection to elastic search and target index
func (dao *dao) Init(def *conf.DataSourceDefinition) {
	// create connector
	dao.Connector = elastic.NewElasticConnector[FeatureSet]()
	// validate data source definition
	if err := dao.Connector.ValidateDataSourceDefinition(def); err != nil {
		panic(err)
	}
	// init connector
	dao.Connector.InitConnection(def)
}

func convertDtoToDao(fs *abstract.FeatureSet) (result *FeatureSet, err error) {
	features := []Feature{}
	var b []byte
	for _, f := range fs.Features {

		if b, err = json.Marshal(f.Value); err != nil {
			return
		}

		features = append(features, Feature{
			Name:     f.Name,
			Value:    string(b),
			DataType: f.DataType,
		})
	}

	result = &FeatureSet{
		Name:        fs.Name,
		InsertedAt:  fs.InsertedAt,
		Version:     fs.Version,
		Features:    features,
		Description: fs.Description,
		Labels:      fs.Labels,
	}

	return
}

// Create ... Create featureset on ES
func (dao *dao) Create(fs *abstract.FeatureSet) error {
	daoFs, err := convertDtoToDao(fs)
	if err != nil {
		return err
	}

	jsonVal, err := json.Marshal(daoFs)
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
		return nil, fmt.Errorf("error encoding query: %s", err)
	}

	searchResponse, err := dao.Connector.Search(&buf)
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
	return nil, fmt.Errorf("no document found for id %s", id)
}

// GetByName ... Retrieve document by given name
func (dao *dao) GetByName(name string, limit int, page int) (*abstract.Paginated[abstract.FeatureSet], error) {
	// todo: add paging using limit and page params

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
		return nil, fmt.Errorf("error encoding query: %s", err)
	}

	searchResponse, err := dao.Connector.Search(&buf)
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
		return &abstract.Paginated[abstract.FeatureSet]{Data: hitDocs}, nil
	}
	// else return an empty feature set
	return nil, fmt.Errorf("no document found for name %s", name)
}

// ListAllFeatureSets ... Return all featuresets in index
func (dao *dao) ListAllFeatureSets(limit int, page int) (*abstract.Paginated[abstract.FeatureSet], error) {
	// prepare search query
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %s", err)
	}

	searchResponse, err := dao.Connector.Search(&buf)
	if err != nil {
		return nil, err
	}

	log.Println("ListAllFeatureSets :: Retrieved", searchResponse.Hits.Total.Value, "documents")
	if searchResponse.Hits.Total.Value > 0 {
		fsColl, err := convertDocumentsToFeatureSetCollection(searchResponse.Hits.Hits)
		if err != nil {
			return nil, err
		}
		return &abstract.Paginated[abstract.FeatureSet]{Data: fsColl}, nil
	}
	// else return an empty feature set
	return nil, fmt.Errorf("no document found in index %s", dao.Connector.IndexName)
}

func (dao *dao) Search(query string, limit int, page int) (*abstract.Paginated[abstract.FeatureSet], error) {

	// match all documents matching any of the values specified in the search field
	var buf bytes.Buffer
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"description": query,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
		return nil, fmt.Errorf("error encoding query: %s", err)
	}

	searchResponse, err := dao.Connector.Search(&buf)
	if err != nil {
		return nil, err
	}

	log.Println("Search :: Retrieved", searchResponse.Hits.Total.Value, "documents")
	if searchResponse.Hits.Total.Value > 0 {
		fsColl, err := convertDocumentsToFeatureSetCollection(searchResponse.Hits.Hits)
		if err != nil {
			return nil, err
		}
		return &abstract.Paginated[abstract.FeatureSet]{Data: fsColl}, nil
	}
	// else return an empty feature set
	return nil, fmt.Errorf("no document found in index %s", dao.Connector.IndexName)
}

func (dao *dao) SearchFeatureSetsByLabels(labels map[string]string, limit int, page int) (*abstract.Paginated[abstract.FeatureSet], error) {

	matchLabels := make([]map[string]interface{}, 0)
	for k, v := range labels {
		matchLabels = append(matchLabels, map[string]interface{}{
			"match": map[string]interface{}{
				"labels." + k: v,
			},
		})
	}

	// match all documents matching any of the values specified in the search field
	var buf bytes.Buffer
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"nested": map[string]interface{}{
				"path":       "labels",
				"score_mode": "avg",
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"must": matchLabels,
					},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
		return nil, fmt.Errorf("error encoding query: %s", err)
	}

	searchResponse, err := dao.Connector.Search(&buf)
	if err != nil {
		return nil, err
	}

	log.Println("SearchFeatureSetsByLabels :: Retrieved", searchResponse.Hits.Total.Value, "documents")
	if searchResponse.Hits.Total.Value > 0 {
		fsColl, err := convertDocumentsToFeatureSetCollection(searchResponse.Hits.Hits)
		if err != nil {
			return nil, err
		}
		return &abstract.Paginated[abstract.FeatureSet]{Data: fsColl}, nil
	}
	// else return an empty feature set
	return nil, fmt.Errorf("no document found in index %s", dao.Connector.IndexName)
}

func convertDocumentsToFeatureSetCollection(documents []elastic.ResponseDoc[FeatureSet]) (*[]abstract.FeatureSet, error) {
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

func convertDocumentToFeatureSet(document elastic.ResponseDoc[FeatureSet]) (*abstract.FeatureSet, error) {
	fs := abstract.FeatureSet{}
	fs.Name = document.Source.Name
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
			var r bool
			json.Unmarshal([]byte(f.Value), &r)
			af.Value = r
		case "int":
			var r int
			json.Unmarshal([]byte(f.Value), &r)
			af.Value = r
		case "float":
			var r float32
			json.Unmarshal([]byte(f.Value), &r)
			af.Value = r
		case "string":
			af.Value = f.Value
		default:
			var r map[string]interface{}
			json.Unmarshal([]byte(f.Value), &r)
			af.Value = r
		}

		result = append(result, af)
	}
	return &result, nil
}

// CloseConnection ... Terminates the connection to ES for the DAO
func (dao *dao) CloseConnection() {
	dao.Connector.CloseConnection()
}
