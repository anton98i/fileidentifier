package fileidentifier

import (
	"io/ioutil"
	"math/big"
	"os"
	"testing"
)

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

func getFileIDFromCommand(t *testing.T, file *os.File) uint64 {
	return cutBigIntToUint64(getFileIDFromCommandEx(t, file))
}

func TestGetFileIdentifierByPath(t *testing.T) {
	file := getTestFile(t)
	defer deferTestFileFunc(t, file)

	f, err := GetFileIdentifierByPath(file.Name())
	if err != nil {
		t.Errorf("GetFileIdentifierByPath(%s) failed: %v", file.Name(), err)
		t.FailNow()
	}

	fileID := getFileIDFromCommand(t, file)
	if fileID != f.GetFileID() {
		t.Errorf("fileID != f.GetFileID(): %v != %v", fileID, f.GetFileID())
	}

	if f.GetFileID() == 0 {
		t.Errorf("f.GetFileID() == 0")
	}
	if f.GetDeviceID() == 0 {
		t.Errorf("f.GetDeviceID() == 0")
	}
	if f.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0 {
		t.Errorf("f.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0")
	}

}

func iterateAllUint64(max uint64, cb func(count uint64)) {
	var i uint64
	var addvalue uint64
	for lastValie := i; i >= lastValie && i < max; i = i + addvalue {
		cb(i)
		addvalue = addvalue*3000/2 + 100
		lastValie = i
	}
	cb(max)
}
