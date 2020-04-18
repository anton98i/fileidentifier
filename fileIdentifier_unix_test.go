// +build !windows

package fileidentifier

import (
	"math/big"
	"os"
	"os/exec"
	"testing"
)

/*
man stat output (cutted):

STAT(1)                                                                                            User Commands                                                                                            STAT(1)
NAME
       stat - display file or file system status
SYNOPSIS
       stat [OPTION]... FILE...
DESCRIPTION
       Display file or file system status.
       Mandatory arguments to long options are mandatory for short options too.
       -L, --dereference
              follow links
       -f, --file-system
              display file system status instead of file status
       -c  --format=FORMAT
              use the specified FORMAT instead of the default; output a newline after each use of FORMAT
       --printf=FORMAT
              like --format, but interpret backslash escapes, and do not output a mandatory trailing newline; if you want a newline, include \n in FORMAT
       -t, --terse
              print the information in terse form
       --help display this help and exit
       --version
              output version information and exit
       The valid format sequences for files (without --file-system):
%a     access rights in octal (note '#' and '0' printf flags)
%A     access rights in human readable form
%b     number of blocks allocated (see %B)
%B     the size in bytes of each block reported by %b
%C     SELinux security context string
%d     device number in decimal
%D     device number in hex
%f     raw mode in hex
%F     file type
%g     group ID of owner
%G     group name of owner
%h     number of hard links
%i     inode number
*/

func getFileIDFromCommandEx(t *testing.T, file *os.File) *big.Int {
	out, err := exec.Command("stat", "-c %i", file.Name()).Output()
	if err != nil {
		t.Errorf("exec.Command(stat -c %i %v).Output() failed: %v", file.Name(), err)
	}
	expectedFileID := new(big.Int)
	expectedFileID.SetString(string(out), 10)
	return expectedFileID
}

func iterateAllFileIdentifier(cb func(globalId, expectedFileID *big.Int, dev, inode uint64)) {
	expected := big.NewInt(0)
	iterateAllUint64(0, func(dev uint64) {
		devBig := getBigInt(dev, 64)
		expected.Add(expected, devBig)
		iterateAllUint64(0, func(inode uint64) {
			inodeBig := getBigInt(inode, 0)
			expected.Add(expected, inodeBig)

			cb(expected, inode, dev, inode)

			expected.Sub(expected, inodeBig)
		})
		expected.Sub(expected, devBig)
	})
}

func TestGetIDAllPossibleValuesUnix(t *testing.T) {
	f := FileIdentifier{}

	iterateAllFileIdentifier(func(globalId *big.Int, expectedFileID uint64, dev, inode uint64) {
		f.device = dev
		f.inode = inode

		if f.GetDeviceID() != vol {
			t.Errorf("f.GetDeviceID() != vol, expected: %d, received: %d", vol, f.GetDeviceID())
		}

		if expected.Cmp(f.GetGlobalFileID()) != 0 {
			t.Errorf("expected.Cmp(f.GetGlobalFileID()) != 0, expected: %s, received: %s", expected.String(), f.GetID().String())
		}

		if inodeBig != f.GetFileID() {
			t.Errorf("inodeBig != f.GetFileID(), expected: %d, received: %d", inodeBig, f.GetFileID())
		}

		newF := GetFileIdentifierFromGetGlobalFileID(f.GetGlobalFileID())
		if newF.device != device {
			t.Errorf("GetFileIdentifierFromGetGlobalFileID device failed, expected: %d, received: %d", device, newF.device)
		}
		if newF.inode != inode {
			t.Errorf("GetFileIdentifierFromGetGlobalFileID inode failed, expected: %d, received: %d", inode, newF.inode)
		}
	})
}
