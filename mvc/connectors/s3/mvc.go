package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/sources/s3"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/mvc/commons"
	"github.com/minio/minio-go/v7"
	"github.com/sger/go-hashdir"
	"gopkg.in/yaml.v2"
)

func NewMvc(manifestFilename string) commons.MvcProvider {
	return &S3Mvc{manifestFilename: manifestFilename, connector: s3.NewS3Connector()}
}

type S3Mvc struct {
	manifestFilename string
	connector        *s3.Connector
}

func (mvc *S3Mvc) SetManifestFilename(manifestFilename string) {
	mvc.manifestFilename = manifestFilename
}

func (mvc *S3Mvc) SetConnector(connector *s3.Connector) {
	mvc.connector = connector
}

func (mvc *S3Mvc) GetConnector() *s3.Connector {
	return mvc.connector
}

func (mvc *S3Mvc) InitConnection(cfg *conf.Config) (commons.MvcProvider, error) {
	//mvc.connector = s3.NewS3Connector()
	if err := mvc.connector.ValidateDataSourceDefinition(&cfg.DataSourceDefinition); err != nil {
		log.Panicln(err)
	}
	// inits connection
	mvc.connector.InitConnection(&cfg.DataSourceDefinition)
	return mvc, nil
}

func (mvc *S3Mvc) ensureBucket(bucketName string) {

	bucketOptions := minio.MakeBucketOptions{}
	if len(mvc.connector.Region) > 0 {
		bucketOptions.Region = mvc.connector.Region
	}

	// create bucket
	err := mvc.connector.GetClient().MakeBucket(context.Background(), bucketName, bucketOptions)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := mvc.connector.GetClient().BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created bucket %s\n", bucketName)
	}
}

func (mvc *S3Mvc) manifestExists(bucketName string) (bool, *minio.ObjectInfo) {
	// check if a metadata file already exists
	if objInfo, err := mvc.connector.GetClient().StatObject(context.Background(), bucketName, mvc.manifestFilename, minio.StatObjectOptions{}); err == nil {
		return true, &objInfo
	}
	return false, nil
}

func (mvc *S3Mvc) getRemoteManifest(bucketName string) (*abstract.Asset, error) {
	reader, err := mvc.connector.GetClient().GetObject(context.Background(), bucketName, mvc.manifestFilename, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
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
	return a, nil
}

func (mvc *S3Mvc) PutLocalManifest(localPath string, destinationPath string) (*minio.UploadInfo, error) {
	fi, _ := os.Stat(localPath)
	fileReader, _ := os.Open(localPath)
	// write asset definition to bucket
	progress := pb.New64(fi.Size()).SetUnits(pb.U_BYTES)
	progress.Start()
	uploadInfo, err := mvc.connector.GetClient().PutObject(context.Background(), destinationPath, mvc.manifestFilename, fileReader, fi.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", Progress: progress})
	return &uploadInfo, err
}

func (mvc *S3Mvc) InitDataset(cmd *commons.InitCmd) {
	bucketName := cmd.DestinationPath
	localManifestPath := cmd.LocalManifestPath

	// make sure the bucket for the dataset exists
	mvc.ensureBucket(bucketName)

	if exists, objInfo := mvc.manifestExists(bucketName); exists {
		log.Fatalf("Error! %s for %s already exists :: Size:%d Bytes, LastModified:%v", mvc.manifestFilename, bucketName, objInfo.Size, objInfo.LastModified)
	}

	// otherwise init manifest file

	// if a manifest is provided load it to the bucket
	if localManifestPath != nil {
		// parse and validate local manifest
		_, err := commons.LoadLocalManifest(*localManifestPath)
		if err != nil {
			log.Fatalf("Error - %v", err)
		}
		// write asset definition to bucket
		_, err = mvc.PutLocalManifest(*localManifestPath, bucketName)
		if err != nil {
			log.Fatalf("Error - %v", err)
		}
	} else {
		// else initialize an empty asset at the current location, so that it can be filled in
		err := commons.InitLocalManifest(mvc.manifestFilename)
		if err != nil {
			log.Fatalf("Error - %v", err)
		}
	}
}

func (mvc *S3Mvc) OverwriteManifest(bucketName string, content string) (*minio.UploadInfo, error) {
	reader := strings.NewReader(content)
	progress := pb.New64(reader.Size()).SetUnits(pb.U_BYTES)
	progress.Start()
	uploadInfo, err := mvc.connector.GetClient().PutObject(context.Background(), bucketName, mvc.manifestFilename, reader, reader.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", Progress: progress})
	return &uploadInfo, err
}

func (mvc *S3Mvc) DeleteVersionFiles(bucketName string, version string) {
	/*
		// iterate and delete all objects
		objectCh := mvc.connector.GetClient().ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
			Prefix:    version,
			Recursive: true,
		})

		opts := minio.RemoveObjectOptions{
			GovernanceBypass: true,
		}
		for object := range objectCh {
			if object.Err != nil {
				fmt.Println(object.Err)
				return
			}
			err := mvc.connector.GetClient().RemoveObject(context.Background(), bucketName, object.Key, opts)
			if err != nil {
				fmt.Println("Error detected during deletion: ", err)
			}
		}
	*/

	// use channels to list and delete object lists
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		listOpts := minio.ListObjectsOptions{
			Prefix:    version,
			Recursive: true,
		}
		for object := range mvc.connector.GetClient().ListObjects(context.Background(), bucketName, listOpts) {
			if object.Err != nil {
				fmt.Println(object.Err)
				return
			}
			objectsCh <- object
		}
	}()

	for err := range mvc.connector.GetClient().RemoveObjects(context.Background(), bucketName, objectsCh, minio.RemoveObjectsOptions{GovernanceBypass: true}) {
		fmt.Println("Error detected during deletion:", err)
	}
}

func GetVersionedPath(version string, filePath string) (*string, error) {
	base, err := url.Parse(version)
	if err != nil {
		return nil, err
	}
	base.Path = path.Join(base.Path, filePath)
	path := base.String()
	return &path, nil
}

func (mvc *S3Mvc) PutFiles(localPath string, bucketName string, version string) {
	// get a ref to the parent's path
	ref := filepath.Dir(filepath.Clean(localPath))

	// walk file system
	var walkFn filepath.WalkFunc = func(currentPath string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		// check if it is a regular file
		if info.Mode().IsRegular() {

			// compute relative path to avoid saving to s3 the absolute path
			rel, err := filepath.Rel(ref, currentPath)
			if err != nil {
				log.Println(err.Error())
				return e
			}
			// visiting the same file being passed as input
			if rel == "." {
				// getting filename only (basename)
				rel = filepath.Base(currentPath)
			}

			// open the file at currentPath
			file, _ := os.Open(currentPath)
			fi, _ := file.Stat()

			defer file.Close()

			progress := pb.New64(fi.Size()).SetUnits(pb.U_BYTES)
			progress.Start()

			p, err := GetVersionedPath(version, rel)
			if err != nil {
				panic(err)
			}
			fmt.Println(*p)

			_, err = mvc.connector.GetClient().PutObject(context.Background(), bucketName, *p, file, fi.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", Progress: progress})
			if err != nil {
				panic(err)
			}
		}
		return nil
	}
	// walk file system
	// https://golang.org/pkg/path/filepath/#Walk
	// https://flaviocopes.com/go-list-files/
	err := filepath.Walk(localPath, walkFn)
	if err != nil {
		panic(err)
	}
}

func (mvc *S3Mvc) NewVersion(cmd *commons.NewCmd) {
	asset, err := mvc.getRemoteManifest(cmd.DestinationPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while retrieving remote manifest from path %s :: %s", cmd.DestinationPath, err))
		return
	}
	version := strconv.FormatInt(commons.GetVersionAsUnixTimeInSeconds(time.Now()), 10)
	versionMetadata := map[string]string{}
	asset.Versions[version] = versionMetadata

	data, err := yaml.Marshal(&asset)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while parsing yaml manifest from path %s :: %s", cmd.DestinationPath, err))
		return
	}

	mvc.OverwriteManifest(cmd.DestinationPath, string(data))

	fmt.Println("\n", version)
}

func (mvc *S3Mvc) editVersionMetadata(localPath string, destinationPath string, version string, append bool) {
	asset, err := mvc.getRemoteManifest(destinationPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, ok := asset.Versions[version]; ok {
		var versionMetadata interface{}
		if append {
			versionMetadata = asset.Versions[version]
		} else {
			versionMetadata = map[string]string{}
		}
		v := reflect.ValueOf(versionMetadata)

		// --
		// compute hash of newly added element bunch
		basename := filepath.Base(filepath.Clean(localPath))
		//h, err := dirhash.HashDir(localPath, "a", dirhash.DefaultHash)
		h, err := hashdir.Create(localPath, "sha256")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(basename, h)
		// add hashes for files/folders to version metadata
		if v.Kind() == reflect.Map {
			key := reflect.ValueOf(basename)
			value := reflect.ValueOf(h)
			v.SetMapIndex(key, value)
		}
		// --
		asset.Versions[version] = versionMetadata
		data, err := yaml.Marshal(&asset)
		if err != nil {
			fmt.Println(err)
			return
		}
		mvc.OverwriteManifest(destinationPath, string(data))
	} else {
		fmt.Println(fmt.Sprintf("No version %s found", version))
	}
}

func (mvc *S3Mvc) Add(cmd *commons.AddCmd) {
	// open manifest and get newest version
	versions, err := mvc.GetVersions(cmd.DestinationPath)
	if err != nil {
		panic(err)
	}

	if len(versions) > 0 {
		latestVersion := versions[0]
		mvc.PutFiles(cmd.LocalPath, cmd.DestinationPath, latestVersion)
		mvc.editVersionMetadata(cmd.LocalPath, cmd.DestinationPath, latestVersion, true)
	} else {
		fmt.Println(fmt.Sprintf("No versions found at %s \n", cmd.DestinationPath))
		return
	}

}

func (mvc *S3Mvc) GetVersions(destinationPath string) ([]string, error) {
	// open manifest and get newest version
	asset, err := mvc.getRemoteManifest(destinationPath)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(asset.Versions))
	for k, _ := range asset.Versions {
		keys = append(keys, k)
	}
	//sort.Strings(keys) // ASC
	sort.Sort(sort.Reverse(sort.StringSlice(keys))) // DESC
	return keys, nil
}

func (mvc *S3Mvc) AllVersions(cmd *commons.VersionsCmd) {
	// open manifest and get newest version
	versions, err := mvc.GetVersions(cmd.DestinationPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while retrieving versions from path %s :: %s", cmd.DestinationPath, err))
		return
	}

	if len(versions) > 0 {
		fmt.Println(versions)
	} else {
		fmt.Println(fmt.Sprintf("No versions found at %s", cmd.DestinationPath))
	}
}

func (mvc *S3Mvc) LatestVersion(cmd *commons.LatestCmd) {
	// open manifest and get newest version
	versions, err := mvc.GetVersions(cmd.DestinationPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while retrieving latest version from path %s :: %s", cmd.DestinationPath, err))
		return
	}

	if len(versions) > 0 {
		fmt.Println(versions[0])
	} else {
		fmt.Printf("No versions found at %s \n", cmd.DestinationPath)
	}
}

func (mvc *S3Mvc) DeleteVersionMetadata(cmd *commons.DeleteCmd) {
	asset, err := mvc.getRemoteManifest(cmd.DestinationPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, ok := asset.Versions[cmd.Version]; ok {
		delete(asset.Versions, cmd.Version)
		data, err := yaml.Marshal(&asset)
		if err != nil {
			fmt.Println(err)
			return
		}
		mvc.OverwriteManifest(cmd.DestinationPath, string(data))
	} else {
		fmt.Println(fmt.Sprintf("No version %s found", cmd.Version))
	}
}

func (mvc *S3Mvc) DeleteVersion(cmd *commons.DeleteCmd) {
	mvc.DeleteVersionFiles(cmd.DestinationPath, cmd.Version)
	mvc.DeleteVersionMetadata(cmd)
}

func (mvc *S3Mvc) OverwriteVersion(cmd *commons.OverwriteCmd) {
	mvc.DeleteVersionFiles(cmd.DestinationPath, cmd.Version)
	mvc.PutFiles(cmd.LocalPath, cmd.DestinationPath, cmd.Version)
	mvc.editVersionMetadata(cmd.LocalPath, cmd.DestinationPath, cmd.Version, false)
}
