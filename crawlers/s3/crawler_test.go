package s3

import (
	"context"
	"os"
	"testing"

	"github.com/datamillcloud/mastro/commons/abstract"
	"github.com/datamillcloud/mastro/commons/utils/conf"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
)

var (
	cfg     *conf.Config
	crawler *s3Crawler
	connErr error
)

func TestMain(m *testing.M) {
	cfg = conf.Load("example_s3.cfg")
	crawler = NewCrawler().(*s3Crawler)
	_, connErr = crawler.InitConnection(cfg)

	// run module tests
	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestConnection(t *testing.T) {
	t.Log("TestConnection running")
	assert := assert.New(t)
	assert.Equal(connErr, nil)
	assert.NotEqual(crawler.GetClient(), nil)
}

func TestWalk(t *testing.T) {
	t.Log("TestWalk running")
	assert := assert.New(t)

	bucketName := cfg.DataSourceDefinition.Settings["bucket"]
	// create bucket if not existing
	crawler.GetClient().MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})

	// Name of the object
	objectName := "exampleObject"
	// Path to file to be uploaded
	filePath := "file.csv"

	f, err := os.Create(filePath)
	_, _ = f.WriteString("hello world\n")
	f.Close()

	size, err := crawler.GetClient().FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: "application/csv"})

	assert.Equal(err, nil)
	assert.NotEqual(size, 0)

	// test walk function
	fs, err := crawler.WalkWithFilter(crawler.connector.Prefix, abstract.DefaultManifestFilename)

	assert.Equal(err, nil)
	assert.NotEqual(fs, nil)

	// remove remote object and bucket and ciao
	//err = crawler.GetClient().RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	assert.Equal(err, nil)
	err = crawler.GetClient().RemoveBucket(context.Background(), bucketName)
	assert.Equal(err, nil)

	err = os.Remove(filePath)
	assert.Equal(err, nil)
}
