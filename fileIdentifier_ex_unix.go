// +build !windows

package fileidentifier

import (
	"math/big"
	"os"
)

// fileIdentEx struct
type fileIdentEx struct {
	device uint64
	inode  uint64
}

// GetDeviceID returns the device id
func (f *fileIdentEx) GetDeviceID() uint64 {
	return f.device
}

// GetFileID returns the file id
func (f *fileIdentEx) GetFileID() *big.Int {
	return getBigInt(f.inode, 0)
}

// GetGlobalFileID returns the file id
func (f *fileIdentEx) GetGlobalFileID() *big.Int {
	n := getBigInt(f.device, 64)
	n.Add(n, getBigInt(f.inode, 0))
	return n
}

// GetFileIdentifierFromGetGlobalFileID returns a FileIdentifier by a GlobalFileID
func GetFileIdentifierFromGetGlobalFileIDEx(n *big.Int) FileIdentEx {
	info := GetFileIdentifierFromGetGlobalFileID(n)
	return &fileIdentEx{
		device: info.(*fileIdentifier).device,
		inode:  info.(*fileIdentifier).inode,
	}
}

// GetFileIdentifierByFile method
func GetFileIdentifierByFileEx(f *os.File) (FileIdentEx, error) {
	info, err := GetFileIdentifierByFile(f)
	if err != nil {
		return nil, err
	}
	return &fileIdentEx{
		device: info.(*fileIdentifier).device,
		inode:  info.(*fileIdentifier).inode,
	}, nil
}
