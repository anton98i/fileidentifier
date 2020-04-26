[![Build Status](https://travis-ci.org/anton98i/fileIdentifier.svg?branch=master)](https://travis-ci.org/anton98i/fileIdentifier)
[![GoDoc](https://godoc.org/github.com/thylong/ian?status.svg)](https://godoc.org/github.com/anton98i/fileidentifier)
# FileIdentifier
FileIdentifier as a go module that read the ID/Device ID of a a file on linux/windows as thais is a os specific operation.

# install dependencies
For windows we need the "golang.org/x/sys" package:
``` bash
go get ./..
```

## Exported interfaces
``` go
// FileIdentifier interface
type FileIdentifier interface {
	// GetGlobalFileID returns the device id + file id combined to one id (a "uint128")
	GetGlobalFileID() *big.Int

	// GetDeviceID returns the device id (on windows it is a uint32 casted as uint64)
	GetDeviceID() uint64

	// GetFileID returns the file id
	GetFileID() uint64
}

// FileIdentEx interface
type FileIdentEx interface {
	// GetGlobalFileID returns the device id + file id combined to one id (a "uint192")
	GetGlobalFileID() *big.Int

	// GetDeviceID returns the device id
	GetDeviceID() uint64

	// GetFileID returns the file id as a "uint128"
	GetFileID() *big.Int
}
```

The difference between these two types is only at windows a difference:
*  FileIdentEx uses a uint64 device ID instead of uint32 => device id is cutted at windows :/
*  FileIdentEx uses a "uint128" file ID instead of a uint64 => according to windows documentation are 128 bit IDs used at ReFS


## Exported Functions to get a file identifier
Ways to get a FileIdentifier:

``` go
func GetFileIdentifierByPath(path string) (FileIdentifier, error)
```

Ways to get a "extended" FileIdentifier:
``` go
func GetFileIdentifierByPathEx(path string) (FileIdentifier, error)
```


## Run build
``` bash
go build
```

different os:
``` bash
# windows and linus 64 bit
env GOOS=windows GOARCH=amd64 go build
env GOOS=linux GOARCH=amd64 go build

# windows and linux 32 bit
env GOOS=windows GOARCH=386 go build
env GOOS=linux GOARCH=386 go build
```

## Run tests
runs the tests
``` bash
go test -v ./...
```

different os (test itself won't run, only if compile of the test works is checked):
``` bash
# windows and linus 64 bit
env GOOS=windows GOARCH=amd64 go test -v ./...
env GOOS=linux GOARCH=amd64 go test -v ./...

# windows and linux 32 bit
env GOOS=windows GOARCH=386 go test -v ./...
env GOOS=linux GOARCH=386 go test -v ./...
```
