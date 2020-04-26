package fileidentifier

import (
	"fmt"
	"math/big"
	"os"
	"reflect"
	"syscall"
)

type fileIdentifier struct {
	vol   uint32
	idxHi uint32
	idxLo uint32
}

// GetGlobalFileID returns the file id (all ids added to one big.Int)
func (f *fileIdentifier) GetGlobalFileID() *big.Int {
	n := getBigInt(uint64(f.vol), 64)
	n.Add(n, getBigInt(uint64(f.idxHi), 32))
	n.Add(n, getBigInt(uint64(f.idxLo), 0))
	return n
}

// GetFileIdentifierFromGetGlobalFileID returns a FileIdentifier by a GlobalFileID
func GetFileIdentifierFromGetGlobalFileID(n *big.Int) FileIdentifier {
	tmpPtr := new(big.Int)
	tmpPtr.Set(n)
	var resultTmp big.Int
	// 4294967295 = 32 bits filled with 1
	andOperator := getBigInt(4294967295, 0)
	idxLo := resultTmp.And(tmpPtr, andOperator).Uint64()
	idxHi := resultTmp.And(tmpPtr.Rsh(tmpPtr, 32), andOperator).Uint64()
	vol := resultTmp.And(tmpPtr.Rsh(tmpPtr, 32), andOperator).Uint64()

	return &fileIdentifier{
		vol:   uint32(vol),
		idxHi: uint32(idxHi),
		idxLo: uint32(idxLo),
	}
}

// GetDeviceID returns the device id
func (f fileIdentifier) GetDeviceID() uint64 {
	return uint64(f.vol)
}

// GetFileID returns the file id
func (f fileIdentifier) GetFileID() uint64 {
	return uint64(f.idxHi)<<32 + uint64(f.idxLo)
}

// GetFileIdentifierByPath method
func GetFileIdentifierByPath(path string) (FileIdentifier, error) {
	h, err := getHandleFromPath(path)
	if err != nil {
		return nil, fmt.Errorf("getHandleFromPath error: %v", err)
	}
	defer syscall.CloseHandle(h)
	return getFileIdentifierByHandle(h)
}

func getFileIdentifierByHandle(handle syscall.Handle) (FileIdentifier, error) {
	// get file id's like in the go implementation newFileStatFromGetFileInformationByHandle https://golang.org/src/os/types_windows.go#L44
	var d syscall.ByHandleFileInformation
	err := syscall.GetFileInformationByHandle(handle, &d)

	if err != nil {
		return nil, fmt.Errorf("GetFileInformationByHandle error: %v", err.Error())
	}

	return &fileIdentifier{
		vol:   d.VolumeSerialNumber,
		idxHi: d.FileIndexHigh,
		idxLo: d.FileIndexLow,
	}, nil
}

// getFileIdentifier returns the platform specific FileIdentifier
// not reliable (for lang paths as loadFileIds do not use fixLongPath) => not used for now, maybe later
func getFileIdentifier(i os.FileInfo) (FileIdentifier, error) {
	// according to that is the file id in fileInfo stored as private value that are used for samefile
	// https://golang.org/src/os/types_windows.go#L65
	// according to that gets the file information allways loaded and skipped if already done =>
	// https://golang.org/src/os/types_windows.go#L216
	// => call SameFile to make the the values are set
	// os.FileInfo is a interface, which is a pointer

	// it is already filled called with stat method of a file it https://golang.org/src/os/stat_windows.go#L15
	if !os.SameFile(i, i) {
		// the implementation of SameFile will return false if on any file occurred an error
		// a error also occurs, if the path is too lang..... not so good go implementation
		return nil, fmt.Errorf("error getting ids")
	}

	// get fileStat through reflection as they are not accessible string because they are private
	// https://golang.org/src/os/types_windows.go#L65
	fileStat := reflect.ValueOf(i).Elem()

	// https://docs.microsoft.com/en-us/windows/win32/api/fileapi/ns-fileapi-by_handle_file_information
	// According to that are the identifier 128 bit (64bit each) on refs, but go uses the "old" api
	// needed go fields are here defined: https://golang.org/src/os/types_windows.go#L16
	// Uint returns a uint64, but the values are actually 32 => safe to cast as uint32
	return &fileIdentifier{
		idxHi: uint32(fileStat.FieldByName("idxhi").Uint()),
		idxLo: uint32(fileStat.FieldByName("idxlo").Uint()),
		vol:   uint32(fileStat.FieldByName("vol").Uint()),
	}, nil
}
