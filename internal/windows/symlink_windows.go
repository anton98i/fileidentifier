// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package windows

// FileInformationClass values
const (
	FileBasicInfo                  = 0    // FILE_BASIC_INFO
	FileStandardInfo               = 1    // FILE_STANDARD_INFO
	FileNameInfo                   = 2    // FILE_NAME_INFO
	FileStreamInfo                 = 7    // FILE_STREAM_INFO
	FileCompressionInfo            = 8    // FILE_COMPRESSION_INFO
	FileAttributeTagInfo           = 9    // FILE_ATTRIBUTE_TAG_INFO
	FileIDBothDirectoryInfo        = 0xa  // FILE_ID_BOTH_DIR_INFO
	FileIDBothDirectoryRestartInfo = 0xb  // FILE_ID_BOTH_DIR_INFO
	FileRemoteProtocolInfo         = 0xd  // FILE_REMOTE_PROTOCOL_INFO
	FileFullDirectoryInfo          = 0xe  // FILE_FULL_DIR_INFO
	FileFullDirectoryRestartInfo   = 0xf  // FILE_FULL_DIR_INFO
	FileStorageInfo                = 0x10 // FILE_STORAGE_INFO
	FileAlignmentInfo              = 0x11 // FILE_ALIGNMENT_INFO
	FileIDInfo                     = 0x12 // FILE_ID_INFO
	FileIDExtdDirectoryInfo        = 0x13 // FILE_ID_EXTD_DIR_INFO
	FileIDExtdDirectoryRestartInfo = 0x14 // FILE_ID_EXTD_DIR_INFO
)

//sys	GetFileInformationByHandleEx(handle syscall.Handle, class uint32, info *byte, bufsize uint32) (err error)
