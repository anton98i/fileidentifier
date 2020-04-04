package windows

import (
	"syscall"
	"unsafe"

	"github.com/anton98i/FileIdentifier/internal/windows/sysdll"
)

var _ unsafe.Pointer

var (
	modkernel32 = syscall.NewLazyDLL(sysdll.Add("kernel32.dll"))

	procGetFileInformationByHandleEx = modkernel32.NewProc("GetFileInformationByHandleEx")
)

// Do the interface allocations only once for common
// Errno values.
const (
	errnoErrortIoPending = 997
)

var (
	errErrorIoPending error = syscall.Errno(errnoErrortIoPending)
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoErrortIoPending:
		return errErrorIoPending
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

// GetFileInformationByHandleEx method
func GetFileInformationByHandleEx(handle syscall.Handle, class uint32, info *byte, bufsize uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procGetFileInformationByHandleEx.Addr(), 4, uintptr(handle), uintptr(class), uintptr(unsafe.Pointer(info)), uintptr(bufsize), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
