[![Build Status](https://travis-ci.org/anton98i/fileIdentifier.svg?branch=master)](https://travis-ci.org/anton98i/fileIdentifier)
# FileIdentifier
FileIdentifier as a go module that read the ID/Device ID of a a file on linux/windows as thais is a os specific operation.

# install dependencies
For windows we need the "golang.org/x/sys" package:
``` bash
go get ./..
```

## Exported Functions to get a file identifier
Ways to get a FileIdentifier:

``` go
func GetFileIdentifierByFile(f *os.File) (FileIdentifier, error)
```
``` go
func GetFileIdentifierByFileEx(f *os.File) (FileIdentEx, error)
```


A FileIdentifier can also be received by a os.FileInfo, but that uses at windows private fields which go might change later.
``` go
func GetFileIdentifier(i os.FileInfo) (FileIdentifier, error)
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
*  FileIdentEx uses a uint64 device ID instead of uint32
*  FileIdentEx uses a "uint128" file ID instead of a uint64 => according to windows documentation are 128 bit IDs used at ReFS

## Exported Functions to get a file identifier
2 Ways to get FileIdentifier:


GetFileIdentifierByFile

The return value of GetGlobalFileID can get used to create the FileIdentifier again. Needs to be on same os to work correctly.
``` go
func GetFileIdentifierFromGetGlobalFileID(n *big.Int) FileIdentifier
```
``` go
func GetFileIdentifierFromGetGlobalFileIDEx(n *big.Int) FileIdentEx
```
