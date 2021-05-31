# MVC - Mastro Version Control
A minimal data versioning tool in Golang

* `mvc init $PATH` - initializes metadata file at $PATH
* `mvc new $PATH` - creates new version and returns full path at $PATH
* `mvc add $LOCALPATH $PATH` - adds content at $LOCALPATH to remote $PATH at current latest version
* `mvc latest $PATH` - retrieves latest version at $PATH and returns full path
* `mvc overwrite $PATH $VERSION $LOCALPATH` - overwrite existing version $VERSION at $PATH and updates metadata

## Example

