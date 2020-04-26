// +build !windows

package fileidentifier

import (
	"fmt"
	"math/big"
	"os"
	"syscall"
)

// fileIdentifier struct
type fileIdentifier struct {
	device uint64
	inode  uint64
}

// GetDeviceID returns the device id
func (f *fileIdentifier) GetDeviceID() uint64 {
	return f.device
}

// GetFileID returns the file id
func (f *fileIdentifier) GetFileID() uint64 {
	return f.inode
}

// GetGlobalFileID returns the file id
func (f *fileIdentifier) GetGlobalFileID() *big.Int {
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

	return &fileIdentifier{
		device: device,
		inode:  inode,
	}
}

// GetFileIdentifierByPath method
func GetFileIdentifierByPath(path string) (FileIdentifier, error) {
	stats, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	return getFileIdentifier(stats)
}

// getFileIdentifier returns the platform specific FileIdentifier
func getFileIdentifier(i os.FileInfo) (FileIdentifier, error) {

	/* not necessary
	if !os.SameFile(i, i) {
		return nil, fmt.Errorf("error getting ids")
	}
	*/

	stat, ok := i.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("Not a syscall.Stat_t")
	}

	// Get the two fields required to uniquely identify file
	// https://golang.org/pkg/syscall/#Stat_t
	return &fileIdentifier{
		device: uint64(stat.Dev),
		inode:  uint64(stat.Ino),
	}, nil
}
