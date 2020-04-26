package fileidentifier

import (
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func cutBigIntToUint64(n *big.Int) uint64 {
	tmp := getBigInt(0, 0)
	return tmp.And(n, getBigInt(18446744073709551615, 0)).Uint64()
}

func TestGetFileIDFromCommand(t *testing.T) {
	overFlowTo0, _ := getBigInt(0, 0).SetString("18446744073709551616", 10)
	overFlowTo1, _ := getBigInt(0, 0).SetString("18446744073709551617", 10)
	testArr := []struct {
		In  *big.Int
		Out uint64
	}{
		{
			In:  getBigInt(18446744073709551615, 0),
			Out: 18446744073709551615,
		},
		{
			In:  overFlowTo0,
			Out: 0,
		},
		{
			In:  overFlowTo1,
			Out: 1,
		},
	}
	for _, te := range testArr {
		got := cutBigIntToUint64(te.In)
		if got != te.Out {
			t.Errorf("cutBigIntToUint64 got false result, expected: %v, got: %v", te.Out, got)
		}
	}
}

func getFileIDFromCommandByName(t *testing.T, filename string) uint64 {
	return cutBigIntToUint64(getFileIDFromCommandExByName(t, filename))
}

func TestFiles(t *testing.T) {
	testArr := []string{"./README.md"}
	for _, testPath := range testArr {
		testFileByPath(t, testPath)
		testFileByPathEx(t, testPath)
	}
}

func TestDirectories(t *testing.T) {
	testArr := []string{"./"}
	for _, testPath := range testArr {
		testFileByPath(t, testPath)
		testFileByPathEx(t, testPath)
	}
}

func expectDirectoryExists(t *testing.T, path string) string {
	_, err := os.Stat(path)
	if err != nil {
		pathAbsolute, err := filepath.Abs(path)
		if err != nil {
			t.Errorf("filepath.Abs(%s) failed: %v", path, err)
			t.FailNow()
		}
		err = os.MkdirAll(pathAbsolute, 0777)
		if err != nil {
			t.Errorf("failed to create test directory: %v", err)
			t.FailNow()
		}
	}
	return path
}

func expectFileExists(t *testing.T, fileName string) string {
	_, err := os.Stat(fileName)
	// if os.IsNotExist(err) {
	if err != nil {
		pathAbsolute, err := filepath.Abs(fileName)
		if err != nil {
			t.Errorf("filepath.Abs(%s) failed: %v", fileName, err)
			t.FailNow()
		}
		file, err := os.Create(pathAbsolute)
		if err != nil {
			t.Errorf("failed to create test file %s:  %v", fileName, err)
			t.FailNow()
		}
		err = file.Close()
		if err != nil {
			t.Errorf("failed to close newly created file: %v", err)
		}
	}
	return fileName
}

func TestLongPath(t *testing.T) {
	longPathRelative := expectDirectoryExists(t, filepath.Join(
		"./test",
		strings.Repeat("1", 50),
		strings.Repeat("2", 50),
		strings.Repeat("3", 50),
		strings.Repeat("4", 50),
		strings.Repeat("5", 50),
		strings.Repeat("6", 50),
		strings.Repeat("7", 50),
		strings.Repeat("8", 50),
		strings.Repeat("9", 50),
	))
	fileName := expectFileExists(t, filepath.Join(longPathRelative, "testfile.txt"))
	// mkdirAll fails on creating a single long directory name => not to long single name
	longPath2 := expectDirectoryExists(t, filepath.Join("./test", strings.Repeat("a", 100), strings.Repeat("b", 100), strings.Repeat("c", 100)))
	fileName2 := expectFileExists(t, filepath.Join(longPath2, "testfile2.txt"))
	testArr := []string{longPathRelative, fileName, longPath2, fileName2}
	for _, testPath := range testArr {
		if filepath.IsAbs(testPath) {
			t.Errorf("Relative path test has relative path: %v", testPath)
		}
		testFileByPathNoConsoleCheck(t, testPath)
		testFileByPathNoConsoleCheckEx(t, testPath)
		pathAbsolute, err := filepath.Abs(testPath)
		if err != nil {
			t.Errorf("filepath.Abs(%s) failed: %v", testPath, err)
		}
		testFileByPathNoConsoleCheck(t, pathAbsolute)
		testFileByPathNoConsoleCheckEx(t, pathAbsolute)
	}
}

func testFileByPath(t *testing.T, path string) {
	fileIdent := testFileByPathNoConsoleCheck(t, path)

	expected := getFileIDFromCommandByName(t, path)
	if expected != fileIdent.GetFileID() {
		t.Errorf("expected.Uint64() != state.GetFileID(); expected: %v, got: %v", expected, fileIdent.GetFileID())
	}
}

func testFileByPathEx(t *testing.T, path string) {
	fileIdent := testFileByPathNoConsoleCheckEx(t, path)

	expected := getFileIDFromCommandExByName(t, path)
	if expected.Cmp(fileIdent.GetFileID()) != 0 {
		t.Errorf("testFileByPathEx/expected.Cmp(fileIdent.GetFileID()) != 0; expected: %v, got: %v", expected.String(), fileIdent.GetFileID().String())
	}
}

func testFileByPathNoConsoleCheck(t *testing.T, path string) FileIdentifier {
	fileIdent, err := GetFileIdentifierByPath(path)
	if err != nil {
		t.Errorf("GetFileIdentifierByPath error: %v; path: %s", err, path)
	}

	if fileIdent.GetFileID() == 0 {
		t.Errorf("testFileByPath/fileIdent.GetFileID() == 0; path: %s", path)
	}
	if fileIdent.GetDeviceID() == 0 {
		t.Errorf("testFileByPath/fileIdent.GetDeviceID() == 0; path: %s", path)
	}
	if fileIdent.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0 {
		t.Errorf("testFileByPath/fileIdent.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0; path: %s", path)
	}

	return fileIdent
}

func testFileByPathNoConsoleCheckEx(t *testing.T, path string) FileIdentEx {
	fileIdent, err := GetFileIdentifierByPathEx(path)
	if err != nil {
		t.Errorf("GetFileIdentifierByPathEx error: %v; path %s", err, path)
	}

	if fileIdent.GetFileID().Cmp(getBigInt(0, 0)) == 0 {
		t.Errorf("testFileByPathEx/fileIdent.GetFileID().Cmp(getBigInt(0, 0)) == 0; path: %s", path)
	}
	if fileIdent.GetDeviceID() == 0 {
		t.Errorf("testFileByPathEx/fileIdent.GetDeviceID() == 0; path: %s", path)
	}
	if fileIdent.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0 {
		t.Errorf("testFileByPathEx/fileIdent.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0; path: %s", path)
	}
	return fileIdent
}

func testIDsBasic(t *testing.T, f FileIdentifier, expectedGlobalID string, expectedFileID, expectedDiviceID uint64) {
	testFunction := func(f FileIdentifier) {
		globalID := f.GetGlobalFileID().String()
		if globalID != expectedGlobalID {
			t.Errorf("testIDsBasic/GetGlobalFileID/expected: %s\ngot %s", expectedGlobalID, globalID)
		}
		fileID := f.GetFileID()
		if fileID != expectedFileID {
			t.Errorf("testIDsBasic/GetFileID/expected: %d\ngot %d", expectedFileID, fileID)
		}
		deviceID := f.GetDeviceID()
		if deviceID != expectedDiviceID {
			t.Errorf("testIDsBasic/GetDeviceID/expected: %v\ngot %v", expectedDiviceID, deviceID)
		}
	}
	testFunction(f)
	f2 := GetFileIdentifierFromGetGlobalFileID(f.GetGlobalFileID())
	testFunction(f2)
}

func testIDsBasicEx(t *testing.T, f FileIdentEx, expectedGlobalID string, expectedFileID string, expectedDiviceID uint64) {
	testFunction := func(f FileIdentEx, testname string) {
		globalID := f.GetGlobalFileID().String()
		if globalID != expectedGlobalID {
			t.Errorf("testIDsBasicEx/GetGlobalFileID/%s/expected: %s\ngot %s", testname, expectedGlobalID, globalID)
		}
		fileID := f.GetFileID().String()
		if fileID != expectedFileID {
			t.Errorf("testIDsBasicEx/GetFileID/%s/expected: %s\ngot %s", testname, expectedFileID, fileID)
		}
		deviceID := f.GetDeviceID()
		if deviceID != expectedDiviceID {
			t.Errorf("testIDsBasicEx/GetDeviceID/%s/expected: %d\ngot %d", testname, expectedDiviceID, deviceID)
		}
	}
	testFunction(f, "default")
	f2 := GetFileIdentifierFromGetGlobalFileIDEx(f.GetGlobalFileID())
	testFunction(f2, "after GetFileIdentifierFromGetGlobalFileIDEx from GetFileIdentifierFromGetGlobalFileIDEx")
}
