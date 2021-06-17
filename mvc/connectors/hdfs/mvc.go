package hdfs

import (
	"bytes"
	"fmt"
	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/sources/hdfs"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/mvc/commons"
	"github.com/sger/go-hashdir"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
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
)

func NewMvc(manifestFilename string) commons.MvcProvider {
	return &HDFSMvc{manifestFilename: manifestFilename, connector: hdfs.NewHDFSConnector()}
}

type HDFSMvc struct {
	manifestFilename string
	connector        *hdfs.Connector
}

func (mvc *HDFSMvc) SetManifestFilename(manifestFilename string) {
	mvc.manifestFilename = manifestFilename
}

func (mvc *HDFSMvc) SetConnector(connector *hdfs.Connector) {
	mvc.connector = connector
}

func (mvc *HDFSMvc) GetConnector() *hdfs.Connector {
	return mvc.connector
}

func (mvc *HDFSMvc) InitConnection(cfg *conf.Config) (commons.MvcProvider, error) {
	//mvc.connector = s3.NewS3Connector()
	if err := mvc.connector.ValidateDataSourceDefinition(&cfg.DataSourceDefinition); err != nil {
		log.Panicln(err)
	}
	// inits connection
	mvc.connector.InitConnection(&cfg.DataSourceDefinition)
	return mvc, nil
}

func (mvc *HDFSMvc) ensurePath(pathName string) {

	err := mvc.connector.GetClient().Mkdir(pathName, os.ModeDir)

	if err != nil {
		log.Fatalf("HDFS path %s already exists", pathName)
	}

}

func (mvc *HDFSMvc) manifestExists(pathName string) (bool, os.FileInfo) {
	// check if a metadata file already exists
	if objInfo, err := mvc.connector.GetClient().Stat(fmt.Sprintf("%s/%s", pathName, mvc.manifestFilename)); err == nil {
		return true, objInfo
	}
	return false, nil
}

func (mvc *HDFSMvc) InitDataset(cmd *commons.InitCmd) {
	hdfsPathName := cmd.DestinationPath
	localManifestPath := cmd.LocalManifestPath

	// make sure the bucket for the dataset exists
	mvc.ensurePath(hdfsPathName)

	if exists, fileInfo := mvc.manifestExists(hdfsPathName); exists {
		log.Fatalf("Error! %s for %s already exists :: Size:%d Bytes, LastModified:%v", mvc.manifestFilename, hdfsPathName, fileInfo.Size(), fileInfo.ModTime())
	}

	// otherwise init manifest file

	// if a manifest is provided load it to the bucket
	if localManifestPath != nil {
		// parse and validate local manifest
		_, err := commons.LoadLocalManifest(*localManifestPath)
		if err != nil {
			log.Fatalf("Error - %v", err)
		}
		// write asset definition to hdfs path
		_, err = mvc.PutLocalManifest(*localManifestPath, hdfsPathName)
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

func (mvc *HDFSMvc) getRemoteManifest(hdfsPathName string) (*abstract.Asset, error) {
	reader, err := mvc.connector.GetClient().Open(fmt.Sprintf("%s/%s", hdfsPathName, mvc.manifestFilename))
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()
	stat := reader.Stat()
	buf := new(bytes.Buffer)
	if _, err := io.CopyN(buf, reader, stat.Size()); err != nil {
		return nil, err
	}

	a, err := abstract.ParseAsset(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (mvc *HDFSMvc) PutLocalManifest(localPath string, destinationPath string) (*os.FileInfo, error) {
	fileBytes, _ := ioutil.ReadFile(localPath)

	fileWriter, err := mvc.connector.GetClient().Create(fmt.Sprintf("%s/%s", destinationPath, mvc.manifestFilename))

	if err != nil {
		log.Fatalf("Error - %v", err)
	}

	_, err = fileWriter.Write(fileBytes)

	if err != nil {
		log.Fatalf("Error - %v", err)
	}

	defer fileWriter.Close()

	// stats created file from hdfs path
	uploadInfo, err := mvc.connector.GetClient().Stat(fmt.Sprintf("%s/%s", destinationPath, mvc.manifestFilename))
	return &uploadInfo, err
}

func (mvc *HDFSMvc) OverwriteManifest(destinationPath string, content string) (*os.FileInfo, error) {
	reader := strings.NewReader(content)
	fileWriter, err := mvc.connector.GetClient().Create(fmt.Sprintf("%s/%s", destinationPath, mvc.manifestFilename))
	if err != nil {
		log.Fatalf("Error - %v", err)
	}

	_, err = reader.WriteTo(fileWriter)

	if err != nil {
		log.Fatalf("Error - %v", err)
	}

	defer fileWriter.Close()

	// stats created file from hdfs path
	uploadInfo, err := mvc.connector.GetClient().Stat(fmt.Sprintf("%s/%s", destinationPath, mvc.manifestFilename))
	return &uploadInfo, err
}

func (mvc *HDFSMvc) PutFiles(localFolder string, hdfsPathName string, version string) {

	// walk file system
	var walkFn filepath.WalkFunc = func(currentPath string, info os.FileInfo, e error) error {
		// get a ref to the parent's path
		ref := filepath.Dir(filepath.Clean(localFolder))
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

			p, err := GetVersionedPath(version, rel)
			if err != nil {
				panic(err)
			}
			fmt.Println(*p)

			fileBytes, _ := ioutil.ReadFile(currentPath)

			fileWriter, err := mvc.connector.GetClient().Create(fmt.Sprintf("%s/%s", hdfsPathName, *p))

			if err != nil {
				log.Fatalf("Error - %v", err)
			}

			_, err = fileWriter.Write(fileBytes)

			if err != nil {
				log.Fatalf("Error - %v", err)
			}

			defer fileWriter.Close()

			_, err = mvc.connector.GetClient().Stat(fmt.Sprintf("%s/%s", hdfsPathName, *p))
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

func (mvc *HDFSMvc) NewVersion(cmd *commons.NewCmd) {
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

func (mvc *HDFSMvc) Add(cmd *commons.AddCmd) {
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

func (mvc *HDFSMvc) GetVersions(destinationPath string) ([]string, error) {
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

func (mvc *HDFSMvc) AllVersions(cmd *commons.VersionsCmd) {
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

func (mvc *HDFSMvc) LatestVersion(cmd *commons.LatestCmd) {
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

func (mvc *HDFSMvc) DeleteVersionFiles(hdfsPathName string, version string) {
	versionedPath, err := GetVersionedPath(version, hdfsPathName)

	if err != nil {
		fmt.Println(fmt.Sprintf("Error - %v", err))
		return
	}

	err = mvc.connector.GetClient().RemoveAll(*versionedPath)
	if err != nil {
		fmt.Println("Error detected during deletion: ", err)
	}
}

func (mvc *HDFSMvc) editVersionMetadata(localPath string, destinationPath string, version string, append bool) {
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

func (mvc *HDFSMvc) DeleteVersionMetadata(cmd *commons.DeleteCmd) {
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

func (mvc *HDFSMvc) DeleteVersion(cmd *commons.DeleteCmd) {
	mvc.DeleteVersionFiles(cmd.DestinationPath, cmd.Version)
	mvc.DeleteVersionMetadata(cmd)
}

func (mvc *HDFSMvc) OverwriteVersion(cmd *commons.OverwriteCmd) {
	mvc.DeleteVersionFiles(cmd.DestinationPath, cmd.Version)
	mvc.PutFiles(cmd.LocalPath, cmd.DestinationPath, cmd.Version)
	mvc.editVersionMetadata(cmd.LocalPath, cmd.DestinationPath, cmd.Version, false)
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
