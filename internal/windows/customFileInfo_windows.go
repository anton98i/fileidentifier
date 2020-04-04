package windows

// FileIDInfoStrcut is representation of FILE_ID_INFO: https://docs.microsoft.com/de-de/windows/win32/api/winbase/ns-winbase-file_id_info
type FileIDInfoStrcut struct {
	// ULONGLONG = typedef unsigned __int64 ULONGLONG;: https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/c57d9fba-12ef-4853-b0d5-a6f472b50388
	// ULONGLONG   VolumeSerialNumber;
	VolumeSerialNumber uint64
	// FILE_ID_128 https://docs.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-file_id_128
	// typedef struct _FILE_ID_128 {
	//  BYTE Identifier[16];
	// }
	// FILE_ID_128 FileId;
	FileID struct {
		arr [16]byte
	}
}

func bytes2String(b []byte) uint64 {
	var ret uint64
	for i := uint64(0); i < 8; i++ {
		ret += uint64(b[i]) << (i * 8)
	}
	return ret
}

// GetFileID method
func (f FileIDInfoStrcut) GetFileID() (uint64, uint64) {
	return bytes2String(f.FileID.arr[8:16]), bytes2String(f.FileID.arr[:8])
}
