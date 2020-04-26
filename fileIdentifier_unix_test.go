// +build !windows

package fileidentifier

import (
	"strconv"
	"testing"
)

var testArr []struct {
	device           uint64
	inode            uint64
	expectedGlobalID string
	expectedFileID   uint64
}

func init() {
	testArr = []struct {
		device           uint64
		inode            uint64
		expectedGlobalID string
		expectedFileID   uint64
	}{
		// result can get calculated by a full precision calculator like: https://www.mathsisfun.com/calculator-precision.html
		{
			device: 1234,
			inode:  5678,
			// (1234 * 2^64) + 5678
			expectedGlobalID: "22763282186957586699822",
			expectedFileID:   5678,
		}, {
			device: 18446744073709551614,
			inode:  18446744073709551613,
			// (18446744073709551614 * 2^64) + 18446744073709551613
			expectedGlobalID: "340282366920938463444927863358058659837",
			expectedFileID:   18446744073709551613,
		},
	}
}

func TestGetID(t *testing.T) {
	f := &fileIdentifier{}

	for _, te := range testArr {
		f.device = te.device
		f.inode = te.inode

		testIDsBasic(t, f, te.expectedGlobalID, te.expectedFileID, te.device)
	}
}

func TestGetIDEx(t *testing.T) {
	f2 := &fileIdentifier{}
	f := &fileIdentEx{f: f2}

	for _, te := range testArr {
		f2.device = te.device
		f2.inode = te.inode

		testIDsBasicEx(t, f, te.expectedGlobalID, strconv.FormatUint(te.expectedFileID, 10), te.device)
	}
}
