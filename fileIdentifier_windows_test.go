package fileidentifier

import (
	"io/ioutil"
	"math/big"
	"os"
	"testing"
)

func checkFileIdentifierBasic(t *testing.T, f, expected FileIdentifier) {
	if f.vol != expected.vol {
		t.Errorf("GetFileIdentifierFromGetGlobalFileID vol failed, expected: %d, received: %d", expected.vol, f.vol)
	}
	if f.idxHi != expected.idxHi {
		t.Errorf("GetFileIdentifierFromGetGlobalFileID idxHi failed, expected: %d, received: %d", expected.idxHi, f.idxHi)
	}
	if f.idxLo != expected.idxLo {
		t.Errorf("GetFileIdentifierFromGetGlobalFileID idxLo failed, expected: %d, received: %d", expected.idxLo, f.idxLo)
	}
	if f.GetFileID().Cmp(expected.GetFileID()) != 0 {
		t.Errorf("GetFileIdentifierFromGetGlobalFileID GetFileID failed, expected: %d, received: %d", expected.GetFileID(), f.GetFileID())
	}
}

func getTestFile(t *testing.T) *os.File {
	file, err := ioutil.TempFile("", "fileidentifier_test")
	if err != nil {
		t.Errorf("failed to create testfile, errormsg: %v", err)
	}
	return file
}

func deferTestFileFunc(t *testing.T, file *os.File) {
	fileName := file.Name()
	err := file.Close()
	if err != nil {
		t.Errorf("failed to close testfile, errormsg: %v", err)
	}
	err = os.Remove(fileName)
	if err != nil {
		t.Errorf("failed to delete testfile, errormsg: %v", err)
	}
}

func TestGetFileIdentifier(t *testing.T) {
	file := getTestFile(t)
	defer deferTestFileFunc(t, file)

	f, err := GetFileIdentifierByFile(file)
	expectNil(t, err)

	fileinfo, err := file.Stat()
	expectNil(t, err)

	state := GetFileIdentifier(fileinfo)

	checkFileIdentifierBasic(t, state, *f)

	// idxHi may be 0
	// expectTrue(t, state.idxHi > 0)
	expectTrue(t, state.idxLo > 0)
	expectTrue(t, state.vol > 0)

	expectTrue(t, state.GetGlobalFileID().String() != "")
}

func TestGetFileIdentifierStat(t *testing.T) {
	file := getTestFile(t)
	defer deferTestFileFunc(t, file)

	f, err := GetFileIdentifierByFile(file)
	expectNil(t, err)

	fileinfo, err := os.Stat(file.Name())
	expectNil(t, err)

	state := GetFileIdentifier(fileinfo)

	checkFileIdentifierBasic(t, state, *f)

	// idxHi may be 0
	// expectTrue(t, state.idxHi > 0)
	expectTrue(t, state.idxLo > 0)
	expectTrue(t, state.vol > 0)

	expectTrue(t, state.GetGlobalFileID().String() != "")
}

/* deferend as refs and normal
func TestGetID(t *testing.T) {
	f := FileIdentifier{
		vol:   1234,
		idxHi: 5678,
		idxLo: 90,
	}
	id := f.GetGlobalFileID()
	// result can get calculated by a full precision calculator like: https://www.mathsisfun.com/calculator-precision.html
	// (1234 * 2^128) + (5678 * 2^64)+90
	expect(t, "419908440780438064018544878421324807012442", id.String())
}

func TestGetID2(t *testing.T) {
	// max: 18446744073709551615
	f := FileIdentifier{
		vol:   18446744073709551614,
		idxHi: 18446744073709551613,
		idxLo: 18446744073709551612,
	}
	id := f.GetGlobalFileID()
	// result can get calculated by a full precision calculator like: https://www.mathsisfun.com/calculator-precision.html
	// (18446744073709551614 * 2^128) + (18446744073709551613 * 2^64)+18446744073709551612
	expect(t, "6277101735386680763495507056286727952602087348884847198204", id.String())
}
*/

var maxVol uint64
var maxIdxHi uint64
var maxIdxLo uint64

func init() {
	f := FileIdentifier{
		vol:   0,
		idxHi: 0,
		idxLo: 0,
	}
	f.vol--
	f.idxHi--
	f.idxLo--
	maxVol = uint64(f.vol)
	maxIdxHi = uint64(f.idxHi)
	maxIdxLo = uint64(f.idxLo)
}

func TestGetIDAllPossibleValues(t *testing.T) {
	iterateAllFileIdentifier(func(expected, expectedFileID *big.Int, vol, idxHi, idxLo uint64, f FileIdentifier) {
		if f.GetDeviceID() != vol {
			t.Errorf("compare failed, expected: %d, received: %d", vol, f.GetDeviceID())
		}

		if expected.Cmp(f.GetGlobalFileID()) != 0 {
			t.Errorf("compare failed, expected: %s, received: %s", expected.String(), f.GetGlobalFileID().String())
		}

		if expectedFileID.Cmp(f.GetFileID()) != 0 {
			t.Errorf("compare failed, expected: %s, received: %s", expectedFileID.String(), f.GetFileID().String())
		}

		checkFileIdentifierBasic(t, GetFileIdentifierFromGetGlobalFileID(f.GetGlobalFileID()), f)
	})
}

func BenchmarkGetGlobalFileIDAndGetFileIdentifierFromGetGlobalFileID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		iterateAllFileIdentifier(func(expected, expectedFileID *big.Int, vol, idxHi, idxLo uint64, f FileIdentifier) {
			GetFileIdentifierFromGetGlobalFileID(f.GetGlobalFileID())
		})
	}
}
