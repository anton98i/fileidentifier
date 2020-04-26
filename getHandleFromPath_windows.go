package fileidentifier

import (
	"fmt"
	"syscall"

	"github.com/anton98i/fileidentifier/internal/longpath"
)

func getHandleFromPath(path string) (syscall.Handle, error) {
	path, err := longpath.Fix(path)
	if err != nil {
		return syscall.Handle(0), fmt.Errorf("longpath.Fix failed: %v", err)
	}
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return syscall.Handle(0), fmt.Errorf("failed to create UTF16PtrFromString from path %v", path)
	}
	attrs := uint32(syscall.FILE_FLAG_BACKUP_SEMANTICS)
	attrs |= syscall.FILE_FLAG_OPEN_REPARSE_POINT
	h, err := syscall.CreateFile(pathPtr, 0, 0, nil, syscall.OPEN_EXISTING, attrs, 0)
	if err != nil {
		return syscall.Handle(0), fmt.Errorf("syscall.CreateFile failed: %v", err)
	}
	return h, nil
}
