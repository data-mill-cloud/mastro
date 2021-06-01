package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/sources/s3"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/commons/utils/strings"
	"github.com/minio/minio-go/v7"
)

type s3Crawler struct {
	connector *s3.Connector
}

// NewCrawler ... returns an instance of the crawler
func NewCrawler() abstract.Crawler {
	return &s3Crawler{}
}

// GetClient ... get client for test purposes
func (crawler *s3Crawler) GetClient() *minio.Client {
	return crawler.connector.GetClient()
}

func (crawler *s3Crawler) InitConnection(cfg *conf.Config) (abstract.Crawler, error) {
	crawler.connector = s3.NewS3Connector()
	if err := crawler.connector.ValidateDataSourceDefinition(&cfg.DataSourceDefinition); err != nil {
		log.Panicln(err)
	}
	// inits connection
	crawler.connector.InitConnection(&cfg.DataSourceDefinition)

	// set filter for the manifest filename
	//crawler.config = &cfg.CrawlerDefinition

	return crawler, nil
}

/*
func (crawler *s3Crawler) Walk(bucket string) ([]minio.ObjectInfo, error) {

	exists, errBucketExists := crawler.GetClient().BucketExists(context.Background(), bucket)
	if errBucketExists != nil {
		return nil, errBucketExists
	}

	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", bucket)
	}

	return crawler.ListObjects(bucket, "", true, abstract.DefaultManifestFilename)
}

func (crawler *s3Crawler) ListBuckets() ([]minio.BucketInfo, error) {
	return crawler.GetClient().ListBuckets(context.Background())
}
*/

func (crawler *s3Crawler) ListObjects(bucket string, prefix string, recursive bool, filter string) ([]minio.ObjectInfo, error) {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	ctx := context.Background()

	opts := minio.ListObjectsOptions{
		Recursive: recursive,
		Prefix:    prefix,
	}

	objectCh := crawler.GetClient().ListObjects(ctx, bucket, opts)
	var slice []minio.ObjectInfo

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		// append iff the file matches the given pattern
		if strings.MatchPattern(object.Key, filter) {
			slice = append(slice, object)
		}
	}

	return slice, nil
}

func (crawler *s3Crawler) WalkWithFilter(root string, filter string) ([]abstract.Asset, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//ctx := context.Background()

	exists, errBucketExists := crawler.GetClient().BucketExists(ctx, crawler.connector.Bucket)
	if errBucketExists != nil {
		return nil, errBucketExists
	}

	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", crawler.connector.Bucket)
	}

	objs, err := crawler.ListObjects(root, crawler.connector.Prefix, true, filter) //crawler.config.FilterFilename)
	if err != nil {
		return nil, err
	}

	var assets []abstract.Asset
	opts := minio.GetObjectOptions{}
	for _, o := range objs {
		log.Println("Found ", o.Key)
		reader, err := crawler.connector.GetClient().GetObject(ctx, crawler.connector.Bucket, o.Key, opts)
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		stat, err := reader.Stat()
		if err != nil {
			return nil, err
		}
		buf := new(bytes.Buffer)
		if _, err := io.CopyN(buf, reader, stat.Size); err != nil {
			return nil, err
		}

		a, err := abstract.ParseAsset(buf.Bytes())
		if err != nil {
			return nil, err
		}
		assets = append(assets, *a)
	}

	return assets, nil
}
