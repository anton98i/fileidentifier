package fileidentifier

import (
	"testing"
)

func TestGetID(t *testing.T) {
	testArr := []struct {
		vol              uint32
		idxHi            uint32
		idxLo            uint32
		expectedGlobalID string
		expectedFileID   uint64
	}{
		// result can get calculated by a full precision calculator like: https://www.mathsisfun.com/calculator-precision.html
		{
			vol:   1234,
			idxHi: 5678,
			idxLo: 90,
			// (1234 * 2^64) + (5678 * 2^32)+90
			expectedGlobalID: "22763282211344411000922",
			expectedFileID:   24386824306778,
		}, {
			vol:   4294967294,
			idxHi: 4294967293,
			idxLo: 4294967292,
			// (4294967294 * 2^64) + (4294967293 * 2^32)+4294967292
			expectedGlobalID: "79228162495817593511244464124",
			expectedFileID:   18446744065119617020,
		},
	}
	f := &fileIdentifier{}

	for _, te := range testArr {
		f.vol = te.vol
		f.idxHi = te.idxHi
		f.idxLo = te.idxLo

		testIDsBasic(t, f, te.expectedGlobalID, te.expectedFileID, uint64(te.vol))
	}
}
