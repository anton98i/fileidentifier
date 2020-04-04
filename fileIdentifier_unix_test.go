// +build !windows,!integration

package fileidentifier

import (
	"math/big"
	"testing"
)

// getind dev/inode not tested as i trust this os :)

func iterateAllFileIdentifier(cb func(globalId, expectedFileID *big.Int, dev, inode uint64)) {
	expected := big.NewInt(0)
	iterateAllUint64(0, func(dev uint64) {
		devBig := getBigInt(dev, 64)
		expected.Add(expected, devBig)
		iterateAllUint64(0, func(inode uint64) {
			inodeBig := getBigInt(inode, 0)
			expected.Add(expected, inodeBig)

			cb(expected, inodeBig, dev, inode)

			expected.Sub(expected, inodeBig)
		})
		expected.Sub(expected, devBig)
	})
}

func TestGetIDAllPossibleValuesUnix(t *testing.T) {
	f := FileIdentifier{}

	iterateAllFileIdentifier(func(globalId, expectedFileID *big.Int, dev, inode uint64) {
		f.device = dev
		f.inode = inode

		if f.GetDeviceID() != vol {
			t.Errorf("compare failed, expected: %d, received: %d", vol, f.GetDeviceID())
		}

		if expected.Cmp(f.GetGlobalFileID()) != 0 {
			t.Errorf("compare failed, expected: %s, received: %s", expected.String(), f.GetID().String())
		}

		if inodeBig.Cmp(f.GetFileID()) != 0 {
			t.Errorf("compare failed, expected: %s, received: %s", inodeBig.String(), f.GetFileID().String())
		}

		newF := GetFileIdentifierFromGetGlobalFileID(f.GetGlobalFileID())
		if newF.device != device {
			t.Errorf("GetFileIdentifierFromGetGlobalFileID device failed, expected: %d, received: %d", device, newF.device)
		}
		if newF.inode != inode {
			t.Errorf("GetFileIdentifierFromGetGlobalFileID inode failed, expected: %d, received: %d", inode, newF.inode)
		}
	})
}
