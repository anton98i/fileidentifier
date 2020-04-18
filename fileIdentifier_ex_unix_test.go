// +build !windows

package fileidentifier

import (
	"math/big"
	"testing"
)

func iterateAllFileIdentifierEx(cb func(globalId, expectedFileID *big.Int, dev, inode uint64)) {
	expected := big.NewInt(0)
	iterateAllUint64(18446744073709551615, func(dev uint64) {
		devBig := getBigInt(dev, 64)
		expected.Add(expected, devBig)
		iterateAllUint64(18446744073709551615, func(inode uint64) {
			inodeBig := getBigInt(inode, 0)
			expected.Add(expected, inodeBig)

			cb(expected, inodeBig, dev, inode)

			expected.Sub(expected, inodeBig)
		})
		expected.Sub(expected, devBig)
	})
}

func checkFileIdentifierBasicEx(t *testing.T, _f, _expected FileIdentEx) {
	f := _f.(*fileIdentEx)
	expected := _expected.(*fileIdentEx)
	if f.device != expected.device {
		t.Errorf("checkFileIdentifierBasic vol failed, expected: %d, received: %d", expected.device, f.device)
	}
	if f.inode != expected.inode {
		t.Errorf("checkFileIdentifierBasic idxHi failed, expected: %d, received: %d", expected.inode, f.inode)
	}
	if f.GetFileID().Cmp(expected.GetFileID()) != 0 {
		t.Errorf("f.Cmp(expected.GetFileID()) != 0 , expected: %s, received: %s", expected.GetFileID().String(), f.GetFileID().String())
	}
}

func TestGetIDAllPossibleValuesUnixEx(t *testing.T) {
	f := &fileIdentEx{}

	iterateAllFileIdentifierEx(func(expected, expectedFileID *big.Int, dev, inode uint64) {
		f.device = dev
		f.inode = inode

		if f.GetDeviceID() != dev {
			t.Errorf("f.GetDeviceID() != dev, expected: %d, received: %d", dev, f.GetDeviceID())
		}

		if expected.Cmp(f.GetGlobalFileID()) != 0 {
			t.Errorf("expected.Cmp(f.GetGlobalFileID()) != 0, expected: %s, received: %s", expected.String(), f.GetGlobalFileID().String())
		}

		if expectedFileID.Cmp(f.GetFileID()) != 0 {
			t.Errorf("expectedFileID.Cmp(f.GetFileID()) != 0 , expected: %s, received: %s", expectedFileID.String(), f.GetFileID().String())
		}

		checkFileIdentifierBasicEx(t, GetFileIdentifierFromGetGlobalFileIDEx(f.GetGlobalFileID()), f)
	})
}
