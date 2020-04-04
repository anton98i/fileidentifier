// +build !refs

package fileidentifier

import (
	"fmt"
	"math/big"
	"os"
	"reflect"
	"syscall"
)

// FileIdentifier struct
type FileIdentifier struct {
	vol   uint32
	idxHi uint32
	idxLo uint32
}

// GetGlobalFileID returns the file id (all ids added to one big.Int)
func (f FileIdentifier) GetGlobalFileID() *big.Int {
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

	return FileIdentifier{
		vol:   uint32(vol),
		idxHi: uint32(idxHi),
		idxLo: uint32(idxLo),
	}
}

// GetDeviceID returns the device id
func (f FileIdentifier) GetDeviceID() uint64 {
	return uint64(f.vol)
}

// GetFileID returns the file id
func (f FileIdentifier) GetFileID() *big.Int {
	return getBigInt((uint64(f.idxHi)<<32)+uint64(f.idxLo), 0)
}

// GetFileIdentifierByFile method
func GetFileIdentifierByFile(f *os.File) (*FileIdentifier, error) {
	// get file id's like in the go implementation newFileStatFromGetFileInformationByHandle https://golang.org/src/os/types_windows.go#L44
	var d syscall.ByHandleFileInformation
	err := syscall.GetFileInformationByHandle(syscall.Handle(f.Fd()), &d)

	if err != nil {
		return nil, fmt.Errorf("GetFileInformationByHandle error: %v", err.Error())
	}

	return &FileIdentifier{
		vol:   d.VolumeSerialNumber,
		idxHi: d.FileIndexHigh,
		idxLo: d.FileIndexLow,
	}, nil
}

// GetFileIdentifier returns the platform specific FileIdentifier
func GetFileIdentifier(i os.FileInfo) FileIdentifier {
	// according to that is the file id in fileInfo stored as private value that are used for samefile
	// https://golang.org/src/os/types_windows.go#L65
	// according to that gets the file information allways loaded and skipped if already done =>
	// https://golang.org/src/os/types_windows.go#L216
	// => call SameFile to make the the values are set
	//
	// it is already filled called with stat method of a file it https://golang.org/src/os/stat_windows.go#L15
	os.SameFile(i, i)

	// get fileStat through reflection as otherwise they not accessible because they are private
	// https://golang.org/src/os/types_windows.go#L65
	fileStat := reflect.ValueOf(i).Elem()

	// https://docs.microsoft.com/en-us/windows/win32/api/fileapi/ns-fileapi-by_handle_file_information
	// According to that are the identifier 128 bit (64bit each) on refs, but go uses the "old" api
	// needed go fields are here defined: https://golang.org/src/os/types_windows.go#L16
	// Uint returns a uint64, but the values are actually 32 => safe to cast as uint32
	fileState := FileIdentifier{
		idxHi: uint32(fileStat.FieldByName("idxhi").Uint()),
		idxLo: uint32(fileStat.FieldByName("idxlo").Uint()),
		vol:   uint32(fileStat.FieldByName("vol").Uint()),
	}
	return fileState
}

// iterateAllFileIdentifier is for for tests
func iterateAllFileIdentifier(cb func(globalId, expectedFileID *big.Int, vol, idxHi, idxLo uint64, f FileIdentifier)) {
	f := FileIdentifier{}

	expected := big.NewInt(0)
	iterateAllUint64(4294967295, func(vol uint64) {
		bigIdxVol := getBigInt(vol, 64)
		expected.Add(expected, bigIdxVol)
		iterateAllUint64(4294967295, func(idxHi uint64) {
			bigIdxHi := getBigInt(idxHi, 32)
			expected.Add(expected, bigIdxHi)
			iterateAllUint64(4294967295, func(idxLo uint64) {
				f.vol = uint32(vol)
				f.idxHi = uint32(idxHi)
				f.idxLo = uint32(idxLo)

				bigIdxLo := getBigInt(idxLo, 0)
				expected.Add(expected, bigIdxLo)

				expectedFileID := big.NewInt(0)
				expectedFileID.Add(bigIdxHi, bigIdxLo)

				cb(expected, expectedFileID, vol, idxHi, idxLo, f)

				expected.Sub(expected, bigIdxLo)
			})
			expected.Sub(expected, bigIdxHi)
		})
		expected.Sub(expected, bigIdxVol)
	})
}
