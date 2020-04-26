// +build !windows

package fileidentifier

import (
	"math/big"
)

// fileIdentEx struct
type fileIdentEx struct {
	f FileIdentifier
}

// GetDeviceID returns the device id
func (f *fileIdentEx) GetDeviceID() uint64 {
	return f.f.GetDeviceID()
}

// GetFileID returns the file id
func (f *fileIdentEx) GetFileID() *big.Int {
	return getBigInt(f.f.GetFileID(), 0)
}

// GetGlobalFileID returns the file id
func (f *fileIdentEx) GetGlobalFileID() *big.Int {
	return f.f.GetGlobalFileID()
}

// GetFileIdentifierFromGetGlobalFileID returns a FileIdentifier by a GlobalFileID
func GetFileIdentifierFromGetGlobalFileIDEx(n *big.Int) FileIdentEx {
	info := GetFileIdentifierFromGetGlobalFileID(n)
	return &fileIdentEx{
		f: info,
	}
}

// GetFileIdentifierByPathEx method
func GetFileIdentifierByPathEx(path string) (FileIdentEx, error) {
	info, err := GetFileIdentifierByPath(path)
	if err != nil {
		return nil, err
	}
	return &fileIdentEx{
		f: info,
	}, nil
}
