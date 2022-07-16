package git

import (
	"log"
	"os"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
)

// NewGitConnector factory
func NewGitConnector() *Connector {
	return &Connector{
		ConfigurableConnector: abstract.ConfigurableConnector{
			RequiredFields: map[string]string{
				"gitPwdOrToken": "pwd-token-var",
				"repo":          "repo",
			},
			OptionalFields: map[string]string{
				"storageType": "storage-type",
				"gitUser":     "username-var",
				"pemFile":     "pem-file",
			},
		},
	}
}

// Connector ... Connector type
type Connector struct {
	abstract.ConfigurableConnector
	storage storage.Storer
	Fs      billy.Filesystem
	Repo    *git.Repository
}

// InitConnection ... Starts a connection with Git
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	var err error

	// get storage to where to clone all repositories
	if storageType, exist := def.Settings[c.OptionalFields["storageType"]]; exist {
		if storageType == "memory" {
			c.storage = memory.NewStorage()
		} else {
			// assume a well-formed path
			// todo: check if the path is valid
			// the repo will have a worktree
			// https://github.com/go-git/go-git/blob/db2bc57350561c4368a8d32c42476699b48d2a09/repository.go#L216
			wt := osfs.New(storageType)
			c.Fs, err = wt.Chroot(git.GitDirName)
			if err != nil {
				log.Fatal(err)
			}
			c.storage = filesystem.NewStorage(c.Fs, cache.NewObjectLRUDefault())
			//c.storage = filesystem.NewStorageWithOptions(c.fs, cache.NewObjectLRUDefault(), filesystem.Options{KeepDescriptors: true})
		}
	} else {
		// default to memory storage
		c.storage = memory.NewStorage()
		c.Fs = memfs.New()
	}

	// retrieve repo auth info
	repoUrl := def.Settings[c.RequiredFields["repo"]]
	password := os.Getenv(def.Settings[c.RequiredFields["gitPwdOrToken"]])

	// basic conn opts
	options := &git.CloneOptions{URL: repoUrl}

	// https://github.com/go-git/go-git/blob/master/_examples/clone/auth/basic/access_token/main.go
	// https://github.com/go-git/go-git/blob/master/_examples/clone/auth/basic/username_password/main.go
	// https://github.com/go-git/go-git/blob/master/_examples/clone/auth/ssh/main.go
	// if a pem and a password are provided then use them
	if pemFile, exist := def.Settings[c.OptionalFields["pemFile"]]; exist {
		publicKeys, err := ssh.NewPublicKeysFromFile("git", pemFile, password)
		if err != nil {
			log.Fatal(err)
		}
		options.Auth = publicKeys
	} else {
		// default to username and password
		if userEnv, exist := def.Settings[c.OptionalFields["gitUser"]]; exist {
			options.Auth = &http.BasicAuth{
				Username: os.Getenv(userEnv),
				Password: password,
			}
		} else {
			log.Fatalf("Unset field %s", c.OptionalFields["gitUser"])
		}
	}

	// clone the repo
	c.Repo, err = git.Clone(c.storage, c.Fs, options)
	if err != nil {
		log.Fatal(err)
	}
}

// CloseConnection ... Closes the connection with Git
func (c *Connector) CloseConnection() {
	// noop
}
