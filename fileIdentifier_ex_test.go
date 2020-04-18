package fileidentifier

import "testing"

func TestGetFileIdentifierByPathEx(t *testing.T) {
	file := getTestFile(t)
	defer deferTestFileFunc(t, file)

	f, err := GetFileIdentifierByPathEx(file.Name())
	if err != nil {
		t.Errorf("GetFileIdentifierByPathEx(%s) failed: %v", file.Name(), err)
		t.FailNow()
	}

	fileID := getFileIDFromCommandEx(t, file)
	if fileID.String() != f.GetFileID().String() {
		t.Errorf("fileID.String() != f.GetFileID().String(): %v != %v", fileID.String(), f.GetFileID().String())
	}

	if f.GetFileID().Cmp(getBigInt(0, 0)) == 0 {
		t.Errorf("f.GetFileID().Cmp(getBigInt(0, 0)) == 0")
	}
	if f.GetDeviceID() == 0 {
		t.Errorf("f.GetDeviceID() == 0")
	}
	if f.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0 {
		t.Errorf("f.GetGlobalFileID().Cmp(getBigInt(0, 0)) == 0")
	}
}
