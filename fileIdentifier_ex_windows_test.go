package fileidentifier

import (
	"math/big"
	"testing"
)

func checkFileIdentifierBasicEx(t *testing.T, _f, _expected FileIdentEx) {
	f := _f.(*fileIdentEx)
	expected := _expected.(*fileIdentEx)
	if f.vol != expected.vol {
		t.Errorf("checkFileIdentifierBasicEx vol failed, expected: %d, received: %d", expected.vol, f.vol)
	}
	if f.idxHi != expected.idxHi {
		t.Errorf("checkFileIdentifierBasicEx idxHi failed, expected: %d, received: %d", expected.idxHi, f.idxHi)
	}
	if f.idxLo != expected.idxLo {
		t.Errorf("checkFileIdentifierBasicEx idxLo failed, expected: %d, received: %d", expected.idxLo, f.idxLo)
	}
	if f.GetFileID().Cmp(expected.GetFileID()) != 0 {
		t.Errorf("checkFileIdentifierBasicEx GetFileID failed, expected: %d, received: %d", expected.GetFileID(), f.GetFileID())
	}
}

func TestGetIDEx(t *testing.T) {
	f := fileIdentEx{
		vol:   1234,
		idxHi: 5678,
		idxLo: 90,
	}
	id := f.GetGlobalFileID()
	// result can get calculated by a full precision calculator like: https://www.mathsisfun.com/calculator-precision.html
	// (1234 * 2^128) + (5678 * 2^64)+90id.String()
	if "419908440780438064018544878421324807012442" != id.String() {
		t.Errorf("419908440780438064018544878421324807012442 != id.String(); got %v", id.String())
	}
}

func TestGetID2Ex(t *testing.T) {
	// max: 18446744073709551615
	f := fileIdentEx{
		vol:   18446744073709551614,
		idxHi: 18446744073709551613,
		idxLo: 18446744073709551612,
	}
	id := f.GetGlobalFileID()
	// result can get calculated by a full precision calculator like: https://www.mathsisfun.com/calculator-precision.html
	// (18446744073709551614 * 2^128) + (18446744073709551613 * 2^64)+18446744073709551612
	if "6277101735386680763495507056286727952602087348884847198204" != id.String() {
		t.Errorf("6277101735386680763495507056286727952602087348884847198204 != id.String(); got %v", id.String())
	}
}
func TestGetFileIdentifierEx(t *testing.T) {
	file := getTestFile(t)
	defer deferTestFileFunc(t, file)

	fileIdent, err := GetFileIdentifierByPathEx(file.Name())
	if err != nil {
		t.Errorf("GetFileIdentifierByPathEx failed: %v", err)
	}

	expected := getFileIDFromCommandEx(t, file)
	if expected.String() != fileIdent.GetFileID().String() {
		t.Errorf("expected.String() != fileIdent.GetFileID().String(); expected: %v, got: %v", expected.String(), fileIdent.GetFileID().String())
	}

	if fileIdent.GetFileID().Cmp(getBigInt(0, 0)) == 0 {
		t.Errorf("fileIdent.GetFileID().Cmp(getBigInt(0, 0)) == 0")
	}
	if fileIdent.GetDeviceID() == 0 {
		t.Errorf("fileIdent.GetDeviceID() == 0")
	}
	if fileIdent.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0 {
		t.Errorf("fileIdent.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0")
	}
}

func TestGetIDAllPossibleValuesEx(t *testing.T) {
	iterateAllFileIdentifierEx(func(expected, expectedFileID *big.Int, vol, idxHi, idxLo uint64, f FileIdentEx) {
		if f.GetDeviceID() != vol {
			t.Errorf("compare failed, expected: %d, received: %d", vol, f.GetDeviceID())
		}

		if expected.Cmp(f.GetGlobalFileID()) != 0 {
			t.Errorf("compare failed, expected: %s, received: %s", expected.String(), f.GetGlobalFileID().String())
		}

		if expectedFileID.Cmp(f.GetFileID()) != 0 {
			t.Errorf("compare failed, expected: %s, received: %s", expectedFileID.String(), f.GetFileID().String())
		}

		checkFileIdentifierBasicEx(t, GetFileIdentifierFromGetGlobalFileIDEx(f.GetGlobalFileID()), f)
	})
}

// iterateAllFileIdentifier is for for tests
func iterateAllFileIdentifierEx(cb func(globalId, expectedFileID *big.Int, vol, idxHi, idxLo uint64, f FileIdentEx)) {
	f := &fileIdentEx{}

	expected := big.NewInt(0)
	iterateAllUint64(18446744073709551615, func(vol uint64) {
		bigIdxVol := getBigInt(vol, 128)
		expected.Add(expected, bigIdxVol)
		iterateAllUint64(18446744073709551615, func(idxHi uint64) {
			bigIdxHi := getBigInt(idxHi, 64)
			expected.Add(expected, bigIdxHi)
			iterateAllUint64(18446744073709551615, func(idxLo uint64) {
				f.vol = vol
				f.idxHi = idxHi
				f.idxLo = idxLo

				bigIdxLo := getBigInt(idxLo, 0)
				expected.Add(expected, bigIdxLo)

				expectedFileID := big.NewInt(0)
				expectedFileID.Add(bigIdxHi, bigIdxLo)

				cb(expected, expectedFileID, vol, idxHi, idxLo, f)

				expected.Sub(expected, bigIdxLo)
			})
			expected.Sub(expected, bigIdxHi)
		})
		expected.Sub(expected, bigIdxVol)
	})
}
