# AR

AR is a Go subpackage for reading AR files. 

## Example
Example can be found [here](../example/README.md)
```
go run main.go ar file/path/to/foo.ar
```

## Dev
### Golden Image
```
ar rcs ../ar/golden-archive.ar README.md bar/baz bar/symlink chew foo
```

## Credit
Cloned from https://github.com/erikgeiser/ar
