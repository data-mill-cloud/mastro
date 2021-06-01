# MVC - Mastro Version Control
A minimal data versioning tool in Golang.

## Manifest creation
* `mvc init -d $PATH` - initializes local metadata file (i.e. manifest) for asset located at $PATH
* `mvc init -d $PATH -f $MANIFESTPATH` - uploads manifest file located at $MANIFESTPATH at $PATH

## Version management
* `mvc new -d $PATH` - creates new version and returns full path at $PATH
* `mvc versions -d $PATH` - retrieves all available versions at $PATH and shows their metadata
* `mvc latest -d $PATH` - retrieves latest version at $PATH
* `mvc delete -d $PATH -v $VERSION` - deletes the specified version and updates the metadata

## File management
* `mvc add -l $LOCALPATH -d $PATH` - adds $LOCALPATH to remote $PATH at current latest version
* `mvc overwrite -d $PATH -v $VERSION -l $LOCALPATH` - overwrite existing version $VERSION at $PATH and updates metadata
