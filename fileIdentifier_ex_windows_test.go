package fileidentifier

import (
	"testing"
)

func TestGetIDEx(t *testing.T) {
	testArr := []struct {
		vol              uint64
		idxHi            uint64
		idxLo            uint64
		expectedGlobalID string
		expectedFileID   string
	}{
		// result can get calculated by a full precision calculator like: https://www.mathsisfun.com/calculator-precision.html
		{
			vol:   1234,
			idxHi: 5678,
			idxLo: 90,
			// (1234 * 2^128) + (5678 * 2^64)+90
			expectedGlobalID: "419908440780438064018544878421324807012442",
			expectedFileID:   "104740612850522834075738",
		},
		{
			vol:   18446744073709551614,
			idxHi: 18446744073709551613,
			idxLo: 18446744073709551612,
			// (18446744073709551614 * 2^128) + (18446744073709551613 * 2^64)+18446744073709551612
			expectedGlobalID: "6277101735386680763495507056286727952602087348884847198204",
			expectedFileID:   "340282366920938463426481119284349108220",
		},
	}
	f := &fileIdentEx{}

	for _, te := range testArr {
		f.vol = te.vol
		f.idxHi = te.idxHi
		f.idxLo = te.idxLo

		testIDsBasicEx(t, f, te.expectedGlobalID, te.expectedFileID, te.vol)
	}
}
