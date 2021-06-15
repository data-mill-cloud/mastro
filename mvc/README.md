# MVC - Mastro Version Control
A minimal data versioning tool in Golang.

TLDR:
* for each data asset keeps a manifest file that can be crawled and stored in a mastro catalogue.  
* based on the `commons.abstract.sources` package - multiple connectors available 

## Prerequisites

An mvc provider is available for the desired backend storage to be used for file versioning. 
Mind that a provider is defined as follows: 
```go
type MvcProvider interface {
	InitConnection(cfg *conf.Config) (MvcProvider, error)
	InitDataset(cmd *InitCmd)
	NewVersion(cmd *NewCmd)
	Add(cmd *AddCmd)
	AllVersions(cmd *VersionsCmd)
	LatestVersion(cmd *LatestCmd)
	OverwriteVersion(cmd *OverwriteCmd)
	DeleteVersion(cmd *DeleteCmd)
}
```
The mvc provider instantiates a [mastro connector](../doc/CONNECTORS.md) within the `InitConnection` function, as specified in the commons module.

## Requirements

In order for **mvc** to work, a Mastro configuration file of kind `mvc` must be specified and referred to using the `MVC_CONFIG` variable, e.g.:

```bash
./mvc -h
required key MVC_CONFIG missing value
```

```bash
export MVC_CONFIG=$PWD/conf/example_s3.yml
```

where `example_s3.yml` refers to the public minio:

```yaml
type: mvc
backend:
  name: public-minio-s3
  type: s3
  settings:
    region: us-east-1
    endpoint: play.min.io
    access-key-id: Q3AM3UQ867SPQQA43P2F
    secret-access-key: zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG
    use-ssl: true
    bucket: ""
```

Let us now list available versions for the path `abcde`:

```bash
./mvc versions -d abcde
2021/06/10 14:44:58 Successfully loaded config mvc public-minio-s3
2021/06/10 14:44:58 Successfully validated data source definition
2021/06/10 14:44:58 Using provided region us-east-1
[1623324009]
```

## Usage

### Manifest creation
* `mvc init -d $PATH` - initializes local metadata file (i.e. manifest) for an asset located at $PATH
* `mvc init -d $PATH -f $MANIFESTPATH` - uploads manifest file located at $MANIFESTPATH at $PATH

### Version management
* `mvc new -d $PATH` - creates new version and returns full path at $PATH
* `mvc versions -d $PATH` - retrieves all available versions at $PATH and shows their metadata
* `mvc latest -d $PATH` - retrieves latest version at $PATH
* `mvc delete -d $PATH -v $VERSION` - deletes the specified version and updates the metadata

### File management
* `mvc add -l $LOCALPATH -d $PATH` - adds $LOCALPATH to remote $PATH at current latest version, includes the sha256 in the version metadata
* `mvc overwrite -d $PATH -v $VERSION -l $LOCALPATH` - overwrite existing version $VERSION at $PATH and overwrites metadata

### Checksum
* `mvc check -l $LOCALPATH` - computes the sha256sum of the entire folder at $LOCALPATH