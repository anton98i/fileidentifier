package fileidentifier

import (
	"fmt"
	"math/big"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
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

const fileIDInfo = 0x12 // FILE_ID_INFO

// fileIDInfoStrcut is representation of FILE_ID_INFO: https://docs.microsoft.com/de-de/windows/win32/api/winbase/ns-winbase-file_id_info
type fileIDInfoStrcut struct {
	// ULONGLONG = typedef unsigned __int64 ULONGLONG;: https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/c57d9fba-12ef-4853-b0d5-a6f472b50388
	// ULONGLONG   VolumeSerialNumber;
	VolumeSerialNumber uint64
	// FILE_ID_128 https://docs.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-file_id_128
	// typedef struct _FILE_ID_128 {
	//  BYTE Identifier[16];
	// }
	// FILE_ID_128 FileId;
	FileID struct {
		arr [16]byte
	}
}

func bytes2String(b []byte) uint64 {
	var ret uint64
	for i := uint64(0); i < 8; i++ {
		ret += uint64(b[i]) << (i * 8)
	}
	return ret
}

// GetFileID method
func (f fileIDInfoStrcut) GetFileID() (uint64, uint64) {
	return bytes2String(f.FileID.arr[8:16]), bytes2String(f.FileID.arr[:8])
}

// GetFileIdentifierByPathEx method
func GetFileIdentifierByPathEx(path string) (FileIdentEx, error) {
	h, err := getHandleFromPath(path)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(h)
	return getFileIdentifierByHandleEx(h)
}

func getFileIdentifierByHandleEx(h syscall.Handle) (FileIdentEx, error) {
	// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-getfileinformationbyhandleex
	// call with FileIdInfo (0x12) to get indexes: https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-getfileinformationbyhandleex#remarks
	// https://docs.microsoft.com/de-de/windows/win32/api/winbase/ns-winbase-file_id_info

	// go definitions:
	// windows.FileIdInfo: https://golang.org/src/internal/syscall/windows/symlink_windows.go
	// GetFileInformationByHandleEx in go src: https://golang.org/src/internal/syscall/windows/zsyscall_windows.go
	// used GetFileInformationByHandleEx in x/sys: https://pkg.go.dev/golang.org/x/sys@v0.0.0-20200413165638-669c56c373c4/windows?tab=doc#GetFileInformationByHandleEx
	var ti fileIDInfoStrcut
	err := windows.GetFileInformationByHandleEx(windows.Handle(h), fileIDInfo, (*byte)(unsafe.Pointer(&ti)), uint32(unsafe.Sizeof(ti)))

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
