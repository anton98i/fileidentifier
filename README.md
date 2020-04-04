# FileIdentifier


## Exported Types
2 Ways to get FileIdentifier:

Get by File safer way
``` go
func GetFileIdentifierByFile(f *os.File) (*FileIdentifier, error)
```
Faster one that uses the values set by fileInfo already
(but if a error occurs, the ids are 0 in windows: https://golang.org/src/os/types_windows.go#L216)
``` go
func GetFileIdentifier(i os.FileInfo) FileIdentifier
```

the return FileIdentifier has the following attributes:
``` go
func (f FileIdentifier) GetFileID() *big.Int
```
GetGlobalFileID returns 64 bit at unix/windows, 128 bit on windows with refs flag.
``` go
func (f FileIdentifier) GetDeviceID() uint64
```
GetDeviceID on windows without refs flag just returns a uint32 casted as a uint64.
``` go
func (f FileIdentifier) GetGlobalFileID() *big.Int
```
GetGlobalFileID returns 128 bit at unix/windows, 192 bit on windows with refs flag.


The return value of GetGlobalFileID can get used to create the FileIdentifier again. Needs to be same os to work correctly.
``` go
func GetFileIdentifierFromGetGlobalFileID(n *big.Int) FileIdentifier
```

## build with windows GetFileInformationByHandleEx
Building with "refs" flag uses the GetFileInformationByHandleEx to get 128 bit file id as used at ReFS filesystem.

The function "GetFileIdentifierByFile" needs to be used to get 128bit, "GetFileIdentifier" uses the go values that are 64 bit (they just get casted as 128 bit).

Go is not providing a function to us => much copy from go internal folder syscall, not really nice, as it may change later but there is no other way.

Additionally constants got changed, so golint to not give errors => do not use the constants.

The file "customFileInfo_windows.go" got created to use GetFileInformationByHandleEx to get file IDs.

```bash
go build -tags refs
```
