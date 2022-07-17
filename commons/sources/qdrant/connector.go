package qdrant

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Connector struct {
	abstract.ConfigurableConnector
	conn             *grpc.ClientConn
	collectionClient pb.CollectionsClient
	pointClient      pb.PointsClient
	collectionName   string
	vectorSize       uint64
	segmentNumber    uint64
}

var reqTimeout = time.Second

func NewQdrantConnector() *Connector {
	return &Connector{
		ConfigurableConnector: abstract.ConfigurableConnector{
			RequiredFields: map[string]string{
				"endpoint":      "endpoint",
				"collection":    "collection",
				"vectorSize":    "vector-size",
				"segmentNumber": "segment-number",
			},
			OptionalFields: map[string]string{
				"certFile": "cert-file",
			},
		},
	}
}

// InitConnection ... inits connection
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	var err error
	endpoint := def.Settings[c.RequiredFields["endpoint"]]
	c.collectionName = def.Settings[c.RequiredFields["collection"]]
	if c.vectorSize, err = strconv.ParseUint(def.Settings[c.RequiredFields["vectorSize"]], 10, 64); err != nil {
		log.Fatalf("failed to parse vector size: %s", err.Error())
	}
	if c.segmentNumber, err = strconv.ParseUint(def.Settings[c.RequiredFields["segmentNumber"]], 10, 64); err != nil {
		log.Fatalf("failed to parse segment number: %s", err.Error())
	}

	var creds credentials.TransportCredentials
	if certFile, exist := def.Settings[c.OptionalFields["certFile"]]; exist {
		if creds, err = credentials.NewClientTLSFromFile(certFile, ""); err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
	} else {
		creds = insecure.NewCredentials()
	}

	if c.conn, err = grpc.Dial(endpoint, grpc.WithTransportCredentials(creds)); err != nil {
		log.Fatalf("Could not connect to %s: %s", endpoint, err.Error())
	}
	c.collectionClient = pb.NewCollectionsClient(c.conn)
	c.pointClient = pb.NewPointsClient(c.conn)

	// check if collection exists already
	c.CreateCollection(context.Background())
}

func (c *Connector) CloseConnection() {
	c.conn.Close()
}

func (c *Connector) DeleteCollection(parentCtx context.Context) error {
	ctx, cancel := context.WithTimeout(parentCtx, reqTimeout)
	defer cancel()

	if _, err := c.collectionClient.Delete(ctx, &pb.DeleteCollection{CollectionName: c.collectionName}); err != nil {
		return fmt.Errorf("Could not delete collection: %v", err)
	}
	log.Println("Collection", c.collectionName, "deleted")
	return nil
}

func (c *Connector) ListCollections(parentCtx context.Context) (collections []*pb.CollectionDescription, err error) {

	ctx, cancel := context.WithTimeout(parentCtx, reqTimeout)
	defer cancel()

	var r *pb.ListCollectionsResponse
	if r, err = c.collectionClient.List(ctx, &pb.ListCollectionsRequest{}); err != nil {
		return nil, fmt.Errorf("Could not get collections: %v", err)
	}
	return r.GetCollections(), nil
}

func (c *Connector) CreateCollection(parentCtx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(parentCtx, reqTimeout)
	defer cancel()

	_, err = c.collectionClient.Create(ctx, &pb.CreateCollection{
		CollectionName: c.collectionName,
		VectorSize:     c.vectorSize,
		Distance:       pb.Distance_Dot,
		OptimizersConfig: &pb.OptimizersConfigDiff{
			DefaultSegmentNumber: &c.segmentNumber,
		},
	})
	return
}

func (c *Connector) CreateFieldIndex(parentCtx context.Context, fieldName string, fieldType *pb.FieldType) (err error) {
	ctx, cancel := context.WithTimeout(parentCtx, reqTimeout)
	defer cancel()

	fieldIndex := &pb.CreateFieldIndexCollection{
		CollectionName: c.collectionName,
		FieldName:      fieldName,
		FieldType:      fieldType,
	}

	_, err = c.pointClient.CreateFieldIndex(ctx, fieldIndex)
	return
}

func (c *Connector) UpsertPoints(parentCtx context.Context, waitUpsert bool, points []*pb.PointStruct) (err error) {
	ctx, cancel := context.WithTimeout(parentCtx, reqTimeout)
	defer cancel()

	_, err = c.pointClient.Upsert(ctx, &pb.UpsertPoints{
		CollectionName: c.collectionName,
		Wait:           &waitUpsert,
		Points:         points,
	})
	return
}

func (c *Connector) GetPointsById(parentCtx context.Context, pointIds ...*pb.PointId) ([]*pb.RetrievedPoint, error) {

	ctx, cancel := context.WithTimeout(parentCtx, reqTimeout)
	defer cancel()

	pointsById, err := c.pointClient.Get(ctx, &pb.GetPoints{
		CollectionName: c.collectionName,
		Ids:            pointIds,
	})
	if err != nil {
		return nil, err
	}
	return pointsById.GetResult(), nil
}

func (c *Connector) DeletePointsByIds(parentCtx context.Context, pointIds ...*pb.PointId) (err error) {
	ctx, cancel := context.WithTimeout(parentCtx, reqTimeout)
	defer cancel()

	_, err = c.pointClient.Delete(ctx, &pb.DeletePoints{
		CollectionName: c.collectionName,
		Points: &pb.PointsSelector{
			PointsSelectorOneOf: &pb.PointsSelector_Points{
				Points: &pb.PointsIdsList{
					Ids: pointIds,
				},
			},
		},
	})

	return
}

func (c *Connector) DeletePointsByName(parentCtx context.Context, name string) (err error) {
	ctx, cancel := context.WithTimeout(parentCtx, reqTimeout)
	defer cancel()

	_, err = c.pointClient.Delete(ctx, &pb.DeletePoints{
		CollectionName: c.collectionName,
		Points: &pb.PointsSelector{
			PointsSelectorOneOf: &pb.PointsSelector_Filter{
				Filter: &pb.Filter{
					Must: []*pb.Condition{
						{
							ConditionOneOf: &pb.Condition_Field{
								Field: &pb.FieldCondition{
									Key: "name",
									Match: &pb.Match{
										MatchValue: &pb.Match_Keyword{
											Keyword: name,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	return
}

func (c *Connector) GetPointsHavingName(parentCtx context.Context, name string) ([]*pb.RetrievedPoint, error) {

	ctx, cancel := context.WithTimeout(parentCtx, reqTimeout)
	defer cancel()

	scrollResponse, err := c.pointClient.Scroll(ctx, &pb.ScrollPoints{
		CollectionName: c.collectionName,
		Filter: &pb.Filter{
			Must: []*pb.Condition{
				{
					ConditionOneOf: &pb.Condition_Field{
						Field: &pb.FieldCondition{
							Key: "name",
							Match: &pb.Match{
								MatchValue: &pb.Match_Keyword{
									Keyword: name,
								},
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}
	return scrollResponse.GetResult(), nil
}

func (c *Connector) SimilarToThis(parentCtx context.Context, point []float32, k uint64, filter *pb.Filter) ([]*pb.ScoredPoint, error) {
	ctx, cancel := context.WithTimeout(parentCtx, reqTimeout)
	defer cancel()

	searchResponse, err := c.pointClient.Search(ctx, &pb.SearchPoints{
		CollectionName: c.collectionName,
		Vector:         point,
		Limit:          k,
		Filter:         filter,
	})
	if err != nil {
		return nil, err
	}
	return searchResponse.GetResult(), nil
}
