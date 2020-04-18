// +build !windows

package fileidentifier

import (
	"math/big"
	"os"
	"os/exec"
	"strings"
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
	/*
		  out, err := exec.Command("stat", "-c%i", file.Name()).Output()
			if err != nil {
				t.Errorf("exec.Command(stat -c%%i %v).Output() failed: %v", file.Name(), err)
		  }
	*/
	out, err := exec.Command("ls", "-i", file.Name()).Output()
	if err != nil {
		t.Errorf("exec.Command(ls -i %v).Output() failed: %v", file.Name(), err)
	}
	expectedFileID := new(big.Int)
	/*
		parseSource := strings.TrimRight(string(out), "\r\n")
	*/
	trimmed := strings.Trim(string(out), " \r\n")
	splitted := strings.Split(trimmed, " ")
	if len(splitted) != 2 {
		t.Logf("len(splitted) != 2, got: %d", len(splitted))
	}
	parseSource := splitted[0]
	_, ok := expectedFileID.SetString(parseSource, 10)
	if !ok {
		t.Logf("filename: %q", file.Name())
		t.Logf("stat id output: %q", string(out))
		t.Logf("trimmed: %q", trimmed)
		t.Logf("parseSource: %q", parseSource)
		t.Logf("stat id parsed: %q", expectedFileID.String())
		t.Errorf("expectedFileID.SetString(%s, 16) failed", string(out))
	}
	return expectedFileID
}
