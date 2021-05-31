package local

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/datamillcloud/mastro/commons/abstract"
	"github.com/datamillcloud/mastro/commons/utils/conf"
	"github.com/datamillcloud/mastro/commons/utils/strings"
)

type localCrawler struct{}

// NewCrawler ... returns an instance of the crawler
func NewCrawler() abstract.Crawler {
	return &localCrawler{}
}

func (crawler *localCrawler) InitConnection(cfg *conf.Config) (abstract.Crawler, error) {
	// init connection by checking whether the target path is available
	if _, err := os.Stat(cfg.DataSourceDefinition.CrawlerDefinition.Root); os.IsNotExist(err) {
		return nil, err
	}
	return crawler, nil
}

func (crawler *localCrawler) WalkWithFilter(root string, filter string) ([]abstract.Asset, error) {
	var assets []abstract.Asset

	// walk file system
	var walkFn filepath.WalkFunc = func(currentPath string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		// check if it is a regular file (not dir) and the name is like the filter
		if info.Mode().IsRegular() && strings.MatchPattern(info.Name(), filter) {
			// open the file
			stringFile, err := ioutil.ReadFile(currentPath)
			if err != nil {
				return e
			}
			a, err := abstract.ParseAsset(stringFile)
			if err != nil {
				return err
			}
			assets = append(assets, *a)
		}
		return nil
	}
	// walk file system
	// https://golang.org/pkg/path/filepath/#Walk
	// https://flaviocopes.com/go-list-files/
	err := filepath.Walk(root, walkFn)
	return assets, err
}
