// +build !windows

package fileidentifier

import (
	"math/big"
	"os"
	"syscall"
)

// FileIdentifier struct
type FileIdentifier struct {
	device uint64
	inode  uint64
}

// GetDeviceID returns the device id
func (f FileIdentifier) GetDeviceID() uint64 {
	return f.device
}

// GetFileID returns the file id
func (f FileIdentifier) GetFileID() *big.Int {
	return getBigInt(f.inode, 0)
}

// GetGlobalFileID returns the file id
func (f FileIdentifier) GetGlobalFileID() *big.Int {
	n := getBigInt(f.device, 64)
	n.Add(n, getBigInt(f.inode, 0))
	return n
}

// GetFileIdentifierFromGetGlobalFileID returns a FileIdentifier by a GlobalFileID
func GetFileIdentifierFromGetGlobalFileID(n *big.Int) FileIdentifier {
	tmpPtr := new(big.Int)
	tmpPtr.Set(n)
	var resultTmp big.Int
	// 18446744073709551615 = 64 bits filled with 1
	andOperator := getBigInt(18446744073709551615, 0)
	inode := resultTmp.And(tmpPtr, andOperator).Uint64()
	device := resultTmp.And(tmpPtr.Rsh(tmpPtr, 64), andOperator).Uint64()

	return FileIdentifier{
		device: device,
		inode:  inode,
	}
}

// GetFileIdentifierByFile method
func GetFileIdentifierByFile(f *os.File) (*FileIdentifier, error) {
	stats, err := f.Stat()
	if err != nil {
		return nil, err
	}
	ret := GetFileIdentifier(stats)
	return &ret, nil
}

// GetFileIdentifier returns the platform specific FileIdentifier
func GetFileIdentifier(i os.FileInfo) FileIdentifier {

	/* not always necessary
	// make sure thd ids are set to force setting them by compared the file (this checks the ids)
	if !os.SameFile(i, i) {
		// no should not happen
		panic("os.SameFile is not the same file for the same file info")
	}
	*/

	stat := i.Sys().(*syscall.Stat_t)

	// Get the two fields required to uniquely identify file
	// https://golang.org/pkg/syscall/#Stat_t
	fileState := FileIdentifier{
		device: uint64(stat.Dev),
		inode:  uint64(stat.Ino),
	}

	return fileState
}
