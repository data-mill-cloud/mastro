package paths

import (
	"os"
	"path/filepath"
)

func GetRootDir() (*string, error) {
	e, err := os.Executable()
	if err != nil {
		return nil, err
	}
	rootPath := filepath.Dir(e)
	return &rootPath, nil
}

func GetProjPath(e ...string) (*string, error) {
	rootPath, err := GetRootDir()
	if err != nil {
		return nil, err
	}
	// prepend project path to user relative paths
	e = append([]string{*rootPath}, e...)
	projPath := filepath.Join(e...)
	return &projPath, nil
}
