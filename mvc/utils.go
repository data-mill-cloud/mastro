package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"gopkg.in/yaml.v2"
)

func getVersionAsUnixTimeInSeconds(t time.Time) int64 {
	return t.Unix()
}

func getVersionAsUnixTimeInNano(t time.Time) int64 {
	return t.UnixNano()
}

func getVersionAsUnixTimeInMillis(t time.Time) int64 {
	return getVersionAsUnixTimeInNano(t) / int64(time.Millisecond)
}

func loadLocalManifest(path string) (*abstract.Asset, error) {
	log.Printf("Loading manifest %s\n", path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	asset, err := abstract.ParseAsset(data)
	if err != nil {
		return nil, err
	}
	if err := asset.Validate(); err != nil {
		return nil, err
	}
	return asset, nil
}

func serializeAsset(asset *abstract.Asset) ([]byte, error) {
	return yaml.Marshal(asset)
}

func initLocalManifest(manifestFilename string) error {
	a := abstract.NewDatasetAsset()
	// write to local file
	d, err := yaml.Marshal(&a)
	if err != nil {
		return err
	}
	log.Printf("Saving to local file %s\n", manifestFilename)
	fmt.Printf("---\n%s\n\n", string(d))
	f, err := os.Create(manifestFilename)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	n, err := w.WriteString(string(d))
	if err != nil {
		return err
	}
	log.Printf("wrote %d bytes\n", n)
	w.Flush()
	return nil
}
