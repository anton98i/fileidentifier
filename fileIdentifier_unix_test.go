// +build !windows

package fileidentifier

import (
	"math/big"
	"testing"
)

func iterateAllFileIdentifier(cb func(globalId *big.Int, expectedFileID, dev, inode uint64)) {
	expected := big.NewInt(0)
	iterateAllUint64(18446744073709551615, func(dev uint64) {
		devBig := getBigInt(dev, 64)
		expected.Add(expected, devBig)
		iterateAllUint64(18446744073709551615, func(inode uint64) {
			inodeBig := getBigInt(inode, 0)
			expected.Add(expected, inodeBig)

			cb(expected, inode, dev, inode)

			expected.Sub(expected, inodeBig)
		})
		expected.Sub(expected, devBig)
	})
}

func checkFileIdentifierBasic(t *testing.T, _f, _expected FileIdentifier) {
	f := _f.(*fileIdentifier)
	expected := _expected.(*fileIdentifier)
	if f.device != expected.device {
		t.Errorf("checkFileIdentifierBasic vol failed, expected: %d, received: %d", expected.device, f.device)
	}
	if f.inode != expected.inode {
		t.Errorf("checkFileIdentifierBasic idxHi failed, expected: %d, received: %d", expected.inode, f.inode)
	}
	if f.GetFileID() != expected.GetFileID() {
		t.Errorf("checkFileIdentifierBasic GetFileID failed, expected: %d, received: %d", expected.GetFileID(), f.GetFileID())
	}
}

func TestGetIDAllPossibleValuesUnix(t *testing.T) {
	f := &fileIdentifier{}

	iterateAllFileIdentifier(func(expected *big.Int, expectedFileID, dev, inode uint64) {
		f.device = dev
		f.inode = inode

		if f.GetDeviceID() != dev {
			t.Errorf("f.GetDeviceID() != dev, expected: %d, received: %d", dev, f.GetDeviceID())
		}

		if expected.Cmp(f.GetGlobalFileID()) != 0 {
			t.Errorf("expected.Cmp(f.GetGlobalFileID()) != 0, expected: %s, received: %s", expected.String(), f.GetGlobalFileID().String())
		}

		if expectedFileID != f.GetFileID() {
			t.Errorf("expectedFileID != f.GetFileID(), expected: %d, received: %d", expectedFileID, f.GetFileID())
		}

		checkFileIdentifierBasic(t, GetFileIdentifierFromGetGlobalFileID(f.GetGlobalFileID()), f)
	})
}
