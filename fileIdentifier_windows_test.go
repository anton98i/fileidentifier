package fileidentifier

import (
	"math/big"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func getFileIDFromCommandEx(t *testing.T, file *os.File) *big.Int {
	out, err := exec.Command("fsutil", "file", "queryfileid", file.Name()).Output()
	if err != nil {
		t.Errorf("exec.Command(fsutil file queryfileid %v).Output() failed: %v", file.Name(), err)
	}
	splitted := strings.Split(string(out), "0x")
	if len(splitted) != 2 {
		t.Errorf("out (%s) is no in the correct format, expected someting like: 'Datei-ID: 0x000000000000000000030000000618a1'", out)
		t.FailNow()
	}
	expectedFileID := new(big.Int)
	expectedFileID.SetString(splitted[1], 16)
	return expectedFileID
}

func checkFileIdentifierBasic(t *testing.T, _f, _expected FileIdentifier) {
	f := _f.(*fileIdentifier)
	expected := _expected.(*fileIdentifier)
	if f.vol != expected.vol {
		t.Errorf("checkFileIdentifierBasic vol failed, expected: %d, received: %d", expected.vol, f.vol)
	}
	if f.idxHi != expected.idxHi {
		t.Errorf("checkFileIdentifierBasic idxHi failed, expected: %d, received: %d", expected.idxHi, f.idxHi)
	}
	if f.idxLo != expected.idxLo {
		t.Errorf("checkFileIdentifierBasic idxLo failed, expected: %d, received: %d", expected.idxLo, f.idxLo)
	}
	if f.GetFileID() != expected.GetFileID() {
		t.Errorf("checkFileIdentifierBasic GetFileID failed, expected: %d, received: %d", expected.GetFileID(), f.GetFileID())
	}
}

func TestGetID(t *testing.T) {
	f := fileIdentifier{
		vol:   1234,
		idxHi: 5678,
		idxLo: 90,
	}
	id := f.GetGlobalFileID()
	// result can get calculated by a full precision calculator like: https://www.mathsisfun.com/calculator-precision.html
	// (1234 * 2^64) + (5678 * 2^32)+90id.String()
	if "22763282211344411000922" != id.String() {
		t.Errorf("22763282211344411000922 != id.String(); got %v", id.String())
	}
}

func TestGetID2(t *testing.T) {
	// max: 18446744073709551615
	f := fileIdentifier{
		vol:   4294967294,
		idxHi: 4294967293,
		idxLo: 4294967292,
	}
	id := f.GetGlobalFileID()
	// result can get calculated by a full precision calculator like: https://www.mathsisfun.com/calculator-precision.html
	// (4294967294 * 2^64) + (4294967293 * 2^32)+4294967292
	if "79228162495817593511244464124" != id.String() {
		t.Errorf("79228162495817593511244464124 != id.String(); got %v", id.String())
	}
}

func TestGetFileIdentifier(t *testing.T) {
	file := getTestFile(t)
	defer deferTestFileFunc(t, file)

	f, err := GetFileIdentifierByFile(file)
	if err != nil {
		t.Errorf("GetFileIdentifierByFile failed: %v", err)
	}

	fileinfo, err := file.Stat()
	if err != nil {
		t.Errorf("file.Stat() failed: %v", err)
	}

	_fileIdent, err := GetFileIdentifier(fileinfo)
	if err != nil {
		t.Errorf("GetFileIdentifier error: %v", err)
	}
	fileIdent := _fileIdent.(*fileIdentifier)

	checkFileIdentifierBasic(t, _fileIdent, f)

	expected := getFileIDFromCommand(t, file)
	if expected != fileIdent.GetFileID() {
		t.Errorf("expected.Uint64() != state.GetFileID(); expected: %v, got: %v", expected, fileIdent.GetFileID())
	}

	if fileIdent.GetFileID() == 0 {
		t.Errorf("fileIdent.GetFileID() == 0")
	}
	if fileIdent.GetDeviceID() == 0 {
		t.Errorf("fileIdent.GetDeviceID() == 0")
	}
	if fileIdent.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0 {
		t.Errorf("fileIdent.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0")
	}
}

func TestGetFileIdentifierStat(t *testing.T) {
	file := getTestFile(t)
	defer deferTestFileFunc(t, file)

	f, err := GetFileIdentifierByFile(file)
	if err != nil {
		t.Errorf("GetFileIdentifierByFile failed: %v", err)
	}

	fileinfo, err := os.Stat(file.Name())
	if err != nil {
		t.Errorf("os.Stat(%v) failed: %v", file.Name(), err)
	}

	_fileIdent, err := GetFileIdentifier(fileinfo)
	if err != nil {
		t.Errorf("GetFileIdentifier error: %v", err)
	}
	fileIdent := _fileIdent.(*fileIdentifier)

	checkFileIdentifierBasic(t, _fileIdent, f)

	expected := getFileIDFromCommand(t, file)
	if expected != fileIdent.GetFileID() {
		t.Errorf("expected.Uint64() != state.GetFileID(); expected: %v, got: %v", expected, fileIdent.GetFileID())
	}

	if fileIdent.GetFileID() == 0 {
		t.Errorf("fileIdent.GetFileID() == 0")
	}
	if fileIdent.GetDeviceID() == 0 {
		t.Errorf("fileIdent.GetDeviceID() == 0")
	}
	if fileIdent.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0 {
		t.Errorf("fileIdent.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0")
	}
}

func TestGetIDAllPossibleValues(t *testing.T) {
	iterateAllFileIdentifier(func(expected *big.Int, expectedFileID uint64, vol, idxHi, idxLo uint64, f FileIdentifier) {
		if f.GetDeviceID() != vol {
			t.Errorf("compare failed, expected: %d, received: %d", vol, f.GetDeviceID())
		}

		if expected.Cmp(f.GetGlobalFileID()) != 0 {
			t.Errorf("compare failed, expected: %s, received: %s", expected, f.GetGlobalFileID())
		}

		if expectedFileID != f.GetFileID() {
			t.Errorf("compare failed, expected: %d, received: %d", expectedFileID, f.GetFileID())
		}

		checkFileIdentifierBasic(t, GetFileIdentifierFromGetGlobalFileID(f.GetGlobalFileID()), f)
	})
}

func BenchmarkGetGlobalFileIDAndGetFileIdentifierFromGetGlobalFileID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		iterateAllFileIdentifier(func(expected *big.Int, expectedFileID uint64, vol, idxHi, idxLo uint64, f FileIdentifier) {
			GetFileIdentifierFromGetGlobalFileID(f.GetGlobalFileID())
		})
	}
}

func iterateAllFileIdentifier(cb func(globalId *big.Int, expectedFileID, vol, idxHi, idxLo uint64, f FileIdentifier)) {
	f := &fileIdentifier{}

	expected := big.NewInt(0)
	iterateAllUint64(4294967295, func(vol uint64) {
		bigIdxVol := getBigInt(vol, 64)
		expected.Add(expected, bigIdxVol)
		iterateAllUint64(4294967295, func(idxHi uint64) {
			bigIdxHi := getBigInt(idxHi, 32)
			expected.Add(expected, bigIdxHi)
			iterateAllUint64(4294967295, func(idxLo uint64) {
				f.vol = uint32(vol)
				f.idxHi = uint32(idxHi)
				f.idxLo = uint32(idxLo)

				bigIdxLo := getBigInt(idxLo, 0)
				expected.Add(expected, bigIdxLo)

				expectedFileID := bigIdxHi.Uint64() + bigIdxLo.Uint64()

				cb(expected, expectedFileID, vol, idxHi, idxLo, f)

				expected.Sub(expected, bigIdxLo)
			})
			expected.Sub(expected, bigIdxHi)
		})
		expected.Sub(expected, bigIdxVol)
	})
}
