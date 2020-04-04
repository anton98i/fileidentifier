// +build windows,refs

package fileidentifier

import (
	"fmt"
	"math/big"
	"os"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/anton98i/FileIdentifier/internal/windows"
)

// FileIdentifier struct
type FileIdentifier struct {
	vol   uint64
	idxHi uint64
	idxLo uint64
}

// GetGlobalFileID returns the file id (all ids added to one big.Int)
func (f FileIdentifier) GetGlobalFileID() *big.Int {
	n := getBigInt(uint64(f.vol), 128)
	n.Add(n, getBigInt(uint64(f.idxHi), 64))
	n.Add(n, getBigInt(uint64(f.idxLo), 0))
	return n
}

// GetFileIdentifierFromGetGlobalFileID returns a FileIdentifier by a GlobalFileID
func GetFileIdentifierFromGetGlobalFileID(n *big.Int) FileIdentifier {
	tmpPtr := new(big.Int)
	tmpPtr.Set(n)
	var resultTmp big.Int
	// 18446744073709551615 = 64 bits filled with 1
	andOperator := getBigInt(18446744073709551615, 0)
	idxLo := resultTmp.And(tmpPtr, andOperator).Uint64()
	idxHi := resultTmp.And(tmpPtr.Rsh(tmpPtr, 64), andOperator).Uint64()
	vol := resultTmp.And(tmpPtr.Rsh(tmpPtr, 64), andOperator).Uint64()

	return FileIdentifier{
		vol:   vol,
		idxHi: idxHi,
		idxLo: idxLo,
	}
}

// GetDeviceID returns the device id
func (f FileIdentifier) GetDeviceID() uint64 {
	return uint64(f.vol)
}

// GetFileID returns the file id
func (f FileIdentifier) GetFileID() *big.Int {
	n := getBigInt(f.idxHi, 64)
	n.Add(n, getBigInt(f.idxLo, 0))
	return n
}

// GetFileIdentifierByFile method
func GetFileIdentifierByFile(f *os.File) (*FileIdentifier, error) {
	// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-getfileinformationbyhandleex
	// call with FileIdInfo (0x12) to get indexes: https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-getfileinformationbyhandleex#remarks
	// https://docs.microsoft.com/de-de/windows/win32/api/winbase/ns-winbase-file_id_info

	// go definitions:
	// windows.FileIdInfo: https://golang.org/src/internal/syscall/windows/symlink_windows.go
	// GetFileInformationByHandleEx: https://golang.org/src/internal/syscall/windows/zsyscall_windows.go
	var ti windows.FileIDInfoStrcut
	err := windows.GetFileInformationByHandleEx(syscall.Handle(f.Fd()), windows.FileIDInfo, (*byte)(unsafe.Pointer(&ti)), uint32(unsafe.Sizeof(ti)))

	if err != nil {
		return nil, fmt.Errorf("GetFileInformationByHandleEx error: %v", err.Error())
	}
	high, low := ti.GetFileID()
	ret := &FileIdentifier{
		vol:   uint64(uint32(ti.VolumeSerialNumber)),
		idxHi: high,
		idxLo: low,
	}
	return ret, nil
}

// GetFileIdentifier returns the platform specific FileIdentifier
func GetFileIdentifier(i os.FileInfo) FileIdentifier {
	// in this scenario it is not possible to get the full file id as go uses the "old" GetFileInformationByHandle and not the new GetFileInformationByHandleEx

	// according to that are the ids are used for samefile
	// https://golang.org/src/os/types_windows.go#L65
	// according to that gets the file information allways loaded and skipped if already done =>
	// https://golang.org/src/os/types_windows.go#L216
	// => call SameFile to make the the values are set
	os.SameFile(i, i)

	// Gathering fileStat through reflection as otherwise not accessible
	// https://golang.org/src/os/types_windows.go#L65
	fileStat := reflect.ValueOf(i).Elem()

	// https://docs.microsoft.com/en-us/windows/win32/api/fileapi/ns-fileapi-by_handle_file_information
	// According to that are the identifier 128 bit (64bit each) on refs, but go uses the "old" api
	// needed go fields are here defined: https://golang.org/src/os/types_windows.go#L16
	// Uint returns a uint64, but the values are actually 32 => safe to cast as uint32
	// to get the full id use GetFileIdentifierByFile instead that uses GetFileInformationByHandleEx
	return FileIdentifier{
		idxHi: 0,
		idxLo: (fileStat.FieldByName("idxhi").Uint() << 32) + fileStat.FieldByName("idxlo").Uint(),
		vol:   fileStat.FieldByName("vol").Uint(),
	}
}

// iterateAllFileIdentifier is for for tests
func iterateAllFileIdentifier(cb func(globalId, expectedFileID *big.Int, vol, idxHi, idxLo uint64, f FileIdentifier)) {
	f := FileIdentifier{}

	expected := big.NewInt(0)
	iterateAllUint64(18446744073709551615, func(vol uint64) {
		bigIdxVol := getBigInt(vol, 128)
		expected.Add(expected, bigIdxVol)
		iterateAllUint64(18446744073709551615, func(idxHi uint64) {
			bigIdxHi := getBigInt(idxHi, 64)
			expected.Add(expected, bigIdxHi)
			iterateAllUint64(18446744073709551615, func(idxLo uint64) {
				f.vol = vol
				f.idxHi = idxHi
				f.idxLo = idxLo

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
