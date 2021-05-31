package main

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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/cheggaaa/pb"

	"github.com/datamillcloud/mastro/commons/abstract"
	"github.com/kelseyhightower/envconfig"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gopkg.in/yaml.v2"
)

type FSConnDetails struct {
	Endpoint  string `required:"true"`
	AccessKey string `required:"true"`
	SecretKey string `required:"true"`
	UseSSL    bool   `required:"true"`
	Location  string `required:"true"`
}

type Client struct {
	fs *FSConnDetails
	mc *minio.Client
}

func connect(fs *FSConnDetails) (*Client, error) {
	// Initialize minio client object
	minioClient, err := minio.New(fs.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(fs.AccessKey, fs.SecretKey, ""),
		Secure: fs.UseSSL,
	})

	return &Client{fs, minioClient}, err
}

const DefaultManifestFilename string = "MANIFEST.yaml"

var manifestFilename = DefaultManifestFilename

func (c *Client) ensureBucket(bucketName string) {
	// create bucket
	err := c.mc.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: c.fs.Location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := c.mc.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created bucket %s\n", bucketName)
	}
}

func (c *Client) manifestExists(bucketName string, manifestFilename string) (bool, *minio.ObjectInfo) {
	// check if a metadata file already exists
	if objInfo, err := c.mc.StatObject(context.Background(), bucketName, manifestFilename, minio.StatObjectOptions{}); err == nil {
		return true, &objInfo
	}
	return false, nil
}

func (c *Client) getRemoteManifest(bucketName string, manifestFilename string) (*abstract.Asset, error) {
	reader, err := c.mc.GetObject(context.Background(), bucketName, manifestFilename, minio.GetObjectOptions{})
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

func (c *Client) putLocalManifest(localPath string, destinationPath string, manifestFilename string) (*minio.UploadInfo, error) {
	fi, _ := os.Stat(localPath)
	fileReader, _ := os.Open(localPath)
	// write asset definition to bucket
	progress := pb.New64(fi.Size()).SetUnits(pb.U_BYTES)
	progress.Start()
	uploadInfo, err := c.mc.PutObject(context.Background(), destinationPath, manifestFilename, fileReader, fi.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", Progress: progress})
	return &uploadInfo, err
}

func (c *Client) initDataset(cmd *InitCmd, manifestFilename string) {
	bucketName := cmd.DestinationPath
	localManifestPath := cmd.LocalManifestPath

	// make sure the bucket for the dataset exists
	c.ensureBucket(bucketName)

	if exists, objInfo := c.manifestExists(bucketName, manifestFilename); exists {
		log.Fatalf("Error! %s for %s already exists :: Size:%d Bytes, LastModified:%v", manifestFilename, bucketName, objInfo.Size, objInfo.LastModified)
	}

	// otherwise init manifest file

	// if a manifest is provided load it to the bucket
	if localManifestPath != nil {
		// parse and validate local manifest
		_, err := loadLocalManifest(*localManifestPath)
		if err != nil {
			log.Fatalf("Error - %v", err)
		}
		// write asset definition to bucket
		_, err = c.putLocalManifest(*localManifestPath, bucketName, manifestFilename)
		if err != nil {
			log.Fatalf("Error - %v", err)
		}
	} else {
		// else initialize an empty asset at the current location, so that it can be filled in
		err := initLocalManifest(manifestFilename)
		if err != nil {
			log.Fatalf("Error - %v", err)
		}
	}
}

func (c *Client) overwriteManifest(bucketName string, manifestFilename string, content string) (*minio.UploadInfo, error) {
	reader := strings.NewReader(content)
	progress := pb.New64(reader.Size()).SetUnits(pb.U_BYTES)
	progress.Start()
	uploadInfo, err := c.mc.PutObject(context.Background(), bucketName, manifestFilename, reader, reader.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", Progress: progress})
	if err == nil {
		log.Println("Uploaded", manifestFilename, "of size:", reader.Size(), "Bytes", "Successfully.")
	}
	return &uploadInfo, err
}

func (c *Client) deleteVersionFiles(bucketName string, version string) {
	objectCh := c.mc.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    version,
		Recursive: true,
	})

	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	for object := range objectCh {
		if object.Err != nil {
			panic(object.Err)
		}
		err := c.mc.RemoveObject(context.Background(), bucketName, object.Key, opts)
		if err != nil {
			fmt.Println("Error detected during deletion: ", err)
		}
	}
}

func getVersionedPath(version string, filePath string) (*string, error) {
	base, err := url.Parse(version)
	if err != nil {
		return nil, err
	}
	base.Path = path.Join(base.Path, filePath)
	path := base.String()
	return &path, nil
}

func (c *Client) putFiles(localFolder string, bucketName string, version string) {

	// walk file system
	var walkFn filepath.WalkFunc = func(currentPath string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		// check if it is a regular file
		if info.Mode().IsRegular() {
			// open the file at currentPath
			// open the file
			file, _ := os.Open(currentPath)
			fi, _ := file.Stat()
			defer file.Close()

			progress := pb.New64(fi.Size()).SetUnits(pb.U_BYTES)
			progress.Start()

			p, err := getVersionedPath(version, file.Name())
			if err != nil {
				panic(err)
			}
			_, err = c.mc.PutObject(context.Background(), bucketName, *p, file, fi.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", Progress: progress})
			if err != nil {
				panic(err)
			}

		}
		return nil
	}
	// walk file system
	// https://golang.org/pkg/path/filepath/#Walk
	// https://flaviocopes.com/go-list-files/
	err := filepath.Walk(localFolder, walkFn)
	if err != nil {
		panic(err)
	}
}

func (c *Client) newVersion(cmd *NewCmd, manifestFilename string) {
	asset, err := c.getRemoteManifest(cmd.DestinationPath, manifestFilename)
	if err != nil {
		panic(err)
	}
	version := strconv.FormatInt(getVersionAsUnixTimeInSeconds(time.Now()), 10)
	versionMetadata := map[string]string{}
	asset.Versions[version] = versionMetadata

	data, err := yaml.Marshal(&asset)
	if err != nil {
		panic(err)
	}
	c.overwriteManifest(cmd.DestinationPath, manifestFilename, string(data))
	fmt.Println(version)
}

func (c *Client) add(cmd *AddCmd, manifestFilename string) {
	// open manifest and get newest version
	versions, err := c.getVersions(cmd.DestinationPath, manifestFilename)
	if err != nil {
		panic(err)
	}

	if len(versions) > 0 {
		latestVersion := versions[0]
		c.putFiles(cmd.LocalPath, cmd.DestinationPath, latestVersion)
	} else {
		panic(fmt.Sprintf("No versions found at %s \n", cmd.DestinationPath))
	}

}

func (c *Client) getVersions(destinationPath string, manifestFilename string) ([]string, error) {
	// open manifest and get newest version
	asset, err := c.getRemoteManifest(destinationPath, manifestFilename)
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

func (c *Client) allVersions(cmd *VersionsCmd, manifestFilename string) {
	// open manifest and get newest version
	versions, err := c.getVersions(cmd.DestinationPath, manifestFilename)
	if err != nil {
		panic(err)
	}

	if len(versions) > 0 {
		fmt.Println(versions)
	} else {
		fmt.Println(fmt.Sprintf("No versions found at %s", cmd.DestinationPath))
	}
}

func (c *Client) latestVersion(cmd *LatestCmd, manifestFilename string) {
	// open manifest and get newest version
	versions, err := c.getVersions(cmd.DestinationPath, manifestFilename)
	if err != nil {
		panic(err)
	}

	if len(versions) > 0 {
		fmt.Println(versions[0])
	} else {
		fmt.Printf("No versions found at %s \n", cmd.DestinationPath)
	}
}

func (c *Client) deleteVersion(cmd *DeleteCmd) {
	c.deleteVersionFiles(cmd.DestinationPath, cmd.Version)
	// todo: update manifest
}

func (c *Client) overwriteVersion(cmd *OverwriteCmd) {
	c.deleteVersionFiles(cmd.DestinationPath, cmd.Version)
	c.putFiles(cmd.LocalPath, cmd.DestinationPath, cmd.Version)
	// todo: update manifest file
}

func main() {
	/*
		fmt.Println(`
		╔╦╗╦  ╦╔═╗
		║║║╚╗╔╝║
		╩ ╩ ╚╝ ╚═╝`)
	*/

	// make sure connection details to bucket are set as env vars
	var fs FSConnDetails
	err := envconfig.Process("mvc", &fs)
	if err != nil {
		panic(err)
	}
	// connect to remote FS
	client, err := connect(&fs)
	if err != nil {
		panic(err)
	}

	// parse command arguments
	arg.MustParse(&args)

	// call specific subcommand handler
	switch {
	case args.Init != nil:
		client.initDataset(args.Init, manifestFilename)
	case args.New != nil:
		client.newVersion(args.New, manifestFilename)
	case args.Add != nil:
		client.add(args.Add, manifestFilename)
	case args.Versions != nil:
		client.allVersions(args.Versions, manifestFilename)
	case args.Latest != nil:
		client.latestVersion(args.Latest, manifestFilename)
	case args.Overwrite != nil:
		client.overwriteVersion(args.Overwrite)
	default:
		panic(fmt.Sprintf("unknown command %q", os.Args[0]))
	}
}
