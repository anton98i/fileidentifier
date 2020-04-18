package fileidentifier

import (
	"fmt"
	"math/big"
	"os"
	"syscall"
	"unsafe"

	"github.com/anton98i/fileIdentifier/internal/windows"
)

// fileIdentEx struct
type fileIdentEx struct {
	vol   uint64
	idxHi uint64
	idxLo uint64
}

// GetGlobalFileID returns the file id (all ids added to one big.Int)
func (f *fileIdentEx) GetGlobalFileID() *big.Int {
	n := getBigInt(uint64(f.vol), 128)
	n.Add(n, getBigInt(uint64(f.idxHi), 64))
	n.Add(n, getBigInt(uint64(f.idxLo), 0))
	return n
}

// GetFileIdentifierFromGetGlobalFileIDEx returns a FileIdentifier by a GlobalFileID
func GetFileIdentifierFromGetGlobalFileIDEx(n *big.Int) FileIdentEx {
	tmpPtr := new(big.Int)
	tmpPtr.Set(n)
	var resultTmp big.Int
	// 18446744073709551615 = 64 bits filled with 1
	andOperator := getBigInt(18446744073709551615, 0)
	idxLo := resultTmp.And(tmpPtr, andOperator).Uint64()
	idxHi := resultTmp.And(tmpPtr.Rsh(tmpPtr, 64), andOperator).Uint64()
	vol := resultTmp.And(tmpPtr.Rsh(tmpPtr, 64), andOperator).Uint64()

	return &fileIdentEx{
		vol:   vol,
		idxHi: idxHi,
		idxLo: idxLo,
	}
}

// GetDeviceID returns the device id
func (f *fileIdentEx) GetDeviceID() uint64 {
	return f.vol
}

// GetFileID returns the file id
func (f *fileIdentEx) GetFileID() *big.Int {
	n := getBigInt(f.idxHi, 64)
	n.Add(n, getBigInt(f.idxLo, 0))
	return n
}

// GetFileIdentifierByFileEx method
func GetFileIdentifierByFileEx(f *os.File) (FileIdentEx, error) {
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
	return &fileIdentEx{
		vol:   uint64(uint32(ti.VolumeSerialNumber)),
		idxHi: high,
		idxLo: low,
	}, nil
}
